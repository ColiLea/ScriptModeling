package scriptModeling

// import "fmt"

type Model struct {
// Model type, Contains all priors; eventhistogram=counts of events; eventInvcounthistogram:eventspecific inversioncounts; wordEventhistogram:wordspecific eventcounts
  eventPosPrior float64
  eventNegPrior float64
  lmPrior float64
  eventtype_histogram Histogram
  word_eventtype_histogram map[string]Histogram
  eventtype_invCount_histogram map[int]Histogram
}

func CreateModel (corpus *Corpus, topics int) *Model {
  model := new(Model)
  model.eventtype_histogram = newHistogram(topics)//newEventtypeHist(corpus,topics)
  model.eventtype_invCount_histogram = make(map[int]Histogram, topics-2)//newEventInvHist(corpus,topics)
  model.word_eventtype_histogram = make(map[string]Histogram, topics-1)//newWordEventHist(corpus, topics)
  model.Initialize(corpus, topics)
  return model
}

func (model *Model) Initialize(corpus *Corpus, topics int) {
  //Initialize histograms from corpus
  for _,esd := range(*corpus) {
    for _,label := range(esd.Label) {
      model.eventtype_histogram[label]++
    }
    for event,icount := range esd.V {
      if _,ok := model.eventtype_invCount_histogram[event]; !ok {
	model.eventtype_invCount_histogram[event]=newHistogram(topics)
      }
      model.eventtype_invCount_histogram[event][icount]++
    }
    for event,word := range esd.Events {
      if _,ok := model.word_eventtype_histogram[word]; !ok {
	model.word_eventtype_histogram[word]=newHistogram(topics)
      }
      model.word_eventtype_histogram[word][esd.Label[event]]++
    }
  }
}
