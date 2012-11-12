 package scriptModeling

import "fmt"
import "math/rand"
import "math"


//randomly select the event we want to resample
func pick_event(tau [numTop]int) int {
  var el int
  for alt:=0 ; alt!=1; {
    el = rand.Intn(len(tau))
    alt = tau[el]
  }
//   fmt.Println("\n\nResampling t=", tau , " for eventtype", el)
  return el
}

func newTargets(tau [numTop]int, target int) []int {
  newTargets := make([]int, numTop)
  newTargets[0] = target
  idx := 1
  for eventtype,realized := range(tau) {
    if realized == 0 {
      newTargets[idx]=eventtype
      idx++
    }
  }
  return newTargets[:idx]
}


func (sampler *Sampler) Resample_t(esd *ESD, target int) {
  var update, lgamma, totalgamma, totaldoclikelihood, distTotal, documentLikelihood, docPositive, docNegative, docNormalize, tmax, dmax, distMax float64
  var newLabel int
  // decrement counts for current target event, and ALL words in ESD, and ALL event-participant counts
  // ALL, since event order might change due to fixed v (ordering)
  sampler.Model.eventtype_histogram[target]--
  sampler.Model.UpdateEventWordCounts(esd.Label, -1)
  sampler.Model.UpdateEventParticipantCounts(esd.Label, -1)
  if sampler.Model.eventtype_histogram[target]<0 {
    panic("Negative Event Count in resample_t")
  }
  // compute switch-likelihood
  distribution := make([]float64, numTop)
  docLikelihoods := make([]float64, numTop)
  tempESDs := make([]ESD, numTop)
  alts := make([]int, numTop)
  eIdx :=0
  
  // iterate over eventtypes
  for tIdx, val := range(esd.Tau) {
    // if eventtypes not realized or eventtype==target
    if val==0 || tIdx==target {
      tempESD := *esd
      // flip if not target
      if val ==0 {
	tempESD.flip(target, tIdx)
	tempESD.UpdateLabelingT()
      }
      lgamma = 0.0
      for k:=0 ; k<numTop ; k++ {
	//do p(e=k|...) = p(model|flip) * p(document|flip)
	update=0.0
	if k==tIdx {update=1.0}
	// compute P(model|flip)
	docPositive,_ = math.Lgamma(float64(sampler.Model.eventtype_histogram[k])+sampler.eventPosPrior+update)
	docNegative,_ = math.Lgamma(float64(sampler.Model.numESDs-sampler.Model.eventtype_histogram[k])+sampler.eventNegPrior-update)
	docNormalize,_ = math.Lgamma(float64(sampler.Model.numESDs)+sampler.eventPosPrior+sampler.eventNegPrior)
	lgamma += ((docPositive+docNegative)-docNormalize)
      }
      // compute P(document|flip)  <<== THIS SHOULD BE INSIDE THE LOOP & THE LOOP INSIDE THE DOCLIKELIHOOD SHOULD BE GONE!?!?!
      documentLikelihood = sampler.documentLikelihood(tempESD.Label)
      distribution[eIdx]=lgamma
      docLikelihoods[eIdx]=documentLikelihood

      tempESDs[eIdx]=tempESD
      alts[eIdx]=tIdx
      eIdx++
    }
  }
  
  distribution=distribution[:eIdx]
  
  docLikelihoods=docLikelihoods[:eIdx]
  tempESDs=tempESDs[:eIdx]
  alts=alts[:eIdx]
  
  tmax, totalgamma = computeNorm(distribution)
  dmax, totaldoclikelihood = computeNorm(docLikelihoods)
  
  for idx,_ := range(distribution) {
    docLikelihoods[idx]=math.Exp(docLikelihoods[idx]-dmax)/totaldoclikelihood
    distribution[idx] = math.Log(math.Exp(distribution[idx]-tmax)/totalgamma) + math.Log(docLikelihoods[idx])
  }
  distMax, distTotal = computeNorm(distribution)
  for idx,_ := range(distribution) {
    distribution[idx] = math.Exp(distribution[idx]-distMax)/distTotal
  }
  // sample new label
  newLabel = sample(distribution)
  // check whether words have changed class; if so: resample eta
  diff := esd.compareTo(tempESDs[newLabel])
  diff2 := tempESDs[newLabel].compareTo(*esd)
  if  len(diff) > 0 {
    fmt.Println("Event", diff, diff2)
    for class,words := range(diff) {
      for _,word := range(words) {
	fmt.Println(sampler.EventlmPriors[class][word])
	sampler.EventEtas[class][word] = sampler.Resample_eta(sampler.EventEtas[class], word, sampler.wordLikelihood(class, "event"))
	sampler.updatePrior(class, "event")
	fmt.Println(sampler.EventlmPriors[class][word], "\n---------\n")
      }
    }
    fmt.Println("diff2 (should DECREASE): ")
    for class,words := range(diff2) {
      for _,word := range(words) {
	fmt.Println(sampler.EventlmPriors[class][word])
	sampler.EventEtas[class][word] = sampler.Resample_eta(sampler.EventEtas[class], word, sampler.wordLikelihood(class, "event"))
	sampler.updatePrior(class, "event")
	fmt.Println(sampler.EventlmPriors[class][word], "\n---------\n")
      }
    }
  }
  // update model & esd
   *esd = tempESDs[newLabel]
   sampler.Model.eventtype_histogram[alts[newLabel]]++
   sampler.Model.UpdateEventWordCounts(esd.Label, 1)
   sampler.Model.UpdateEventParticipantCounts(esd.Label, 1)
}