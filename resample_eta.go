package scriptModeling

// import "fmt"
import "sliceSampler"

func (sampler *Sampler)Resample_eta(eta []float64, i int, docLikelihood float64) (newEta float64) {
   newEta = sliceSampler.SampleEta(3, 5.0, eta[i], false, eta, i, sampler.vocabulary.Covariances.Matrix, docLikelihood)  
  return
}


func (sampler *Sampler) updateEta(diff map[int][]int, wordLikelihood float64,  mode string) {
  // new class
  for class,words := range(diff) {
    if mode == "event" {
      for _,target := range(words) {
	for _,word := range(sampler.vocabulary.Equivalences[target]) {
	  sampler.eventEtas[class][word] = sampler.Resample_eta(sampler.eventEtas[class], word, wordLikelihood)
	}
      }
    // old class   
    } else {
      for _,target := range(words) {
	for _,word := range(sampler.vocabulary.Equivalences[target]) {
	  sampler.participantEtas[class][word] = sampler.Resample_eta(sampler.participantEtas[class], word, wordLikelihood)
	}
      }
    }
    sampler.updatePrior(class, mode)
  }
}
