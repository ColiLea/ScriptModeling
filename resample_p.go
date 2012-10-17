 package scriptModeling
 
 import "math/rand"
 import "math"
//  import "fmt"
 
 func (esd *ESD) hasParticipants() bool {
   // check whether there are any participants in the esd
   for _,event := range(esd.Label) {
     if len(event.Participants) > 0 {
       return true
     }
   }
   return false
 }
 
 func Pick_participant(label Label) int {
   // Pick a random participant type to resample from the esd labeling
   events := make([]int, numTop*numPar)
   var idx int
   for eID, event := range(label) {
     if len(event.Participants) > 0 {
       events[idx]=eID
       idx++
     }
   }
   target := rand.Intn(len(events[:idx]))
//    fmt.Println("Resampling for participants in eventtype", events[target], "=", label[events[target]].Participants)
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
     if sampler.Model.participanttype_histogram[target]<0 {
       panic("Negative Participant Count")
     }
     sampler.Model.UpdateParticipantWordCounts(target, esd.Label[eventID].Participants[target], -1)
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
       distribution[idx] = math.Log(math.Exp(distribution[idx]-pmax)/totalgamma) + math.Log(math.Exp(docLikelihoods[idx]-dmax)/totaldoclikelihood)
     }
     distMax,distTotal = computeNorm(distribution)
     for idx,_ := range(distribution) {
       distribution[idx]=math.Exp(distribution[idx]-distMax)/distTotal
     }
     newV = sample(distribution)
//      fmt.Println(distribution)
//      fmt.Println(newV, "  = participanttype", alts[newV])
     //update esd and model
     *esd = tempESDs[newV]
     //      esd.UpdateLabelingP(eventID, alternatives[len(alternatives)-1], alternatives[newV])
     sampler.Model.participanttype_histogram[alts[newV]]++
     sampler.Model.participanttype_eventtype_histogram[alts[newV]][eventID]++
     sampler.Model.UpdateParticipantWordCounts(alts[newV], esd.Label[eventID].Participants[alts[newV]], 1)
   }
 }
 
//   func getAlternatives(participant int, label map[int][]string) []int {
//    // Get alternative participant types ; TODO: ugly function!!
//    var add bool
//    idx := 1
//    alts := make([]int, numPar-len(label)+1)
//    alts[0] = participant
//    for ii:=0 ; ii<numPar ; ii++ {
//      add=false
//      if ii!= alts[0] {
//        add=true
//        for pID,_ := range(label) {
// 	 if pID==ii {
// 	   add=false
// 	   break
// 	 }
//        }
//      }
//      if add == true {
//        alts[idx]=ii
//        idx++
//      }
//    }
//    return alts
//  }
