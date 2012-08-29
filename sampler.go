package scriptModeling

// import "fmt"
import "math"
import "math/rand"

type Sampler struct {
  eventPosPrior float64
  eventNegPrior float64
  eventlmPrior float64
  participantPosPrior float64
  participantNegPrior float64
  participantlmPrior float64
  nu_0 float64
  v_0 []float64
  rho []float64
  eventProbCache [][]float64
  model Model
}

func NewSampler(ePprior float64, eNprior float64, elmprior float64, pPprior float64, pNprior float64, plmprior float64, rho0 float64, nu0 float64, model Model) *Sampler {
  sampler := new(Sampler)
  sampler.eventPosPrior = ePprior
  sampler.eventNegPrior = eNprior
  sampler.eventlmPrior = elmprior
  sampler.participantPosPrior = pPprior
  sampler.participantNegPrior = pNprior
  sampler.participantlmPrior = plmprior
  sampler.nu_0 = nu0
  sampler.v_0 = vPrior(rho0)
  // NOTE: Correct?? Initially rho = v_0 ??
  sampler.rho = sampler.v_0
  sampler.model = model
  return sampler
}

func (sampler *Sampler)PickVariable(esd *ESD) {
  //select which random variable to resample; 0:t  1:v  2:rho
  rr := rand.Intn(4)
  if rr ==0 && esd.hasParticipants() {
    sampler.Resample_p(*esd, Pick_participant(*esd))
  } else if rr==1{
    sampler.Resample_v(*esd, pick_invcount(esd.V))
  } else if rr==2{
    sampler.Resample_t(*esd, pick_event(esd.Events.Tau))
  } else {
    Resample_rho()
  }
}

func (sampler *Sampler) Resample_t(esd ESD, target int) {
  var update, lgamma, documentLikelihood, docPositive, docNegative, docNormalize float64
  var newLabel int
  var alternatives []int
  var labels [][]int
  alternatives = newTargets(esd, target)
  labels = getLabels(esd, target, alternatives)
  // decrement counts for current target event, and all words in ESD
  (sampler.model).DecrementCounts(target, -1, esd.Events.Label, esd.Participants.Label, -1, esd.Events.Words, esd.Participants.Words, "event")
  // compute switch-likelihood
  distribution := make([]float64, len(alternatives))
  for idx,alternative := range(alternatives) {
    lgamma = 0.0
    for k:=0 ; k<numTop ; k++ {
      update=0.0
      if k==alternative {
	update=1.0
      }
      docPositive,_ = math.Lgamma(float64(sampler.model.eventtype_histogram[k])+sampler.eventPosPrior+update)
      docNegative,_ = math.Lgamma(float64(sampler.model.numESDs-sampler.model.eventtype_histogram[k])+sampler.eventNegPrior-update)
      docNormalize,_ = math.Lgamma(float64(sampler.model.numESDs)+sampler.eventPosPrior+sampler.eventNegPrior+update)
      lgamma += ((docPositive+docNegative)-docNormalize)
    }
    documentLikelihood = sampler.documentLikelihood(esd.Events.Words, labels[idx])
    distribution[idx]=lgamma+documentLikelihood
  }
  // sample new label
  newLabel = getAccumulativeSample(distribution)
  // update model & esd
   esd.flipEvent(target, alternatives[newLabel])
   esd.Events.Label = labels[newLabel]
   (sampler.model).ReassignCounts(esd, target, alternatives[newLabel], -1, "event")
}

func (sampler *Sampler) Resample_v(esd ESD, target int) {
  var label []int
  var proposedV [numTop-1]int
  var newV int
  // decrement global inversion count for target eventtype
  (sampler.model).DecrementCounts(target, -1, esd.Events.Label, esd.Participants.Label, esd.V[target], esd.Events.Words, esd.Participants.Words, "inversion")
  proposedV = esd.V
  distribution := make([]float64, numTop-target)
  // try every possible value
  for k:=0 ; k<numTop-target ; k++ {
    proposedV[target]=k
    label = computeZ(esd.Events.Tau, computePi(proposedV))
    // NOTE: I am using the **unnormalized log of GMM(target; rho_target)** (Chen does the same!!)
    // NOTE: I am using 'k+1' below as my topicIDs start with 0 ...shouldn't matter, right???
    distribution[k] = -sampler.rho[target] * float64(k+1) + sampler.documentLikelihood(esd.Events.Words, label)
  }
  // sample new value
  newV = getAccumulativeSample(distribution)
  // update model & esd
  esd.V[target] = newV
  esd.Pi = computePi(esd.V)
  esd.Events.Label = computeZ(esd.Events.Tau, esd.Pi)
  (sampler.model).ReassignCounts(esd, target, newV, -1, "inversion")
}

func (sampler *Sampler) Resample_p(esd ESD, targets [2]int) {
  var lgamma, update, pPositive, pNegative, pNormalize, documentLikelihood float64
  var distribution []float64
  var newV int
  eventID := targets[0]
  participantID := targets[1]
  event := esd.Events.Label[eventID]
  target := esd.Participants.Label[eventID][participantID]
  // Get alternative participant types
  alternatives := getAlternatives(target, esd.Participants.Label[eventID])
  // Decrement Counts
  (sampler.model).DecrementCounts(eventID, target, esd.Events.Label, esd.Participants.Label, -1, esd.Events.Words, esd.Participants.Words, "participant")
  // Compute likelihood for every types
  distribution = make([]float64, len(alternatives))
  for idx, proposedV := range(alternatives) {
    lgamma = 0.0
    for i:=0 ; i<numPar ; i++ {
      update = 0.0
      if i==proposedV {update = 1.0}
      pPositive, _ = math.Lgamma(float64(sampler.model.participanttype_eventtype_histogram[proposedV][event]) + sampler.participantPosPrior + update)
      pNegative, _ = math.Lgamma(float64(sampler.model.participanttype_histogram[proposedV]-sampler.model.participanttype_eventtype_histogram[proposedV][event]) + sampler.participantNegPrior - update)
      pNormalize, _ = math.Lgamma(float64(sampler.model.participanttype_histogram[proposedV])+sampler.participantPosPrior+sampler.participantNegPrior+update)
      lgamma = ((pPositive+pNegative)-pNormalize)
    }
    documentLikelihood = sampler.documentLikelihood(esd.Events.Words, esd.Events.Label)/* + sampler.documentLikelihood(esd.Participants.Words, esd.Participants.Label)*/
    distribution[idx]=lgamma+documentLikelihood
  }
  newV = getAccumulativeSample(distribution)
 (sampler.model).ReassignCounts(esd, target, newV, eventID, "participant")
  esd.Participants.Label[eventID][participantID] = newV
}