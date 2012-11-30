 package scriptModeling
 
 import "math/rand"
 import "math"
 import "fmt"
 import "sort"
 
 
 // check whether there are any participants in the esd
 func (esd *ESD) hasParticipants() bool {
   for _,event := range(esd.Label) {
     if len(event.Participants) > 0 {
       return true
     }
   }
   return false
 }
 
 
 // Pick a random participant type to resample from the esd labeling
 func Pick_participant(label Label) int {
   events := make([]int, numTop*numPar)
   var idx int
   
   for eID, event := range(label) {
     if len(event.Participants) > 0 {
       events[idx]=eID
       idx++
     }
   }
   events=events[:idx]
   sort.Ints(events)
   target := rand.Intn(len(events))
//    fmt.Println(events, events[target])
   return events[target]
 }
 
 
 func (sampler *Sampler) Resample_p(esd *ESD, eventID int) {
   var pIdx int
   var newV int
   participants := make([]int, len(esd.Label[eventID].Participants))
   for idx,val := range(esd.Label[eventID].Tau) {
     if val==1 {
       participants[pIdx]=idx
       pIdx++
     }
   }
   participants = participants[:pIdx]
   fmt.Println("participants", participants)
   for _, pID := range(participants) {
     modelLikelihoods, docLikelihoods, alts, tempESDs, /*oldV*/ _ := sampler.getDistributionsP(*esd, pID, eventID)
     //get final distribution
     distribution:=make([]float64, len(modelLikelihoods))
     for idx,_ := range(modelLikelihoods) {
       distribution[idx]=math.Exp(math.Log(modelLikelihoods[idx])+math.Log(docLikelihoods[idx]))
     }
     newV = sample(normalized(distribution))
     fmt.Println("SPEC", normalized(distribution))
     fmt.Println(modelLikelihoods, docLikelihoods, "\n\n\n")

     if newV == -1 {
       esd.Print()
     } 
     // resample eta
//      newWordLabels := esd.compareToP(tempESDs[newV]) 
//      oldWordLabels := tempESDs[newV].compareToP(*esd)
//      if  len(newWordLabels) > 0 {
//        sampler.updateEta(newWordLabels, math.Log(docLikelihoods[newV]), "participant")
//        sampler.updateEta(oldWordLabels, math.Log(docLikelihoods[oldV]), "participant")
//      }
     //update esd and model
     *esd = tempESDs[newV]
     sampler.Model.Participanttype_histogram[alts[newV]]++
     sampler.Model.Participanttype_eventtype_histogram[alts[newV]][eventID]++
     sampler.Model.UpdateParticipantWordCounts(alts[newV], esd.Label[eventID].Participants[alts[newV]], 1)
   }
 }
 
 
 func (sampler *Sampler) getDistributionsP(esd ESD, target, eventID int) (distribution, docLikelihoods []float64, alts []int, tempESDs []ESD, oldV int){
     eIdx := 0
     distribution = make([]float64, numPar)
     docLikelihoods = make([]float64, numPar)
     fullPosterior := make([][5]float64, numPar)
     tempESDs = make([]ESD, numPar)
     alts = make([]int, numPar)
     // Decrement Counts
     sampler.Model.Participanttype_histogram[target]--
     sampler.Model.Participanttype_eventtype_histogram[target][eventID]--
     sampler.Model.UpdateParticipantWordCounts(target, esd.Label[eventID].Participants[target], -1)
     if sampler.Model.Participanttype_histogram[target]<0 {
       panic("Negative Participant Count")
     }
     if sampler.Model.Participanttype_eventtype_histogram[target][eventID]<0 {
       panic("Negative Event Participant Count in resample_p")
     }
     // Compute likelihood for every alternative
     for idx, val := range(esd.Label[eventID].Tau) {
       // relabel ESD
       if val==0 || idx==target {
	 tempESD := esd.Copy()
	 if val==0 {
	   tempESD.flipp(target, idx, eventID)
	   tempESD.UpdateLabelingP(eventID, target, idx)
	 } else {oldV=eIdx}
	 tempESDs[eIdx]=tempESD
	 alts[eIdx]=idx
	 // update Model
	 sampler.Model.Participanttype_histogram[idx]++
	 sampler.Model.UpdateParticipantWordCounts(idx, tempESD.Label[eventID].Participants[idx], 1)
	 // compute score
	 if sampler.Model.Eventtype_histogram[eventID]>0 {
	   sampler.Model.Participanttype_eventtype_histogram[idx][eventID]++
	   distribution[eIdx] = sampler.participantLikelihood(tempESD.Label)
	   sampler.Model.Participanttype_eventtype_histogram[idx][eventID]--
	 }
	 for eID, event := range(tempESD.Label) {
	   for pID, _ := range(event.Participants) {
	     docLikelihoods[eIdx] += sampler.documentLikelihoodP(eID, pID, tempESD.Label)
	   }
	 }
	 // de-update Model
	  sampler.Model.Participanttype_eventtype_histogram[idx][eventID]++
	  fullPosterior[eIdx] = sampler.FullPosterior(tempESD)
	  sampler.Model.Participanttype_eventtype_histogram[idx][eventID]--
	 
	 sampler.Model.Participanttype_histogram[idx]--
	 sampler.Model.UpdateParticipantWordCounts(idx, tempESD.Label[eventID].Participants[idx], -1)
	 eIdx++
       }
     }
     fmt.Println("FULL", normalizeFullPosterior(fullPosterior[:eIdx]))
     return expNormalized(distribution[:eIdx]), expNormalized(docLikelihoods[:eIdx]), alts[:eIdx], tempESDs[:eIdx], oldV
   }
   
func (sampler *Sampler) participantLikelihood(esd Label) (score float64) {
  for eventType, event := range(esd) {
    for pType, _ := range(event.Participants) {
      score += sampler.updateComponentP(pType, eventType)
    }
  }
  return
}
 
 func (sampler *Sampler) updateComponentP(participantID, eventID int) (lgamma float64) {
   // for each alternative participanttype
   var pPositive, pNegative, pNormalize float64
   for i:=0 ; i<numPar ; i++ {
     pPositive,_ = math.Lgamma(float64(sampler.Model.Participanttype_eventtype_histogram[i][eventID]) + sampler.participantPrior)
     pNegative,_ = math.Lgamma(float64(sampler.Model.Eventtype_histogram[eventID]-sampler.Model.Participanttype_eventtype_histogram[i][eventID]) + sampler.participantPrior)
     pNormalize,_ = math.Lgamma(float64(sampler.Model.Eventtype_histogram[eventID]) + 2*sampler.participantPrior)
     lgamma += ((pPositive+pNegative)-pNormalize)
   }
   return
 }