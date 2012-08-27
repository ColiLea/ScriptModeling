package scriptModeling

import "fmt"
import "math"
import "math/rand"

type Sampler struct {
  eventPosPrior float64
  eventNegPrior float64
  lmPrior float64
  nu_0 float64
  v_0 []float64
  rho []float64
  eventProbCache [][]float64
  model Model
}

func NewSampler(ePprior float64, eNprior float64, lmprior float64, rho0 float64, nu0 float64, model Model) *Sampler {
  sampler := new(Sampler)
  sampler.eventPosPrior = ePprior
  sampler.eventNegPrior = eNprior
  sampler.nu_0 = nu0
  sampler.v_0 = vPrior(rho0)
  // NOTE: Correct?? Initially rho = v_0 ??
  sampler.rho = sampler.v_0
  sampler.lmPrior = lmprior
  sampler.model = model
  return sampler
}

func (sampler *Sampler)PickVariable(esd *ESD) {
  //select which random variable to resample; 0:t  1:v  2:rho
  rr := rand.Intn(3)
  if rr==0{
    sampler.resample_t(*esd, pick_event(esd.Tao))
  } else if rr==1{
    sampler.resample_v(*esd, pick_invcount(esd.V))
  } else {
    resample_rho()
  }
}

func (sampler *Sampler) resample_t(esd ESD, target int) {
  var update, lgamma, documentLikelihood, docPositive, docNegative, docNormalize float64
  var newLabel int
  var alternatives []int
  var labels [][]int
  alternatives = newTargets(esd, target)
  labels = getLabels(esd, target, alternatives)
  // decrement counts for current target event, and all words in ESD
  (sampler.model).DecrementCounts(target, esd.Label, esd.Events)
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
    documentLikelihood = sampler.documentLikelihood(esd.Events, labels[idx])
    distribution[idx]=lgamma+documentLikelihood
    fmt.Println("eventPrior: ", lgamma, "DocLikelihood", documentLikelihood)
  }
  // sample new label
  newLabel = getAccumulativeSample(distribution)
  // update model & esd
   (sampler.model).ReassignCounts(alternatives[newLabel], labels[newLabel], esd.Events)
   esd.flipEvent(target, alternatives[newLabel])
   esd.Label = labels[newLabel]
}



func (sampler *Sampler) resample_v(esd ESD, target int) {
  var label []int
  var proposedV [numTop-1]int
  var newV int
  // decrement global inversion count for target eventtype
  sampler.model.invcount_histogram[target] -= esd.V[target]
  proposedV = esd.V
  distribution := make([]float64, numTop-target)
  // try every possible value
  for k:=0 ; k<numTop-target ; k++ {
    proposedV[target]=k
    label = computeZ(esd.Tao, computePi(proposedV))
    // NOTE: I am using the **unnormalized log of GMM(target; rho_target)** (Chen does the same!!)
    // NOTE: I am using 'k+1' below as my topicIDs start with 0 ...shouldn't matter, right???
    distribution[k] = -sampler.rho[target] * float64(k+1) + sampler.documentLikelihood(esd.Events, label)
    fmt.Println("inversionPrior: ", -sampler.rho[target] * float64(k+1), "DocLikelihood", sampler.documentLikelihood(esd.Events, label))
  }
  // sample new value
  newV = getAccumulativeSample(distribution)
  // update model & esd
  sampler.model.invcount_histogram[target] += newV
  esd.V[target] = newV
  esd.Pi = computePi(esd.V)
  esd.Label = computeZ(esd.Tao, esd.Pi)
}