 package scriptModeling

import "fmt"
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
  v_0 [numTop-1]float64
  rho [numTop-1]float64
  eventProbCache [][]float64
  Model Model
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
  sampler.Model = model
  return sampler
}

func (sampler *Sampler)PickVariable(esd *ESD) {
  //select which random variable to resample; 0:t  1:v  2:rho
  rr := rand.Intn(4)
  if rr ==0 && esd.hasParticipants() {
    fmt.Println("Resampling P...\n")
//     sampler.Resample_p(esd, Pick_participant(esd.Label))
  } else if rr==1{
    fmt.Println("Resampling V...\n")
//     sampler.Resample_v(esd, pick_invcount(esd.V))
  } else if rr==2{
    fmt.Println("Resampling T...\n")
    sampler.Resample_t(esd, pick_event(esd.Tau))
  } else {
    fmt.Println("Resampling Rho...\n")
    sampler.Resample_rho()
  }
}

func (sampler *Sampler) Resample_t(esd *ESD, target int) {
  var update, lgamma, documentLikelihood, docPositive, docNegative, docNormalize float64
  var newLabel int
  var alternatives []int
  alternatives = newTargets(esd.Tau, target)
  proposedLabels := make([]Label, len(alternatives))
  // decrement counts for current target event, and all words in ESD
  sampler.Model.eventtype_histogram[target]--
  sampler.Model.UpdateEventWordCounts(esd.Label, -1)
  sampler.Model.UpdateEventParticipantCountsAll(esd.Label, -1)
  // compute switch-likelihood
  distribution := make([]float64, len(alternatives))
  for idx,alternative := range(alternatives) {
    proposedLabels[idx] = updateLabelingT(target, alternative, esd.Label)
    lgamma = 0.0
    for k:=0 ; k<numTop ; k++ {
      update=0.0
      if k==alternative {update=1.0}
      docPositive,_ = math.Lgamma(float64(sampler.Model.eventtype_histogram[k])+sampler.eventPosPrior+update)
      docNegative,_ = math.Lgamma(float64(sampler.Model.numESDs-sampler.Model.eventtype_histogram[k])+sampler.eventNegPrior-update)
      docNormalize,_ = math.Lgamma(float64(sampler.Model.numESDs)+sampler.eventPosPrior+sampler.eventNegPrior+update)
      lgamma += ((docPositive+docNegative)-docNormalize)
    }
    documentLikelihood = sampler.documentLikelihood(proposedLabels[idx])
    distribution[idx]=lgamma+documentLikelihood
  }
  // sample new label
  newLabel = getAccumulativeSample(distribution)
  // update model & esd
   esd.flipEvent(target, alternatives[newLabel])
   esd.UpdateLabelingT(target, alternatives[newLabel])
   sampler.Model.eventtype_histogram[alternatives[newLabel]]++
   sampler.Model.UpdateEventWordCounts(esd.Label, 1)
   sampler.Model.UpdateEventParticipantCountsAll(esd.Label, 1)
}

func (sampler *Sampler) Resample_v(esd *ESD, target int) {
  var proposedV [numTop-1]int
  var newV int
  var realized bool
  var documentLikelihood float64
  proposedLabels := make([]Label, numTop-target)
  // check if target is realized in esd 
  for eID,_ := range(esd.Label) {
    if eID==target {realized=true}
  }
  if realized==true {
    sampler.Model.UpdateEventWordCounts(esd.Label, -1)
    sampler.Model.UpdateEventParticipantCountsAll(esd.Label, -1)
  }
  // decrement global inversion count for target eventtype
  sampler.Model.invcount_histogram[target] -= esd.V[target]
  proposedV = esd.V
  distribution := make([]float64, numTop-target)
  // try every possible value
  for k:=0 ; k<numTop-target ; k++ {
    proposedV[target]=k 
    // NOTE: I am using the **unnormalized log of GMM(target; rho_target)** (Chen does the same!!)
    // NOTE: I am using 'k+1' below as my topicIDs start with 0 ...shouldn't matter, right???
    distribution[k] = -sampler.rho[target] * float64(k+1)
    // compute documentLikelihood if eventtype for which inv count is resampled is realized in esd
    if realized==true{
      proposedLabels[k] = UpdateLabelingV(esd.Tau, computePi(proposedV), esd.EventLabel, esd.Label)
      documentLikelihood = sampler.documentLikelihood(proposedLabels[k])
      distribution[k]+=documentLikelihood
    }
  }
  // sample new value
  newV = getAccumulativeSample(distribution)
  // update model & esd
  esd.V[target] = newV
  esd.Pi = computePi(esd.V)
  esd.UpdateLabelingV()
  sampler.Model.invcount_histogram[target] += esd.V[target]
  if realized==true {
    sampler.Model.UpdateEventWordCounts(esd.Label, 1)
    sampler.Model.UpdateEventParticipantCountsAll(esd.Label, 1)
  }
}

func (sampler *Sampler) Resample_p(esd *ESD, targets [2]int) {
  var lgamma, update, pPositive, pNegative, pNormalize, documentLikelihood float64
  var distribution []float64
  var newV int
  event := targets[0]
  target := targets[1]
  // Get alternative participant types
  alternatives := getAlternatives(target, esd.Label[event].Participants) 
  proposedLabels := make([]Label, len(alternatives))  
  // Decrement Counts
  sampler.Model.participanttype_histogram[target]--
  sampler.Model.UpdateParticipantWordCounts(target, esd.Label[event].Participants[target], -1)
  sampler.Model.participanttype_eventtype_histogram[target][event]--
  // Compute likelihood for every types
  distribution = make([]float64, len(alternatives))
  for idx, proposedP := range(alternatives) {
    if idx==0 {
      proposedLabels[idx]=esd.Label
    } else {
      esd.UpdateLabelingP(event, alternatives[idx-1], proposedP)
      proposedLabels[idx]=esd.Label
    }
    target=alternatives[idx]
    lgamma = 0.0
    for i:=0 ; i<numPar ; i++ {
      update = 0.0
      if i==proposedP {update = 1.0}
      pPositive, _ = math.Lgamma(float64(sampler.Model.participanttype_eventtype_histogram[proposedP][event]) + sampler.participantPosPrior + update)
      pNegative, _ = math.Lgamma(float64(sampler.Model.participanttype_histogram[proposedP]-sampler.Model.participanttype_eventtype_histogram[proposedP][event]) + sampler.participantNegPrior - update)
      pNormalize, _ = math.Lgamma(float64(sampler.Model.participanttype_histogram[proposedP])+sampler.participantPosPrior+sampler.participantNegPrior+update)
      lgamma += ((pPositive+pNegative)-pNormalize)
    }
    documentLikelihood = sampler.documentLikelihood(proposedLabels[idx])
    distribution[idx]=lgamma+documentLikelihood
  }
  newV = getAccumulativeSample(distribution)
  //update esd and model
  esd.UpdateLabelingP(event, alternatives[len(alternatives)-1], alternatives[newV])
  sampler.Model.participanttype_histogram[alternatives[newV]]++
  sampler.Model.participanttype_eventtype_histogram[alternatives[newV]][event]++
  sampler.Model.UpdateParticipantWordCounts(alternatives[newV], esd.Label[event].Participants[alternatives[newV]], 1)
}