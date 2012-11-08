 package scriptModeling
 
//  import "fmt"
 import "math"
 
 
 // compute document likelihood of the events in the current esd
 // all participant labelings will stay constant -> no need to compute them!
 func (sampler *Sampler) documentLikelihood(label Label) float64 {
   var wordTypeFactor, wordFactor, wordNorm, priorDenominator float64
   var typeWordTotal, update int
   documentLikelihood := 0.0
   // iterate over eventtypes
     for k := 0 ; k<numTop ; k++ {
       priorDenominator = priorExpSum(sampler.EventlmPriors[k])
       wordFactor = 0.0
       typeWordTotal = 0
       // iterate over terms in event-vocab
       for term, histogram := range(sampler.Model.word_eventtype_histogram) {
	 typeWordTotal += histogram[k]
	 update = 0
	 // check if eventtype is realized as term and set 'update' accordingly
	 if _,ok := label[k]; ok {
	   update = computeDelta(term, label[k].Words)
	 }
	 // compute LGamma(N(word,event) + prior + udpate)
	 wordTypeFactor,_ = math.Lgamma(float64(histogram[k])+(math.Exp(sampler.EventlmPriors[k][term])/priorDenominator)+float64(update))
	 wordFactor += wordTypeFactor
	 }
       // normalize LGamma(N(words_by_event) + V*prior + total_update)
       wordNorm,_ = math.Lgamma(float64(typeWordTotal) + priorSum(sampler.EventlmPriors[k], priorDenominator) + float64(len(label[k].Words)))
       documentLikelihood += (wordFactor - wordNorm)
     }
   return documentLikelihood
 }

 
  // compute document likelihood of the participant realization in question, given the proposed label
  // all event doc likelihoods will stay constant w.r.t. change -> no need to compute them!
  func (sampler *Sampler) documentLikelihoodP(event int, participant int, label Label) float64 {
   var wordTypeFactor, wordFactor, wordNorm, priorDenominator float64
   var typeWordTotal, update int
   documentLikelihood := 0.0
     // iterate over participanttypes
     for i:= 0 ; i<numPar ; i++ {
      priorDenominator = priorExpSum(sampler.ParticipantlmPriors[i])
      wordFactor = 0.0
      typeWordTotal = 0
      // iterate over terms in participant vocab
      for term, histogram := range(sampler.Model.word_participanttype_histogram) {
	update = 0
        typeWordTotal += histogram[i]
        // set 'update' according to the number of times term is present in current particip descr
        if i==participant {
	  update = computeDelta(term, label[event].Participants[participant])
	}
        // compute LGamma(N(word,part) + prior + update)
        wordTypeFactor,_ = math.Lgamma(float64(histogram[i])+(math.Exp(sampler.ParticipantlmPriors[i][term])/priorDenominator)+float64(update))
        wordFactor += wordTypeFactor
      }
      // normalize
      wordNorm,_ = math.Lgamma(float64(typeWordTotal) + priorSum(sampler.ParticipantlmPriors[i], priorDenominator) + float64(len(label[event].Participants[participant])))
      documentLikelihood += (wordFactor - wordNorm)
     }
   return documentLikelihood
 }

 
func computeDelta(term int, words []int) (update int) {
   for _,word := range(words) {
     if word==term {
       update++
     }
   }
   return
}

func priorSum(priors []float64, norm float64) float64 {
  sum := 0.0
  for _, value := range(priors) {
    sum += (math.Exp(value)/norm)
  }
  return sum
}

func priorExpSum(priors []float64) float64 {
  sum := 0.0
  for _, value := range(priors) {
    sum += math.Exp(value)
  }
  return sum
}