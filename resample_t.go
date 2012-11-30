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
  return el
}




func (sampler *Sampler) Resample_t(esd *ESD, target int) {
  var newLabel int  
  eventLikelihoods, participantLikelihoods, docLikelihoods, alts, tempESDs, /*oldLabel*/ _ := sampler.getDistributionsE(esd, target)  
  // get final distribution
  distribution := make([]float64, len(eventLikelihoods))
  for idx,_ := range(eventLikelihoods) {
    distribution[idx]=math.Exp(math.Log(eventLikelihoods[idx])+math.Log(participantLikelihoods[idx])+math.Log(docLikelihoods[idx]))
  }
  newLabel = sample(normalized(distribution))
  fmt.Println("SPEC", normalized(distribution))
  fmt.Println(eventLikelihoods, participantLikelihoods, docLikelihoods)
  if newLabel == -1 {
       esd.Print()
  } 
  
  // check whether words have changed class; if so: resample eta
//   newWordLabels := esd.compareTo(tempESDs[newLabel])
//   oldWordLabels := tempESDs[newLabel].compareTo(*esd)
//   if  len(newWordLabels) > 0 {
//      sampler.updateEta(newWordLabels, math.Log(docLikelihoods[newLabel]), "event")
//      sampler.updateEta(oldWordLabels, math.Log(docLikelihoods[oldLabel]), "event")
//   }
  // update model & esd
   *esd = tempESDs[newLabel]
   sampler.Model.Eventtype_histogram[alts[newLabel]]++
   sampler.Model.UpdateEventWordCounts(esd.Label, 1)
   sampler.Model.UpdateEventParticipantCounts(esd.Label, 1)
   
   fmt.Println("\n\n\n\n")
}


func (sampler *Sampler) getDistributionsE(esd *ESD, target int) (_, _, _ []float64, _ []int, _ []ESD, oldLabel int) {
  eventDist := make([]float64, numTop)
  participantDist := make([]float64, numTop)
  docLikelihoods := make([]float64, numTop)
  fullPosterior := make([][5]float64, numTop)
  tempESDs := make([]ESD, numTop)
  alts := make([]int, numTop)
  eIdx :=0
  // decrement counts
  sampler.Model.Eventtype_histogram[target]--
  sampler.Model.UpdateEventWordCounts(esd.Label, -1)
  sampler.Model.UpdateEventParticipantCounts(esd.Label, -1)
  if sampler.Model.Eventtype_histogram[target]<0 {
    panic("Negative Event Count in resample_t")
  }
  for tIdx, val := range(esd.Tau) {
    // update ESD Labeling
    if val==0 || tIdx==target {
      tempESD := *esd
      if val ==0 {
	tempESD.flip(target, tIdx)
	tempESD.UpdateLabelingT()
      } else {oldLabel = eIdx}
      tempESDs[eIdx]=tempESD
      alts[eIdx]=tIdx
      // update Model
      sampler.Model.Eventtype_histogram[tIdx]++
      sampler.Model.UpdateEventParticipantCounts(tempESD.Label, 1)
      sampler.Model.UpdateEventWordCounts(tempESD.Label, 1)
      // compute Scores
      eventDist[eIdx] = sampler.eventLikelihood(tempESD.Label)
      participantDist[eIdx] = sampler.participantLikelihood(tempESD.Label)
      docLikelihoods[eIdx] = sampler.documentLikelihood(tempESD.Label)
      
      fullPosterior[eIdx] = sampler.FullPosterior(tempESD)
      
      // de-update Model
      sampler.Model.UpdateEventWordCounts(tempESD.Label, -1)
      sampler.Model.UpdateEventParticipantCounts(tempESD.Label, -1)
      sampler.Model.Eventtype_histogram[tIdx]--
      eIdx++
    }
  }
  fmt.Println("FULL", normalizeFullPosterior(fullPosterior[:eIdx]))
  fmt.Println("Part-log", eventDist[:eIdx], participantDist[:eIdx], docLikelihoods[:eIdx])
  return expNormalized(eventDist[:eIdx]), expNormalized(participantDist[:eIdx]), expNormalized(docLikelihoods[:eIdx]), alts[:eIdx], tempESDs[:eIdx], oldLabel
}


func (sampler *Sampler) eventLikelihood(esd Label) (score float64) {
  for eventType, _ := range(esd) {
    score += sampler.updateComponentE(eventType)
  }
  return
}

func (sampler *Sampler) updateComponentE(targetEvent int) (lgamma float64) {
  var tPositive, tNegative, tNormalize float64
  for k:=0 ; k<numTop ; k++ {
    tPositive,_ = math.Lgamma(float64(sampler.Model.Eventtype_histogram[k])+sampler.eventPrior)
    tNegative,_ = math.Lgamma(float64((sampler.Model.NumESDs)-sampler.Model.Eventtype_histogram[k])+sampler.eventPrior)
    tNormalize,_ = math.Lgamma(float64(sampler.Model.NumESDs)+2*sampler.eventPrior)
    lgamma += ((tPositive+tNegative)-tNormalize)
  }
  return lgamma
}