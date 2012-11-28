 package scriptModeling

import "math"
import "math/rand"
import "fmt"

type Sampler struct {
  eventPrior float64
  participantPrior float64
  eventEtas [][]float64
  participantEtas [][]float64
  EventlmPriors [][]float64
  ParticipantlmPriors [][]float64
  vocabulary Vocabulary
  nu_0 float64
  v_0 [numTop-1]float64
  Model Model
}

func NewSampler(ePrior float64, pPrior float64, rho0 float64, nu0 float64, model Model, scenario string, mode int) *Sampler {
  if mode == 2 {
    vocabulary = LoadVocabulary("/local/lea/thesis/data/vocabularies/"+scenario+".bin")
    fmt.Println("Dict: \n", vocabulary.Dictionary)
    fmt.Println("Matr: \n", vocabulary.Covariances.Matrix)
    fmt.Println("Inve: \n", vocabulary.Covariances.Inverse)
    fmt.Println("Equi: ")
  } else if mode==4 {
    vocabulary.Covariances = new(CovarianceStruct)
    vocabulary.Equivalences = EquivalenceStruct{}
    vocabulary.Dictionary = Dictionary{map[string]int{"cook":0,"add":1,"serve":2,"water":3,"salt":4,"pasta":5}, map[int]string{0:"cook",1:"add",2:"serve",3:"water",4:"salt",5:"pasta"},[]string{"cook","add","serve","water","salt","pasta"}, []string{"v","v","v","n","n","n"}}
    vocabulary.Covariances.Matrix = [][]float64{[]float64{1, 0.0192307692308, 0, 0, 0, 0}, []float64{0.0192307692308, 1, 0, 0, 0, 0}, []float64{0, 0, 1, 0, 0, 0}, []float64{0, 0, 0, 1, 0, 0}, []float64{0, 0, 0, 0, 1, 0}, []float64{0, 0, 0, 0, 0, 1}}
    vocabulary.Covariances.Inverse = [][]float64{[]float64{1.0003699593, -0.0192378838328, 0, 0, 0, 0}, []float64{-0.0192378838328, 1.0003699593, 0, 0, 0, 0}, []float64{0, 0, 1, 0, 0, 0}, []float64{0, 0, 0, 1, 0, 0}, []float64{0, 0, 0, 0, 1, 0}, []float64{0, 0, 0, 0, 0, 1}}
    vocabulary.Equivalences = map[int][]int{0:[]int{0}, 1:[]int{1}, 2:[]int{2}, 3:[]int{3}, 4:[]int{4}, 5:[]int{5}}
  }
  sampler := new(Sampler)
  sampler.Model = model
  sampler.vocabulary = vocabulary
  sampler.vocabulary.Equivalences.Print()
  sampler.eventPrior = ePrior
  sampler.participantPrior = pPrior  
  sampler.eventEtas, sampler.EventlmPriors  = sampler.initializeEta(numTop)
  sampler.participantEtas, sampler.ParticipantlmPriors = sampler.initializeEta(numPar)
  sampler.nu_0 = 0.1/*2*/*float64(sampler.Model.NumESDs)
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
  } else{
    sampler.Resample_v(esd)
  }
}

func vPrior (rho0 float64) [numTop-1]float64 {
  var vPrior [numTop-1]float64
  for j:=0 ; j<numTop-1 ; j++ {
    vPrior[j] = 1.0/rho0
//     vPrior[j] = (1.0/(math.Exp(rho0)-1.0))-((float64(numTop)-float64(j)+1.0)/(math.Exp((float64(numTop)-float64(j)+1.0)*rho0)-1.0))
  }
  return vPrior
}

func (sampler *Sampler)initializeEta(classes int) (eta, prior [][]float64) {
  eta = make([][]float64, classes)
  prior = make([][]float64, classes)
  for classIdx,_ := range(eta){
    eta[classIdx] = make([]float64, len(vocabulary.Dictionary.VList))
    prior[classIdx] = make([]float64, len(vocabulary.Dictionary.VList))
    for wordIdx, _ := range(eta[classIdx]) {
      prior[classIdx][wordIdx] = 1.0/10*(float64(len(vocabulary.Dictionary.VList)))
//       prior[classIdx][wordIdx] = 0.01
    }
  }
  return eta, prior
}


func (sampler *Sampler)updatePrior(class int, mode string) {
  var normalizer, exponent float64
  if mode == "event" {
    normalizer = expSum(sampler.eventEtas[class])
    for wordIdx, _ := range(sampler.EventlmPriors[class]) {
      exponent = math.Exp(sampler.eventEtas[class][wordIdx])
      sampler.EventlmPriors[class][wordIdx] = exponent/normalizer
    }
  } else {
    normalizer = expSum(sampler.participantEtas[class])
    for wordIdx, _ := range(sampler.ParticipantlmPriors[class]) {
      exponent = math.Exp(sampler.participantEtas[class][wordIdx])
      sampler.ParticipantlmPriors[class][wordIdx] = exponent/normalizer
    }
  }
}


// func (sampler *Sampler)FullPosterior(esds []ESD, labels []int) {
//   eventWord := make([]float64, len(esds))
//   ptcptWord := make([]float64, len(esds))
//   eventLike := make([]float64, len(esds))
//   ptcptLike := make([]float64, len(esds))
//   for idx,_ := range(esds) {
//     sampler.Model.Eventtype_histogram[alts[idx]]++
//     sampler.Model.UpdateEventParticipantCounts(esds[idx],1)
//     eventWord[idx] = sampler.updateComponentE(alts[idx])
//     for
//     ptcptWord[idx] += 
//     sampler.Model.Eventtype_histogram[alts[idx]]--
//     sampler.Model.UpdateEventParticipantCounts(esds[idx],-1)
//   }
// }


