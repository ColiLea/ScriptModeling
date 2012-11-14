 package scriptModeling
 
//  import "fmt"
 import "math/rand"
 import "math"
 
 
 func pick_invcount(v [numTop-1]int) int {
   newV := rand.Intn(len(v))
//      fmt.Println("\n\nResampling v=", v , " for eventtype", newV)
   return newV
 }
 
 func (sampler *Sampler) Resample_v(esd *ESD) {
   oldESD := esd.Copy()
   for target, _ := range(esd.V) {
     var proposedV [numTop-1]int
     var newV int
     var documentLikelihood, gmm, distTotal, totalgmm, totaldoclikelihood, gmax, dmax, distMax float64
     proposedLabels := make([]Label, numTop-target)
     sampler.Model.UpdateEventWordCounts(esd.Label, -1)
     sampler.Model.UpdateEventParticipantCounts(esd.Label, -1)
     // decrement global inversion count for target eventtype
     sampler.Model.invcount_histogram[target] -= esd.V[target]
     if sampler.Model.invcount_histogram[target]<0 {
       panic("Negative Inversion Count")
     }
     proposedV = esd.V
     distribution := make([]float64, numTop-target)
     docLikelihoods := make([]float64, numTop-target)
     // try every possible value
     for k:=0 ; k<numTop-target ; k++ {
       proposedV[target]=k 
       // NOTE: I am using 'k+1' below as my topicIDs start with 0 ...shouldn't matter, right???
       gmm = -sampler.Model.rho[target] * float64(k+1)
       // compute documentLikelihood if eventtype for which inv count is resampled is realized in esd
       proposedLabels[k] = UpdateLabelingV(esd.Tau, computePi(proposedV), esd.EventLabel, esd.Label)
       documentLikelihood = sampler.documentLikelihood(proposedLabels[k])
       distribution[k] = gmm
       docLikelihoods[k]=documentLikelihood
     }
     
     gmax, totalgmm = computeNorm(distribution)
     dmax, totaldoclikelihood = computeNorm(docLikelihoods)
     
     for idx,_ := range(distribution) {
       docLikelihoods[idx]=math.Exp(docLikelihoods[idx]-dmax)/totaldoclikelihood
       distribution[idx] = math.Log(math.Exp(distribution[idx]-gmax)/totalgmm) + math.Log(docLikelihoods[idx])
     }
     distMax, distTotal = computeNorm(distribution)

     for idx,_ := range(distribution) {
       distribution[idx]=math.Exp(distribution[idx]-distMax)/distTotal
     }
     // sample new value
     newV = getAccumulativeSample(distribution)
     // update esd
     esd.V[target] = newV
     esd.Pi = computePi(esd.V)
     esd.UpdateLabelingT()
     //update model
     sampler.Model.invcount_histogram[target] += esd.V[target]
     sampler.Model.UpdateEventWordCounts(esd.Label, 1)
     sampler.Model.UpdateEventParticipantCounts(esd.Label, 1)
     // check whether words have changed class; if so: resample eta
   }
   diff := oldESD.compareTo(*esd)
   diff2 := esd.compareTo(oldESD)
   if  len(diff) > 0 {
     sampler.updateEta(diff, diff2, "event")
   }
 }