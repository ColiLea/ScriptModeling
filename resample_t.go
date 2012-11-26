 package scriptModeling
 
// import "fmt"
import "math/rand"
import "math"


//randomly select the event we want to resample
func pick_event(tau [numTop]int) int {
  var el int
  for alt:=0 ; alt!=1; {
    el = rand.Intn(len(tau))
    alt = tau[el]
  }
//   fmt.Println("\n\nResampling t=", tau , " for eventtype", el)
  return el
}

func (sampler *Sampler) Resample_t(esd *ESD, target int) {
    
  var update, lgamma, totalgamma, totaldoclikelihood, distTotal, documentLikelihood, tPositive, tNegative, tNormalize, tmax, dmax, distMax float64
  var newLabel, oldLabel int
  
  // decrement counts for current target event, and ALL words in ESD, and ALL event-participant counts
  // ALL, since event order might change due to fixed v (ordering)
  sampler.Model.Eventtype_histogram[target]--
  sampler.Model.UpdateEventWordCounts(esd.Label, -1)
  sampler.Model.UpdateEventParticipantCounts(esd.Label, -1)
  if sampler.Model.Eventtype_histogram[target]<0 {
    panic("Negative Event Count in resample_t")
  }
  
  // compute switch-likelihood
  distribution := make([]float64, numTop)
  docLikelihoods := make([]float64, numTop)
  tempESDs := make([]ESD, numTop)
  alts := make([]int, numTop)
  eIdx :=0
  
  // iterate over eventtypes
  for tIdx, val := range(esd.Tau) {
    // if eventtypes not realized or eventtype==target
    if val==0 || tIdx==target {
      tempESD := *esd
      // flip if not target
      if val ==0 {
	tempESD.flip(target, tIdx)
	tempESD.UpdateLabelingT()
      } else {
	oldLabel = eIdx
      }
      lgamma = 0.0
      for k:=0 ; k<numTop ; k++ {
	//do p(e=k|...) = p(model|flip) * p(document|flip)
	update=0.0
	if k==tIdx {update=1.0}
	// compute P(model|flip)
	tPositive,_ = math.Lgamma(float64(sampler.Model.Eventtype_histogram[k])+sampler.eventPrior+update)
	tNegative,_ = math.Lgamma(float64((sampler.Model.NumESDs)-sampler.Model.Eventtype_histogram[k])+sampler.eventPrior-update)
	tNormalize,_ = math.Lgamma(float64(sampler.Model.NumESDs)+2*sampler.eventPrior)
	lgamma += ((tPositive+tNegative)-tNormalize)
      }
      // compute P(document|flip)  <<== THIS SHOULD BE INSIDE THE LOOP & THE LOOP INSIDE THE DOCLIKELIHOOD SHOULD BE GONE!?!?!
      documentLikelihood = sampler.documentLikelihood(tempESD.Label)
      distribution[eIdx]=lgamma
      docLikelihoods[eIdx]=documentLikelihood
      tempESDs[eIdx]=tempESD
      alts[eIdx]=tIdx
      eIdx++
    }
  }
  
  distribution=distribution[:eIdx]
  docLikelihoods=docLikelihoods[:eIdx]
  tempESDs=tempESDs[:eIdx]
  alts=alts[:eIdx]

  tmax, totalgamma = computeNorm(distribution)
  dmax, totaldoclikelihood = computeNorm(docLikelihoods)
  
  
  tmpDist:= make([]float64, len(distribution))
  for idx,_ := range(distribution) {
    docLikelihoods[idx] = math.Exp(docLikelihoods[idx]-dmax)/totaldoclikelihood
    tmpDist[idx] = math.Exp(distribution[idx]-tmax)/totalgamma
    distribution[idx] = math.Log(math.Exp(distribution[idx]-tmax)/totalgamma) + math.Log(docLikelihoods[idx])
  }
  
  distMax, distTotal = computeNorm(distribution)
  for idx,_ := range(distribution) {
    distribution[idx] = math.Exp(distribution[idx]-distMax)/distTotal
  }
  // sample new label
  newLabel = sample(distribution)
  if newLabel == -1 {
       esd.Print()
  }
  /*
     fmt.Println(sampler.vocabulary.Dictionary.Itov[esd.Label[target].Words[0]], target, alts[newLabel])
     fmt.Println(distribution)
     fmt.Println(docLikelihoods)
     fmt.Println(tmpDist)*/
  
//   fmt.Println(oldLabel)
  // check whether words have changed class; if so: resample eta
  diff := esd.compareTo(tempESDs[newLabel])
  diff2 := tempESDs[newLabel].compareTo(*esd)
  
  if  len(diff) > 0 {
     sampler.updateEta(diff, math.Log(docLikelihoods[newLabel]), "event")
     sampler.updateEta(diff2, math.Log(docLikelihoods[oldLabel]), "event")
  }
  // update model & esd
   *esd = tempESDs[newLabel]
   sampler.Model.Eventtype_histogram[alts[newLabel]]++
   sampler.Model.UpdateEventWordCounts(esd.Label, 1)
   sampler.Model.UpdateEventParticipantCounts(esd.Label, 1)
}