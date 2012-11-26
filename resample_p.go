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
   var pIdx, oldV int
   var lgamma, distTotal, totalgamma, documentLikelihood, totaldoclikelihood, update, pPositive, pNegative/*, pTerm, nTerm*/, pNormalize, pmax, dmax, distMax float64
   var distribution []float64
   var docLikelihoods []float64
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

     fmt.Println("\n\nUpdating event ", eventID, " and participant", pID, "= ", sampler.vocabulary.Dictionary.Itov[esd.Label[eventID].Participants[pID][0]])
     
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
     tempESDs := make([]ESD, numPar)
     alts := make([]int, numPar)
     
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
	 // for each alternative participanttype
	 for i:=0 ; i<numPar ; i++ {
	   update = 0.0
	   if i==idx {
	     update = 1.0
	   }
	   pPositive,_ = math.Lgamma(float64(sampler.Model.Participanttype_eventtype_histogram[i][eventID]) + sampler.participantPrior + update)
	   pNegative,_ = math.Lgamma(float64(sampler.Model.Eventtype_histogram[eventID]-sampler.Model.Participanttype_eventtype_histogram[i][eventID]) + sampler.participantPrior - update)
	   pNormalize,_ = math.Lgamma(float64(sampler.Model.Eventtype_histogram[eventID]) + 2.0*sampler.participantPrior)
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

     tmpDist:= make([]float64, len(distribution))
     for idx,_ := range(distribution) {
       docLikelihoods[idx]=math.Exp(docLikelihoods[idx]-dmax)/totaldoclikelihood
       tmpDist[idx]= math.Exp(distribution[idx]-pmax)/totalgamma
       distribution[idx] = math.Log(math.Exp(distribution[idx]-pmax)/totalgamma) + math.Log(docLikelihoods[idx])
     }
     
     distMax,distTotal = computeNorm(distribution)
     for idx,_ := range(distribution) {
       distribution[idx]=math.Exp(distribution[idx]-distMax)/distTotal
     }
     newV = sample(distribution)
     
     if newV == -1 {
       esd.Print()
     }
    
     // resample eta
     diff := esd.compareToP(tempESDs[newV]) 
     diff2 := tempESDs[newV].compareToP(*esd)
     if  len(diff) > 0 {
       sampler.updateEta(diff, math.Log(docLikelihoods[newV]), "participant")
       sampler.updateEta(diff2, math.Log(docLikelihoods[oldV]), "participant")
     }
     //update esd and model
     *esd = tempESDs[newV]
     sampler.Model.Participanttype_histogram[alts[newV]]++
     sampler.Model.Participanttype_eventtype_histogram[alts[newV]][eventID]++
     sampler.Model.UpdateParticipantWordCounts(alts[newV], esd.Label[eventID].Participants[alts[newV]], 1)
//      fmt.Println("NEW")
//      for t:=0 ; t<numPar; t++ {
//        fmt.Println(sampler.Model.Eventtype_histogram[eventID],sampler.Model.Participanttype_eventtype_histogram[t][eventID])
//      }
     
   }
 }
// 	 // for each alternative participanttype
// 	 for i:=0 ; i<numPar ; i++ {
// 	   // for each eventtype
// 	   for e := 0 ; e<numTop ; e++ {
// 	     update = 0.0
// 	     if i==idx && e==eventID {
// 	       update = 1.0
// 	       fmt.Println(i,e,sampler.Model.Participanttype_eventtype_histogram[i][e],sampler.Model.Eventtype_histogram[e]-sampler.Model.Participanttype_eventtype_histogram[i][e],sampler.Model.Eventtype_histogram[e])
// 	     }
// 	     pPositive,_ = math.Lgamma(float64(sampler.Model.Participanttype_eventtype_histogram[i][e]) + sampler.participantPrior + update)
// 	     pNegative,_ = math.Lgamma(float64(sampler.Model.Eventtype_histogram[e]-sampler.Model.Participanttype_eventtype_histogram[i][e]) + sampler.participantPrior-update)
// 	     pNormalize,_ = math.Lgamma(float64(sampler.Model.Eventtype_histogram[e])+2*sampler.participantPrior)
// 	     lgamma += ((pPositive+pNegative)-pNormalize)
// 	   }
// 	 }


/*
	 // for each alternative participanttype
	 for i:=0 ; i<numPar ; i++ {
	   // for each eventtype
	   pPositive = 0
	   pNegative = 0
	   for e := 0 ; e<numTop ; e++ {
	     update = 0.0
	     if i==idx && e==eventID {
	       update = 1.0
	       fmt.Println(i,e,sampler.Model.Participanttype_eventtype_histogram[i][e],sampler.Model.Eventtype_histogram[e]-sampler.Model.Participanttype_eventtype_histogram[i][e],sampler.Model.Eventtype_histogram[e])
	     }
	     pTerm,_ = math.Lgamma(float64(sampler.Model.Participanttype_eventtype_histogram[i][e]) + sampler.participantPrior + update)
	     nTerm,_ = math.Lgamma(float64(sampler.Model.Eventtype_histogram[e]-sampler.Model.Participanttype_eventtype_histogram[i][e]) + sampler.participantPrior-update)
	     pPositive += pTerm
	     pNegative += nTerm
	   }
	     pNormalize,_ = math.Lgamma(float64(sampler.Model.NumEvents) + float64(sampler.Model.NumEvents)*sampler.participantPrior)
	     lgamma += ((pPositive+pNegative)-pNormalize)
	 }*/