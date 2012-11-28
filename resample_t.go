 package scriptModeling
 
// import "fmt"
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
  var newLabel/*, oldLabel*/ int
  
  eventLikelihoods, participantLikelihoods, docLikelihoods, alts, tempESDs, /*oldLabel*/ _ := sampler.getDistributionsE(esd, target)
  
  // get final distribution
  distribution := make([]float64, len(eventLikelihoods))
  for idx,_ := range(eventLikelihoods) {
    distribution[idx]=math.Log(eventLikelihoods[idx])+math.Log(participantLikelihoods[idx])+math.Log(docLikelihoods[idx])
  }
  distMax, distTotal := computeNorm(distribution)
  for idx,_ := range(distribution) {
    distribution[idx] = math.Exp(distribution[idx]-distMax)/distTotal
  }
  
  // sample new label
  newLabel = sample(distribution)
  
//   fmt.Println("event", eventLikelihoods)
//   fmt.Println("ptcpt", participantLikelihoods)
//   fmt.Println("words", docLikelihoods)
  
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
  
   esd.Print()
  
   *esd = tempESDs[newLabel]
   sampler.Model.Eventtype_histogram[alts[newLabel]]++
   sampler.Model.UpdateEventWordCounts(esd.Label, 1)
   sampler.Model.UpdateEventParticipantCounts(esd.Label, 1)
   
//    fmt.Println("T=", vocabulary.Dictionary.Itov[esd.Label[alts[newLabel]].Words[0]], target, alts[newLabel], "(with event", esd.Label[alts[newLabel]].Participants, ")")
//    sampler.Model.Print()
//    fmt.Println("\n\n")
   
//    sampler.GetFullPosterior(tempESDs, alts)
//    
//    fmt.Println("=======================================\n\n\n\n")
}


func (sampler *Sampler) getDistributionsE(esd *ESD, target int) (_, _, _ []float64, _ []int, _ []ESD, oldLabel int) {
  
  var totalgamma, totaldoclikelihood, totalptcpt, tmax, dmax, pmax float64
  
  // decrement counts for current target event, and ALL words in ESD, and ALL event-participant counts
  // ALL, since event order might change due to fixed v (ordering)
  sampler.Model.Eventtype_histogram[target]--
  sampler.Model.UpdateEventWordCounts(esd.Label, -1)
  sampler.Model.UpdateEventParticipantCounts(esd.Label, -1)
  if sampler.Model.Eventtype_histogram[target]<0 {
    panic("Negative Event Count in resample_t")
  }
  
  
  // compute switch-likelihood
  eventDist := make([]float64, numTop)
  participantDist := make([]float64, numTop)
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
      } else {
	oldLabel = eIdx
      }
      
      sampler.Model.Eventtype_histogram[tIdx]++
      sampler.Model.UpdateEventParticipantCounts(tempESD.Label, 1)
      
      for ee,vv := range(tempESD.Label) {
	eventDist[eIdx] += sampler.updateComponentE(tIdx)
	for pID,_ := range(vv.Participants) {
	  participantDist[eIdx] += sampler.updateComponentP(pID, ee)	  
	}
      }
      sampler.Model.UpdateEventParticipantCounts(tempESD.Label, -1)
      sampler.Model.Eventtype_histogram[tIdx]--
      
      docLikelihoods[eIdx] = sampler.documentLikelihood(tempESD.Label)
      tempESDs[eIdx]=tempESD
      alts[eIdx]=tIdx
      eIdx++
    }
  }
  
  eventDist=eventDist[:eIdx]
  participantDist=participantDist[:eIdx]
  docLikelihoods=docLikelihoods[:eIdx]
  tempESDs=tempESDs[:eIdx]
  alts=alts[:eIdx]
  
  tmax, totalgamma = computeNorm(eventDist)
  dmax, totaldoclikelihood = computeNorm(docLikelihoods)
  pmax, totalptcpt = computeNorm(participantDist)
  
  for idx,_ := range(eventDist) {
    docLikelihoods[idx] = math.Exp(docLikelihoods[idx]-dmax)/totaldoclikelihood
    eventDist[idx] = math.Exp(eventDist[idx]-tmax)/totalgamma
    participantDist[idx] = math.Exp(participantDist[idx]-pmax)/totalptcpt
  }
  return eventDist, participantDist, docLikelihoods, alts, tempESDs, oldLabel
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