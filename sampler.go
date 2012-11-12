 package scriptModeling

import "math"
import "math/rand"
import "strings"
import "leaMatrix"
import "fmt"

type Sampler struct {
  eventPosPrior float64
  eventNegPrior float64
  participantPosPrior float64
  participantNegPrior float64
  EventEtas [][]float64
  ParticipantEtas [][]float64
  EventlmPriors [][]float64
  ParticipantlmPriors [][]float64
  covariances leaMatrix.Matrix
  nu_0 float64
  v_0 [numTop-1]float64
  Model Model
}


func NewSampler(ePprior float64, eNprior float64, pPprior float64, pNprior float64, rho0 float64, nu0 float64, model Model, cov string) *Sampler {
  
  covarianceFlag := strings.Split(cov , " ")
  
  sampler := new(Sampler)
  sampler.Model = model
  sampler.eventPosPrior = ePprior
  sampler.eventNegPrior = eNprior
  sampler.participantPosPrior = pPprior
  sampler.participantNegPrior = pNprior
  
  if covarianceFlag[0]=="load" {
    sampler.covariances = leaMatrix.LoadCovariance(covarianceFlag[1])
  } else {
    fmt.Println(covarianceFlag, len(vocabulary.VList), vocabulary)
    sampler.covariances = *GetCovarianceMatrix(vocabulary.VList, covarianceFlag[1])
  }
  fmt.Println(sampler.covariances.DataStr)  
  sampler.EventEtas, sampler.EventlmPriors  = sampler.initializeEta(numTop)
  sampler.ParticipantEtas, sampler.ParticipantlmPriors = sampler.initializeEta(numPar)
  sampler.nu_0 = nu0*float64(sampler.Model.numESDs)
  sampler.v_0 = vPrior(rho0)
  sampler.Resample_rho()
  return sampler
}


//   select which random variable to resample; 0:p  1:t  2:v  3:rho
func (sampler *Sampler)PickVariable(esd *ESD) {	
  rr := rand.Intn(3)
  if rr <=0 && esd.hasParticipants() {
    sampler.Resample_p(esd, Pick_participant(esd.Label))
  } else if rr<=1 && len(esd.Label) < numTop {
    sampler.Resample_t(esd, pick_event(esd.Tau))
  } else if rr<=2 {
    sampler.Resample_v(esd)
  } else {
    sampler.Resample_rho()
  }
}

func vPrior (rho0 float64) [numTop-1]float64 {
  var vPrior [numTop-1]float64
  for j:=0 ; j<numTop-1 ; j++ {
    vPrior[j] = (1.0/(math.Exp(rho0)-1.0))-((float64(numTop)-float64(j)+1.0)/(math.Exp((float64(numTop)-float64(j)+1.0)*rho0)-1.0))
  }
  return vPrior
}

func (sampler *Sampler)initializeEta(classes int) (eta, prior [][]float64) {
  eta = make([][]float64, classes)
  prior = make([][]float64, classes)
  for classIdx,_ := range(eta){
    eta[classIdx] = make([]float64, len(vocabulary.VList))
    prior[classIdx] = make([]float64, len(vocabulary.VList))
    for wordIdx, _ := range(eta[classIdx]) {
      prior[classIdx][wordIdx] = 1.0/float64(len(vocabulary.VList))
    }
  }
  return eta, prior
}


func (sampler *Sampler)updatePrior(class int, mode string) {
  var normalizer float64
  if mode == "event" {
    normalizer = expSum(sampler.EventEtas[class])
    for wordIdx, _ := range(sampler.EventlmPriors[class]) {
      sampler.EventlmPriors[class][wordIdx] = math.Exp(sampler.EventEtas[class][wordIdx])/normalizer
    }
  } else {
    normalizer = expSum(sampler.ParticipantEtas[class])
    for wordIdx, _ := range(sampler.ParticipantlmPriors[class]) {
      sampler.ParticipantlmPriors[class][wordIdx] = math.Exp(sampler.ParticipantEtas[class][wordIdx])/normalizer
    }
  }
}