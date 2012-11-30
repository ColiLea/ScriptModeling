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
   var newV int
   oldESD := esd.Copy()
   oldDL = sampler.documentLikelihood(oldESD.Label)
   for target, _ := range(esd.V) {
     fmt.Println("V for", target)
     gmmLikelihoods, participantLikelihoods, docLikelihoods := sampler.getDistributionV(oldESD, *esd, target)
     // get final distribution
     distribution := make([]float64, len(gmmLikelihoods))
     for idx, _ := range(gmmLikelihoods) {
       distribution[idx] = math.Exp(math.Log(gmmLikelihoods[idx]) + math.Log(participantLikelihoods[idx]) + math.Log(docLikelihoods[idx]))
     }
     // sample new value
     newV = sample(normalized(distribution))
     fmt.Println("SPEC", normalized(distribution))
     fmt.Println(participantLikelihoods, docLikelihoods, gmmLikelihoods, "\n\n\n")
    // update esd
     esd.V[target] = newV
     esd.Pi = computePi(esd.V)
     esd.UpdateLabelingT()
     // update model
     sampler.Model.Invcount_histogram[target] += esd.V[target]
     sampler.Model.UpdateEventWordCounts(esd.Label, 1)
     sampler.Model.UpdateEventParticipantCounts(esd.Label, 1)
     
     // resample eta
//      if isIn(target, oldESD.EventLabel) {
//        // check whether words have changed class; if so: resample eta
//        newWordLabels := oldESD.compareTo(*esd)
//        oldWordLabels := esd.compareTo(oldESD)
//        if  len(newWordLabels) > 0 {
// 	 sampler.updateEta(newWordLabels, math.Log(docLikelihoods[newV]), "event")
//          sampler.updateEta(oldWordLabels, math.Log(oldDL), "event")
//        }
//      }
   }
   sampler.Resample_rho()
 }
 
 func (sampler *Sampler) getDistributionV(oldESD, esd ESD, target int) (distribution, participantLikelihoods, docLikelihoods []float64){
   proposedLabels := make([]Label, numTop-target)
   distribution = make([]float64, numTop-target)
   docLikelihoods = make([]float64, numTop-target)
   participantLikelihoods = make([]float64, numTop-target)
   fullPosterior := make([][5]float64, numTop-target)
   proposedV:=esd.V
   // Decrement counts
   sampler.Model.UpdateEventWordCounts(esd.Label, -1)
   sampler.Model.UpdateEventParticipantCounts(esd.Label, -1)
   sampler.Model.Invcount_histogram[target] -= esd.V[target]
   if sampler.Model.Invcount_histogram[target]<0 {
     panic("Negative Inversion Count")
   }
   // Compute Scores
   for k:=0 ; k<numTop-target ; k++ {
     // Relabel ESD
     proposedV[target]=k 
     proposedLabels[k] = UpdateLabelingV(esd.Tau, computePi(proposedV), esd.EventLabel, esd.Label)
     if target==0 && proposedV[target]==esd.V[target] {
       oldLabelIdx = k
     }
     distribution[k] = -sampler.Model.Rho[target] * float64(k) 
     participantLikelihoods[k]=1.0
     docLikelihoods[k]=1.0
     if isIn(target, oldESD.EventLabel) {
       // Update Model
       sampler.Model.UpdateEventParticipantCounts(proposedLabels[k], 1)
       sampler.Model.UpdateEventWordCounts(proposedLabels[k], 1)
       // Compute remaining Scores
       participantLikelihoods[k] = sampler.participantLikelihood(proposedLabels[k])
       docLikelihoods[k] = sampler.documentLikelihood(proposedLabels[k])
       
       fullPosterior[k] = sampler.FullPosterior(ESD{proposedLabels[k], 2, [numTop]int{}, [numTop-1]int{}, [numTop]int{}, []int{}})
       
       // De-update Model
       sampler.Model.UpdateEventParticipantCounts(proposedLabels[k], -1)
       sampler.Model.UpdateEventWordCounts(proposedLabels[k], -1)
     }
   }
   for idx,_ := range(distribution) {
     if target ==0 && idx == oldLabelIdx {
       oldDL = docLikelihoods[idx]
     }
   }
   fmt.Println("FULL", normalizeFullPosterior(fullPosterior))
   return expNormalized(distribution), expNormalized(participantLikelihoods), expNormalized(docLikelihoods)
 }