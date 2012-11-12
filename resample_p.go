 package scriptModeling
 
 import "math/rand"
 import "math"
 import "fmt"
 
 
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
   participants := make([]int, len(esd.Label[eventID].Participants))
   for idx,val := range(esd.Label[eventID].Tau) {
     if val==1 {
       participants[pIdx]=idx
       pIdx++
     }
   }
   participants = participants[:pIdx]
   for _, pID := range(participants) {
     var lgamma, distTotal, totalgamma, documentLikelihood, totaldoclikelihood, update, pPositive, pNegative, pNormalize, pmax, dmax, distMax float64
     var distribution []float64
     var docLikelihoods []float64
     var newV int
     target := pID
     eIdx := 0
//      fmt.Println("...participant type", target)
     // Decrement Counts
     sampler.Model.participanttype_histogram[target]--
     sampler.Model.participanttype_eventtype_histogram[target][eventID]--
     sampler.Model.UpdateParticipantWordCounts(target, esd.Label[eventID].Participants[target], -1)
     if sampler.Model.participanttype_histogram[target]<0 {
       panic("Negative Participant Count")
     }
     if sampler.Model.participanttype_eventtype_histogram[target][eventID]<0 {
       panic("Negative Event Participant Count in resample_p")
     }
     // Compute likelihood for every type
     distribution = make([]float64, numPar)
     docLikelihoods = make([]float64, numPar)
     tempESDs := make([]ESD, numPar)
     alts := make([]int, numPar)
     
     for idx, val := range(esd.Label[eventID].Tau) {
       if val==0 || idx==target {
	 tempESD := esd.Copy()
	 if val==0 {
	   tempESD.flipp(target, idx, eventID)
	   tempESD.UpdateLabelingP(eventID, target, idx)
	 }
	 lgamma = 0.0
	 for i:=0 ; i<numPar ; i++ {
	   update = 0.0
	   if i==idx {update = 1.0}
	   pPositive,_ = math.Lgamma(float64(sampler.Model.participanttype_eventtype_histogram[i][eventID]) + sampler.participantPosPrior + update)
	   pNegative,_ = math.Lgamma(float64(sampler.Model.participanttype_histogram[i]-sampler.Model.participanttype_eventtype_histogram[i][eventID]) + sampler.participantNegPrior - update)
	   pNormalize,_ = math.Lgamma(float64(sampler.Model.participanttype_histogram[i])+sampler.participantPosPrior+sampler.participantNegPrior)
	   lgamma += ((pPositive+pNegative)-pNormalize)
	 }
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
     
     for idx,_ := range(distribution) {
       docLikelihoods[idx]=math.Exp(docLikelihoods[idx]-dmax)/totaldoclikelihood
       distribution[idx] = math.Log(math.Exp(distribution[idx]-pmax)/totalgamma) + math.Log(docLikelihoods[idx])
     }
     distMax,distTotal = computeNorm(distribution)
     for idx,_ := range(distribution) {
       distribution[idx]=math.Exp(distribution[idx]-distMax)/distTotal
     }
     newV = sample(distribution)
     // resample eta
     diff := esd.compareToP(tempESDs[newV])
     diff2 := tempESDs[newV].compareToP(*esd)
     if  len(diff) > 0 {
      for class,words := range(diff) {
	for _,word := range(words) {
	  fmt.Println(sampler.ParticipantlmPriors[class][word])
	  sampler.ParticipantEtas[class][word] = sampler.Resample_eta(sampler.ParticipantEtas[class], word, sampler.wordLikelihood(class, "participant"))
	  sampler.updatePrior(class, "participant")
	  fmt.Println(sampler.ParticipantlmPriors[class][word], "\n---------\n")
	}
      }
      fmt.Println("diff2 (should DECREASE): ")
      for class,words := range(diff2) {
	for _,word := range(words) {
	  fmt.Println(sampler.ParticipantlmPriors[class][word])
	  sampler.ParticipantEtas[class][word] = sampler.Resample_eta(sampler.ParticipantEtas[class], word, sampler.wordLikelihood(class, "participant"))
	  sampler.updatePrior(class, "participant")
	  fmt.Println(sampler.ParticipantlmPriors[class][word], "\n---------\n")
	}
      }
    }
     //update esd and model
     *esd = tempESDs[newV]
     sampler.Model.participanttype_histogram[alts[newV]]++
     sampler.Model.participanttype_eventtype_histogram[alts[newV]][eventID]++
     sampler.Model.UpdateParticipantWordCounts(alts[newV], esd.Label[eventID].Participants[alts[newV]], 1)
   }
 }
