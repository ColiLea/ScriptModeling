package scriptModeling

// import "fmt"
import "math"

func (sampler *Sampler) documentLikelihood(events [][]string, label []int, participants [][]string, plabel [][]int) float64 {
  ll := sampler.documentLikelihoodEvents(events, label)
  ll += sampler.documentLikelihoodParticipants(participants, plabel)
  return ll
}

func(sampler *Sampler) documentLikelihoodEvents(events [][]string, label []int) float64 {
  // compute log document likelihood for each eventtype in the model given current labeling
  var labelIdx, tokensGeneratedByEventtype, delta int
  var realized bool
  var dcm, update, likelihood, wordEventFactor, wordNormalize float64
  var event []string
  // *range over eventtypes j'*
  for k:=0 ; k<numTop ; k++ {
    delta = 0
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
	wordEventFactor,_ = math.Lgamma(float64(histogram[k])+sampler.eventlmPrior)
	tokensGeneratedByEventtype+=histogram[k]
	dcm += wordEventFactor
      }
    } else {
      event = events[labelIdx]
      delta = len(event)
      
      // *iterate over terms in vocabulary*
      for term, histogram := range(sampler.model.word_eventtype_histogram) {
	update = 0.0
	// *check whether term is realized*
	for token:=0 ; token<len(event) && update!=1.0; token++ {
	  if term==event[token] {update = 1.0}
	}
	wordEventFactor,_ = math.Lgamma(float64(histogram[k])+sampler.eventlmPrior+update)
	tokensGeneratedByEventtype+=histogram[k]
	dcm += wordEventFactor
      }
    }
    wordNormalize,_ = math.Lgamma(float64(tokensGeneratedByEventtype) + float64(sampler.model.vocabSize)*sampler.eventlmPrior + float64(delta))
    likelihood += (dcm-wordNormalize)
  }
  return likelihood
}

func (sampler *Sampler) documentLikelihoodParticipants(participants [][]string, label [][]int) float64 {
  var wordParticipantFactor, wordNormalize, dcm, update, likelihood float64
  var tokensGeneratedByParticipanttype, delta int
  for i:=0 ; i<numPar ; i++ {
    delta = 0
    dcm = 0.0
    tokensGeneratedByParticipanttype = 0
    for term, histogram := range(sampler.model.word_participanttype_histogram) {
      tokensGeneratedByParticipanttype+=histogram[i]
      update=0.0
      for e, events := range(participants) {
	for p, _ := range(events) {
	  if participants[e][p]==term && label[e][p]==i {
	    update+=1.0
	    delta++
	  }
	}
      }
      wordParticipantFactor, _ = math.Lgamma(float64(histogram[i])+sampler.participantlmPrior+update)
      dcm += wordParticipantFactor
    }
    wordNormalize,_ = math.Lgamma(float64(tokensGeneratedByParticipanttype) + float64(len(sampler.model.word_participanttype_histogram)) * sampler.participantlmPrior + float64(delta))
    likelihood += dcm-wordNormalize
  }
  return likelihood
}