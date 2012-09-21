 package scriptModeling

  import "fmt"
  import "math/rand"
  import "math"


func pick_invcount(v [numTop-1]int) int {
  newV := rand.Intn(len(v))
  fmt.Println("Resampling v=", v , " for eventtype", newV)
  return newV
}

func (sampler *Sampler) Resample_v(esd *ESD, target int) {
  var proposedV [numTop-1]int
  var newV int
  var documentLikelihood, gmm, distTotal, totalgmm, totaldoclikelihood float64
  proposedLabels := make([]Label, numTop-target)
  sampler.Model.UpdateEventWordCounts(esd.Label, -1, "v", target)
  sampler.Model.UpdateEventParticipantCountsAll(esd.Label, -1)
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
      totalgmm += math.Exp(gmm)
      totaldoclikelihood += math.Exp(documentLikelihood)
  }
  for idx,_ := range(distribution) {
//     fmt.Println(math.Exp(distribution[idx])/totalgmm , math.Exp(docLikelihoods[idx])/totaldoclikelihood)
    distribution[idx] = math.Log(math.Exp(distribution[idx])/totalgmm) + math.Log(math.Exp(docLikelihoods[idx])/totaldoclikelihood)
    distTotal += math.Exp(distribution[idx])
  }
  for idx,_ := range(distribution) {
    distribution[idx]=math.Exp(distribution[idx])/distTotal
  }
  // sample new value
  newV = getAccumulativeSample(distribution)
  // update model & esd
  esd.V[target] = newV
  fmt.Println(distribution)
  fmt.Println(newV, "  = ordering", esd.V)
  esd.Pi = computePi(esd.V)
  esd.UpdateLabelingT()
  sampler.Model.invcount_histogram[target] += esd.V[target]
  sampler.Model.UpdateEventWordCounts(esd.Label, 1, "v", target)
  sampler.Model.UpdateEventParticipantCountsAll(esd.Label, 1)
}