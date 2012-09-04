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
  sampler.Model = model
  sampler.eventPosPrior = ePprior
  sampler.eventNegPrior = eNprior
  sampler.eventlmPrior = elmprior
  sampler.participantPosPrior = pPprior
  sampler.participantNegPrior = pNprior
  sampler.participantlmPrior = plmprior
  sampler.nu_0 = nu0*float64(sampler.Model.numESDs)
  sampler.v_0 = vPrior(rho0)
  // TODO this is wrong!!
  sampler.rho = sampler.v_0
  // TODO this is wrong!!
  return sampler
}

func (sampler *Sampler)PickVariable(esd *ESD) {
  //select which random variable to resample; 0:t  1:v  2:rho
  rr := rand.Intn(3)
  if rr ==0 && esd.hasParticipants() {
    fmt.Println("Resampling P...\n")
    sampler.Resample_p(esd, Pick_participant(&esd.Label))
  } else if rr==1{
    fmt.Println("Resampling V...\n")
    sampler.Resample_v(esd, pick_invcount(esd.V))
  } else if rr==2{
    fmt.Println("Resampling T...\n")
    sampler.Resample_t(esd, pick_event(esd.Tau))
  } else {
    fmt.Println("Resampling Rho...\n")
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


// func (sampler *Sampler) InitializeRho() [numTop-1]int{
// //    prob = -rho * (totalV + v_0 * nu_0) - (nu_0 + numDocs) * (log1p(-exp(-nminusj * rho)) - log1p(-exp(-rho)));
//   var initRho [numTop-1]ints
//   for i:=0 ; i<numTop-1 ; i++ {
//     
//   }
// }



