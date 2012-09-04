 package scriptModeling
// 
import "fmt"
import "math/rand"
import "math"
// 
func pick_event(tau [numTop]int) int {
  //randomly select the event we want to resample
  var el int
  for alt:=0 ; alt!=1; {
    el = rand.Intn(len(tau))
    alt = tau[el]
  }
  fmt.Println("Resampling t=", tau , " for eventtype", el)
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
  var update, lgamma, documentLikelihood, docPositive, docNegative, docNormalize float64
  var newLabel int
  // decrement counts for current target event, and all words in ESD
  sampler.Model.eventtype_histogram[target]--
  sampler.Model.UpdateEventWordCounts(esd.Label, -1)
  sampler.Model.UpdateEventParticipantCountsAll(esd.Label, -1)
  // compute switch-likelihood
  distribution := make([]float64, numTop)
  tempESDs := make([]ESD, numTop)
  alts := make([]int, numTop)
  eIdx :=0
  for tIdx, val := range(esd.Tau) {
    if val==0 || tIdx==target {
      tempESD := *esd
      if val ==0 {
	tempESD.flipEvent(target, tIdx)
	tempESD.UpdateLabelingT()
      }
      lgamma = 0.0
      for k:=0 ; k<numTop ; k++ {
	update=0.0
	if k==tIdx {update=1.0}
	docPositive,_ = math.Lgamma(float64(sampler.Model.eventtype_histogram[k])+sampler.eventPosPrior+update)
	docNegative,_ = math.Lgamma(float64(sampler.Model.numESDs-sampler.Model.eventtype_histogram[k])+sampler.eventNegPrior-update)
	docNormalize,_ = math.Lgamma(float64(sampler.Model.numESDs)+sampler.eventPosPrior+sampler.eventNegPrior+update)
	lgamma += ((docPositive+docNegative)-docNormalize)
      }
      documentLikelihood = sampler.documentLikelihood(tempESD.Label)
      distribution[eIdx]=lgamma+documentLikelihood
      tempESDs[eIdx]=tempESD
      alts[eIdx]=tIdx
      eIdx++
    }
  }
  distribution=distribution[:eIdx]
  tempESDs=tempESDs[:eIdx]
  alts=alts[:eIdx]
  // sample new label
  newLabel = getAccumulativeSample(distribution)
  // update model & esd
   *esd = tempESDs[newLabel]
   sampler.Model.eventtype_histogram[alts[newLabel]]++
   sampler.Model.UpdateEventWordCounts(esd.Label, 1)
   sampler.Model.UpdateEventParticipantCountsAll(esd.Label, 1)
}

















// func (sampler *Sampler) Resample_t(esd *ESD, target int) {
//   var update, lgamma, documentLikelihood, docPositive, docNegative, docNormalize float64
//   var newLabel int
//   var alternatives []int
//   alternatives = newTargets(esd.Tau, target)
//   proposedLabels := make([]Label, len(alternatives))
//   // decrement counts for current target event, and all words in ESD
//   sampler.Model.eventtype_histogram[target]--
//   sampler.Model.UpdateEventWordCounts(esd.Label, -1)
//   sampler.Model.UpdateEventParticipantCountsAll(esd.Label, -1)
//   // compute switch-likelihood
//   distribution := make([]float64, len(alternatives))
//   for idx,alternative := range(alternatives) {
//     proposedLabels[idx] = updateLabelingT(target, alternative, esd.Label)
//     fmt.Println("**", target, "**", proposedLabels[idx])
//     lgamma = 0.0
//     for k:=0 ; k<numTop ; k++ {
//       update=0.0
//       if k==alternative {update=1.0}
//       docPositive,_ = math.Lgamma(float64(sampler.Model.eventtype_histogram[k])+sampler.eventPosPrior+update)
//       docNegative,_ = math.Lgamma(float64(sampler.Model.numESDs-sampler.Model.eventtype_histogram[k])+sampler.eventNegPrior-update)
//       docNormalize,_ = math.Lgamma(float64(sampler.Model.numESDs)+sampler.eventPosPrior+sampler.eventNegPrior+update)
//       lgamma += ((docPositive+docNegative)-docNormalize)
//     }
//     documentLikelihood = sampler.documentLikelihood(proposedLabels[idx])
//     distribution[idx]=lgamma+documentLikelihood
//   }
//   // sample new label
//   newLabel = getAccumulativeSample(distribution)
//   // update model & esd
//    esd.flipEvent(target, alternatives[newLabel])
//    esd.UpdateLabelingT(target, alternatives[newLabel])
//    sampler.Model.eventtype_histogram[alternatives[newLabel]]++
//    sampler.Model.UpdateEventWordCounts(esd.Label, 1)
//    sampler.Model.UpdateEventParticipantCountsAll(esd.Label, 1)
// }
