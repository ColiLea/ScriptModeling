 package scriptModeling
 
 import "math/rand"
 import "math"
//  import "fmt"
 
 
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
   target := rand.Intn(len(events[:idx]))
//    fmt.Println("\n\nResampling for participants in eventtype", events[target], "=", label[events[target]].Participants)
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
   for _, pID := range(participants) {
     
     modelLikelihoods, docLikelihoods, alts, tempESDs, /*oldV*/ _:= sampler.getDistributionsP(*esd, pID, eventID)
     
     
     //get final distribution
     distribution:=make([]float64, len(modelLikelihoods))
     for idx,_ := range(modelLikelihoods) {
       distribution[idx]=math.Log(modelLikelihoods[idx])+math.Log(docLikelihoods[idx])
     }
     distMax,distTotal := computeNorm(distribution)
     for idx,_ := range(distribution) {
       distribution[idx]=math.Exp(distribution[idx]-distMax)/distTotal
     }
     
     
//      fmt.Println(modelLikelihoods)
//      fmt.Println(docLikelihoods)
     
     newV = sample(distribution)
          
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
     
//      fmt.Println("P = ", vocabulary.Dictionary.Itov[esd.Label[eventID].Participants[alts[newV]][0]] , pID, alts[newV], "in event", eventID )
//      sampler.Model.Print()
//      fmt.Println("\n\n\n")

     
   }
 }
 
 
 func (sampler *Sampler) getDistributionsP(esd ESD, pID, eventID int) (distribution, docLikelihoods []float64, alts []int, tempESDs []ESD, oldV int){

   var lgamma, totalgamma, documentLikelihood, totaldoclikelihood, pmax, dmax float64
     
     target := pID
     eIdx := 0
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
     // Compute likelihood for every type
     distribution = make([]float64, numPar)
     docLikelihoods = make([]float64, numPar)
     tempESDs = make([]ESD, numPar)
     alts = make([]int, numPar)
     
     for idx, val := range(esd.Label[eventID].Tau) {
       if val==0 || idx==target {
	 
	 tempESD := esd.Copy()
	 if val==0 {
	   tempESD.flipp(target, idx, eventID)
	   tempESD.UpdateLabelingP(eventID, target, idx)
	   
	 } else {
	   oldV=eIdx
	 }
	 lgamma = 0.0
	 sampler.Model.Participanttype_histogram[idx]++
	 if sampler.Model.Eventtype_histogram[eventID]>0 {
	   sampler.Model.Participanttype_eventtype_histogram[idx][eventID]++
	   lgamma = sampler.updateComponentP(idx, eventID)
	   sampler.Model.Participanttype_eventtype_histogram[idx][eventID]--
	 }
	 sampler.Model.Participanttype_histogram[idx]--
	 documentLikelihood = sampler.documentLikelihoodP(eventID, idx, tempESD.Label)
	 distribution[eIdx]=lgamma
	 docLikelihoods[eIdx]=documentLikelihood
	 tempESDs[eIdx]=tempESD
	 alts[eIdx]=idx
	 eIdx++
       }
     }
     distribution=distribution[:eIdx]
     docLikelihoods=docLikelihoods[:eIdx]
     tempESDs=tempESDs[:eIdx]
     alts=alts[:eIdx]
     
     pmax, totalgamma = computeNorm(distribution)
     dmax, totaldoclikelihood = computeNorm(docLikelihoods)
     
     tmpDist:= make([]float64, len(distribution))
     for idx,_ := range(distribution) {
       docLikelihoods[idx]=math.Exp(docLikelihoods[idx]-dmax)/totaldoclikelihood
       tmpDist[idx]= math.Exp(distribution[idx]-pmax)/totalgamma
       distribution[idx] =math.Exp(distribution[idx]-pmax)/totalgamma
     }
     
     return distribution, docLikelihoods, alts, tempESDs, oldV
   }
 
 func (sampler *Sampler) updateComponentP(participantID, eventID int) (lgamma float64) {
   // for each alternative participanttype
   var pPositive, pNegative, pNormalize float64
   for i:=0 ; i<numPar ; i++ {
     pPositive,_ = math.Lgamma(float64(sampler.Model.Participanttype_eventtype_histogram[i][eventID]) + sampler.participantPrior)
     pNegative,_ = math.Lgamma(float64(sampler.Model.Eventtype_histogram[eventID]-sampler.Model.Participanttype_eventtype_histogram[i][eventID]) + sampler.participantPrior)
     pNormalize,_ = math.Lgamma(float64(sampler.Model.Eventtype_histogram[eventID]) + 2*sampler.participantPrior)
//      fmt.Println("p", i, "e", eventID, "+" ,sampler.Model.Participanttype_eventtype_histogram[i][eventID], "-", sampler.Model.Eventtype_histogram[eventID]-sampler.Model.Participanttype_eventtype_histogram[i][eventID], sampler.Model.Eventtype_histogram[eventID], "||", pPositive, pNegative, pNormalize, "||", ((pPositive+pNegative)-pNormalize))
     lgamma += ((pPositive+pNegative)-pNormalize)
   }
   return
 }