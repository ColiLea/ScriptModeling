 package scriptModeling
 
//  import "fmt"
 import "math"
 
 func (sampler *Sampler) documentLikelihood(label Label) float64 {
   // compute document likelihood of the events in the current esd
   // all participant labelings will stay constant -> no need to compute them!
   var wordTypeFactor, wordFactor, wordNorm float64
   var typeWordTotal, update int
   documentLikelihood := 1.0
   // iterate over eventtypes
     for k := 0 ; k<numTop ; k++ {
       wordFactor = 1.0
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
// 	 wordTypeFactor,_ = math.Lgamma(float64(histogram[k])+sampler.eventlmPrior+float64(update))
// 	 wordFactor += wordTypeFactor
	 wordTypeFactor = math.Gamma(float64(histogram[k])+sampler.eventlmPrior+float64(update))
	 wordFactor *= wordTypeFactor
// 	 fmt.Println(">", k, term, histogram[k]+update)
       }

       // normalize
//        wordNorm,_ = math.Lgamma(float64(typeWordTotal) + float64(sampler.Model.eventVocabulary)*sampler.eventlmPrior + float64(len(label[k].Words)))
//        documentLikelihood += (wordFactor - wordNorm)
       wordNorm = math.Gamma(float64(typeWordTotal) + float64(sampler.Model.eventVocabulary)*sampler.eventlmPrior + float64(len(label[k].Words)))
       documentLikelihood *= (wordFactor / wordNorm)
//        fmt.Println(">>>", typeWordTotal+len(label[k].Words), documentLikelihood, "\n")
     }
   return documentLikelihood
 }

  func (sampler *Sampler) documentLikelihoodP(event int, participant int, label Label) float64 {
  // compute document likelihood of the participant realization in question, given the proposed label
  // other participants & all event doc likelihoods will stay constant w.r.t. change -> no need to compute them!
   var wordTypeFactor, wordFactor, wordNorm float64
   var typeWordTotal, update, totalUpdate int
   documentLikelihood := 1.0
     wordFactor = 1.0
     typeWordTotal = 0
     totalUpdate = 0
     // iterate over terms in participant vocab
     for term, histogram := range(sampler.Model.word_participanttype_histogram) {
       typeWordTotal += histogram[participant]
       update = 0
       // set 'update' according to the number of times term is present in current particip descr
       update = computeDelta(term, label[event].Participants[participant])
       totalUpdate += update
       // compute LGamma(N(word,part) + prior + update)
//        wordTypeFactor,_ = math.Lgamma(float64(histogram[participant])+sampler.participantlmPrior+float64(update))
//        wordFactor += wordTypeFactor
       wordTypeFactor = math.Gamma(float64(histogram[participant])+sampler.participantlmPrior+float64(update))
       wordFactor *= wordTypeFactor
//        fmt.Println(">", histogram[participant]+update)
     }
     // normalize
//      wordNorm,_ = math.Lgamma(float64(typeWordTotal) + float64(sampler.Model.participantVocabulary)*sampler.participantlmPrior + float64(totalUpdate))
//      documentLikelihood += (wordFactor - wordNorm)
     wordNorm = math.Gamma(float64(typeWordTotal) + float64(sampler.Model.participantVocabulary)*sampler.participantlmPrior + float64(totalUpdate))
     documentLikelihood *= (wordFactor / wordNorm)
//      fmt.Println(">>>", typeWordTotal+totalUpdate, documentLikelihood, "\n")
   return documentLikelihood
 }

 
 func computeDelta(term string, words []string) (update int) {
   for _,word := range(words) {
     if word==term {
       update++
     }
   }
   return
}