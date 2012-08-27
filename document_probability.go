package scriptModeling

// import "fmt"
import "math"

func (sampler *Sampler) documentLikelihood(events [][]string, label []int) float64 {
  // compute log document likelihood for each eventtype in the model given current labeling
  var labelIdx, tokensGeneratedByEventtype, eventLength int
  var realized bool
  var dcm, update, likelihood, wordEventFactor, wordNormalize float64
  var event []string
  // *range over eventtypes j'*
  for k:=0 ; k<numTop ; k++ {
    eventLength = 0
    dcm = 0.0
    tokensGeneratedByEventtype = 0
    
    // *check whether eventtype is realized to identify its label*
    realized = false
    for idx,e := range(label) {
      if k==e {
	realized = true
	labelIdx=idx
      }
    }

    if !realized {
      // *iterate over terms in vocabulary*
      for _, histogram := range(sampler.model.word_eventtype_histogram) {
	wordEventFactor,_ = math.Lgamma(float64(histogram[k])+sampler.lmPrior)
	tokensGeneratedByEventtype+=histogram[k]
	dcm += wordEventFactor
      }
    } else {
      event = events[labelIdx]
      eventLength = len(event)
      
      // *iterate over terms in vocabulary*
      for term, histogram := range(sampler.model.word_eventtype_histogram) {
	update = 0.0
	// *check whether term is realized*
	for token:=0 ; token<len(event) && update!=1.0; token++ {
	  if term==event[token] {update = 1.0}
	}
	wordEventFactor,_ = math.Lgamma(float64(histogram[k])+sampler.lmPrior+update)
	tokensGeneratedByEventtype+=histogram[k]
	dcm += wordEventFactor
      }
    }
    wordNormalize,_ = math.Lgamma(float64(tokensGeneratedByEventtype) + float64(sampler.model.vocabSize)*sampler.lmPrior + float64(eventLength))
    likelihood += (dcm-wordNormalize)
  }
  return likelihood
}