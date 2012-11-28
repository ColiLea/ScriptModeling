 package scriptModeling
 
 import "fmt"
 import "math/rand"
 import "math"
 
 var oldLabelIdx int
 var oldDL float64
 
 func pick_invcount(v [numTop-1]int) int {
   newV := rand.Intn(len(v))
   return newV
 }
 
 func (sampler *Sampler) Resample_v(esd *ESD) {
   
//    fmt.Println("Resampling v")

   var newV int
   oldESD := esd.Copy()
   oldDL = sampler.documentLikelihood(oldESD.Label)
   for target, _ := range(esd.V) {
     gmmLikelihoods, participantLikelihoods, docLikelihoods := sampler.getDistributionV(oldESD, *esd, target)
     
     // get final distribution
     distribution := make([]float64, len(gmmLikelihoods))
     for idx, _ := range(gmmLikelihoods) {
       distribution[idx] = math.Log(gmmLikelihoods[idx]) + math.Log(participantLikelihoods[idx]) + math.Log(docLikelihoods[idx])
     }
     distMax, distTotal := computeNorm(distribution)
     for idx,_ := range(distribution) {
       distribution[idx]=math.Exp(distribution[idx]-distMax)/distTotal
     }
     
     // sample new value
     newV = /*max(distribution)*/sample(distribution)
//      fmt.Println("V", newV)
     
//      fmt.Println("\n", "resampling inv count for eventtype", target, "was", esd.V[target], "is", newV)
//      fmt.Println(gmmLikelihoods)
//      fmt.Println(participantLikelihoods)
//      fmt.Println(docLikelihoods)
     
     // update esd
     esd.V[target] = newV
     esd.Pi = computePi(esd.V)
     esd.UpdateLabelingT()
     
     fmt.Println("V", newV)
     
     // update model
     sampler.Model.Invcount_histogram[target] += esd.V[target]
     sampler.Model.UpdateEventWordCounts(esd.Label, 1)
     sampler.Model.UpdateEventParticipantCounts(esd.Label, 1)
     
     // resample eta
     if isIn(target, oldESD.EventLabel) {
       // check whether words have changed class; if so: resample eta
       newWordLabels := oldESD.compareTo(*esd)
       oldWordLabels := esd.compareTo(oldESD)
       if  len(newWordLabels) > 0 {
	 sampler.updateEta(newWordLabels, math.Log(docLikelihoods[newV]), "event")
         sampler.updateEta(oldWordLabels, math.Log(oldDL), "event")
       }
     }
   }
   sampler.Resample_rho()
 }
 
 func (sampler *Sampler) getDistributionV(oldESD, esd ESD, target int) (distribution, participantLikelihoods, docLikelihoods []float64){
   
   var gmm, gmax, totalgmm, dmax, totaldoclikelihood, pmax, totalp float64
   proposedLabels := make([]Label, numTop-target)
   distribution = make([]float64, numTop-target)
   docLikelihoods = make([]float64, numTop-target)
   participantLikelihoods = make([]float64, numTop-target)
   proposedV:=esd.V
   
   // Decrement counts
   sampler.Model.UpdateEventWordCounts(esd.Label, -1)
   sampler.Model.UpdateEventParticipantCounts(esd.Label, -1)
   sampler.Model.Invcount_histogram[target] -= esd.V[target]
   if sampler.Model.Invcount_histogram[target]<0 {
     panic("Negative Inversion Count")
   }
   
   for k:=0 ; k<numTop-target ; k++ {
     
     proposedV[target]=k 
     
     if target==0 && proposedV[target]==esd.V[target] {
       oldLabelIdx = k
     }
     
     // NOTE: I am using 'k+1' below as my topicIDs start with 0 ...shouldn't matter, right???
     gmm = -sampler.Model.Rho[target] * float64(k+1)
     // compute documentLikelihood if eventtype for which inv count is resampled is realized in esd
     proposedLabels[k] = UpdateLabelingV(esd.Tau, computePi(proposedV), esd.EventLabel, esd.Label)
     // updated label is ok
     
     distribution[k] = gmm
     participantLikelihoods[k]=1.0
     docLikelihoods[k]=1.0
     
     if isIn(target, oldESD.EventLabel) {
       sampler.Model.UpdateEventParticipantCounts(proposedLabels[k], 1)
       sampler.Model.UpdateEventWordCounts(proposedLabels[k], 1)
       
       docLikelihoods[k]=sampler.documentLikelihood(proposedLabels[k])
       
       for ee,vv := range(proposedLabels[k]) {
	 for pID,_ := range(vv.Participants) {
	   participantLikelihoods[k] += sampler.updateComponentP(pID, ee)	  
	 }
       }
       
       docLikelihoods[k] = sampler.documentLikelihood(proposedLabels[k])
       sampler.Model.UpdateEventParticipantCounts(proposedLabels[k], -1)
       sampler.Model.UpdateEventWordCounts(proposedLabels[k], -1)
     }
   }
   gmax, totalgmm = computeNorm(distribution)
   dmax, totaldoclikelihood = computeNorm(docLikelihoods)
   pmax, totalp = computeNorm(participantLikelihoods)
   
   for idx,_ := range(distribution) {
     docLikelihoods[idx]=math.Exp(docLikelihoods[idx]-dmax)/totaldoclikelihood
     distribution[idx] = math.Exp(distribution[idx]-gmax)/totalgmm
     participantLikelihoods[idx]=math.Exp(participantLikelihoods[idx]-pmax)/totalp
     
     if target ==0 && idx == oldLabelIdx {
       oldDL = docLikelihoods[idx]
     }
   }
   return
 }