 package scriptModeling

// import "fmt"
import "math"

func (sampler *Sampler) documentLikelihood(label Label) float64 {
  var wordTypeFactor, wordFactor, wordNorm, documentLikelihood float64
  var typeWordTotal, update, totalUpdate int
  // iterate over eventtypes
  for k := 0 ; k<numTop ; k++ {
    wordFactor = 0.0
    typeWordTotal = 0
    // iterate over terms in event-vocab
    for term, histogram := range(sampler.Model.word_eventtype_histogram) {
      typeWordTotal += histogram[k]
      update = 0
      // check if eventtype is realized as term and set 'update' accordingly
      for eID,_ := range(label) {
	if eID==k {
	  update = computeDelta(term, label[eID].Words)
	}
      }
      // compute LGamma(N(word,event) + prior + udpate)
      wordTypeFactor,_ = math.Lgamma(float64(histogram[k])+sampler.eventlmPrior+float64(update))
      wordFactor += wordTypeFactor
    }
    // normalize
    wordNorm,_ = math.Lgamma(float64(typeWordTotal) + float64(sampler.Model.eventVocabulary)*sampler.eventlmPrior + float64(len(label[k].Words)))
    documentLikelihood += (wordFactor-wordNorm)
  }
  // iterate over participanttypes
  for i:=0 ; i<numPar ; i++ {
     wordFactor = 0.0
     typeWordTotal = 0
     totalUpdate = 0
     // iterate over terms in participant vocab
     for term, histogram := range(sampler.Model.word_participanttype_histogram) {
       typeWordTotal += histogram[i]
       update = 0
       // check if participanttype is realized as term and set 'update' accordingly
       for eID, event := range(label) {
	 for pID,_ := range(event.Participants) {
	   if pID == i {
	     update = computeDelta(term, label[eID].Participants[pID])
	     totalUpdate += update
	   }
	 }
       }
       // compute LGamma(N(word,part) + prior + update)
       wordTypeFactor,_ = math.Lgamma(float64(histogram[i])+sampler.participantlmPrior+float64(update))
       wordFactor += wordTypeFactor
     }
     // normalize
     wordNorm,_ := math.Lgamma(float64(typeWordTotal) + float64(sampler.Model.participantVocabulary)*sampler.participantlmPrior + float64(totalUpdate))
     documentLikelihood += (wordFactor - wordNorm)
  }
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