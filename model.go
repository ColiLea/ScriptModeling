package scriptModeling

// import "fmt"
// import "strconv"
// import "math"

const vocabsize float64 = 5.0

type Model struct {
// Model type, Contains all priors; eventhistogram=counts of events; eventInvcounthistogram:eventspecific inversioncounts; wordEventhistogram:wordspecific eventcounts
  eventtype_histogram Histogram
  word_eventtype_histogram map[string]Histogram
//   eventtype_invCount_histogram map[int]Histogram
  invcount_histogram Histogram
  numESDs int
  vocabSize int
}

func NewModel(topics int) *Model {
  model := new(Model)
  model.eventtype_histogram = newHistogram(topics)
  model.invcount_histogram = newHistogram(topics-1)
//   model.eventtype_invCount_histogram = make(map[int]Histogram, topics-2)//newEventInvHist(corpus,topics)
  model.word_eventtype_histogram = make(map[string]Histogram, topics-1)//newWordEventHist(corpus, topics)
  return model
}

func CreateModel (corpus *Corpus, topics int) *Model {
  model := NewModel(topics)
  model.Initialize(corpus, topics)
  return model
}

func (model *Model) Initialize(corpus *Corpus, topics int) {
  //Initialize histograms from corpus
  model.numESDs = len(*corpus)
  for _,esd := range(*corpus) {
    for _,label := range(esd.Label) {
      model.eventtype_histogram[label]++
    }
    for event,icount := range esd.V {
      if event<numTop {model.invcount_histogram[event]+=icount}
//       if _,ok := model.eventtype_invCount_histogram[event]; !ok {
// 	model.eventtype_invCount_histogram[event]=newHistogram(topics)
//       }
//       model.eventtype_invCount_histogram[event][icount]++
    }
    for event,words := range esd.Events {
      for _,word := range(words) {
	if _,ok := model.word_eventtype_histogram[word]; !ok {
	  model.word_eventtype_histogram[word]=newHistogram(topics)
	}
	model.word_eventtype_histogram[word][esd.Label[event]]++
      }
    }
  }
  model.vocabSize = len(model.word_eventtype_histogram)
}

func (model *Model) IncrementWordCount(event int, word string, count int) {
  model.word_eventtype_histogram[word][event]+=count
}

func (model *Model) IncrementEventCount(event int, count int) {
   model.eventtype_histogram[event]+=count
}

func (model *Model) ReassignCounts(newEventtype int, newAssignments []int, esd [][]string){
  model.IncrementEventCount(newEventtype, 1)
  for idx,event := range(esd) {
    for _, word := range(event) {
      model.IncrementWordCount(newAssignments[idx], word, 1)
    }
  }
}

func (model *Model) DecrementCounts(target int, label []int, esd [][]string){
  model.IncrementEventCount(target, -1)
  for idx,event := range(esd) {
    for _, word := range(event) {
      model.IncrementWordCount(label[idx], word, -1)
    }
  }
}