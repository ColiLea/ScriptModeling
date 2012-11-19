 package scriptModeling
 
 import "fmt"
 import "math/rand"
 import "math"
 
 
 func pick_invcount(v [numTop-1]int) int {
   newV := rand.Intn(len(v))
   return newV
 }
 
 func (sampler *Sampler) Resample_v(esd *ESD) {
   fmt.Println("===================================================================\n\nResampling v=", esd.V)
   var docLikelihoods []float64
   var newV, oldLabelIdx int
   oldESD := esd.Copy()
   oldDL := sampler.documentLikelihood(oldESD.Label)
   for target, _ := range(esd.V) {
     fmt.Println("\n-------------------------------------------------\nResampling for position", target, "so we have ", numTop-target, " possibilities.")
     esd.Print()
     var proposedV [numTop-1]int
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
     docLikelihoods = make([]float64, numTop-target)
     // try every possible value
     for k:=0 ; k<numTop-target ; k++ {
       
       proposedV[target]=k 
       
       if target==0 && proposedV[target]==esd.V[target] {
	 oldLabelIdx = k
       }
       
       // NOTE: I am using 'k+1' below as my topicIDs start with 0 ...shouldn't matter, right???
       gmm = -sampler.Model.rho[target] * float64(k+1)
       // compute documentLikelihood if eventtype for which inv count is resampled is realized in esd
       proposedLabels[k] = UpdateLabelingV(esd.Tau, computePi(proposedV), esd.EventLabel, esd.Label)
       // updated label is ok
       
       docLikelihoods[k]=1.0
       
       if isIn(target, oldESD.EventLabel) {
	documentLikelihood = sampler.documentLikelihood(proposedLabels[k])
	fmt.Println("!!", documentLikelihood, "!!")
	distribution[k] = gmm
	docLikelihoods[k]=documentLikelihood
       }
     }
     
     gmax, totalgmm = computeNorm(distribution)
     dmax, totaldoclikelihood = computeNorm(docLikelihoods)
     
     fmt.Println("(1) DocumentLikelihoods	(2) GMM Likelihoods")
     for idx,_ := range(distribution) {
       docLikelihoods[idx]=math.Exp(docLikelihoods[idx]-dmax)/totaldoclikelihood
       fmt.Println(docLikelihoods[idx], math.Exp(distribution[idx]-gmax)/totalgmm)
       distribution[idx] = math.Log(math.Exp(distribution[idx]-gmax)/totalgmm) + math.Log(docLikelihoods[idx])
       
       if target ==0 && idx == oldLabelIdx {
	 oldDL = docLikelihoods[idx]
       }
     }
     
     distMax, distTotal = computeNorm(distribution)

     for idx,_ := range(distribution) {
       distribution[idx]=math.Exp(distribution[idx]-distMax)/distTotal
     }
     
     // sample new value
     newV = getAccumulativeSample(distribution)
     
     fmt.Println("Final Distribution: \n", distribution, "\nAnd we pick component: ", newV)
     // update esd
     esd.V[target] = newV
     esd.Pi = computePi(esd.V)
     esd.UpdateLabelingT()
     
     esd.Print()
     fmt.Println("--------------------------------------------------------------")
     
     //update model
     sampler.Model.invcount_histogram[target] += esd.V[target]
     sampler.Model.UpdateEventWordCounts(esd.Label, 1)
     sampler.Model.UpdateEventParticipantCounts(esd.Label, 1)
     
   }
   // check whether words have changed class; if so: resample eta
//    fmt.Println(oldDL)
   diff := oldESD.compareTo(*esd)
   diff2 := esd.compareTo(oldESD)
   if  len(diff) > 0 {
     sampler.updateEta(diff, math.Log(docLikelihoods[newV]), diff2, math.Log(oldDL), "event")
   }
 }