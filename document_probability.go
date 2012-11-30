 package scriptModeling
 
//  import "fmt"
 import "math" 
 
 // compute document likelihood of the events in the current esd
 // all participant labelings will stay constant -> no need to compute them!
 func (sampler *Sampler) documentLikelihood(label Label) float64 {
   var wordTypeFactor, wordFactor, wordNorm float64
   var typeWordTotal int
   documentLikelihood := 0.0
   // iterate over eventtypes
     for k := 0 ; k<numTop ; k++ {
       wordFactor = 0.0
       typeWordTotal = 0
       // iterate over terms in event-vocab
       for term, histogram := range(sampler.Model.Word_eventtype_histogram) {
	 typeWordTotal += histogram[k]
	 // compute LGamma(N(word,event) + prior + udpate)
	 wordTypeFactor,_ = math.Lgamma(float64(histogram[k])+sampler.EventlmPriors[k][term])
	 wordFactor += wordTypeFactor
       }
       // normalize LGamma(N(words_by_event) + V*prior + total_update)
       wordNorm,_ = math.Lgamma(float64(typeWordTotal) + sum(sampler.EventlmPriors[k]))
       documentLikelihood += (wordFactor - wordNorm)
     }
   return documentLikelihood
 }

 
  // compute document likelihood of the participant realization in question, given the proposed label
  // all event doc likelihoods will stay constant w.r.t. change -> no need to compute them!
  func (sampler *Sampler) documentLikelihoodP(event int, participant int, label Label) float64 {
   var wordTypeFactor, wordFactor, wordNorm float64
   var typeWordTotal/*, update*/ int
   documentLikelihood := 0.0
     // iterate over participanttypes
     for i:= 0 ; i<numPar ; i++ {
      wordFactor = 0.0
      typeWordTotal = 0
      // iterate over terms in participant vocab
      for term, histogram := range(sampler.Model.Word_participanttype_histogram) {
        typeWordTotal += histogram[i]
        // set 'update' according to the number of times term is present in current particip descr
        // compute LGamma(N(word,part) + prior + update)
        wordTypeFactor,_ = math.Lgamma(float64(histogram[i])+sampler.ParticipantlmPriors[i][term])
        wordFactor += wordTypeFactor
      }
      // normalize
      wordNorm,_ = math.Lgamma(float64(typeWordTotal) + sum(sampler.ParticipantlmPriors[i]))
      documentLikelihood += (wordFactor - wordNorm)
     }
   return documentLikelihood
 }
