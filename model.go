 package scriptModeling
// 
import "fmt"
// 
const vocabsize float64 = 5.0
const numTop int = 7
const numPar int = 5

type Model struct {
// Model type, Contains all priors; eventhistogram=counts of events; eventInvcounthistogram:eventspecific inversioncounts; wordEventhistogram:wordspecific eventcounts
  eventtype_histogram Histogram
  participanttype_histogram Histogram
  participanttype_eventtype_histogram map[int]Histogram	//key=participanttype
  word_eventtype_histogram map[string]Histogram		//key=word
  word_participanttype_histogram map[string]Histogram 	//key=word
  invcount_histogram Histogram
  numESDs int
  eventVocabulary int
  participantVocabulary int
}

func NewModel() *Model {
  model := new(Model)
  model.eventtype_histogram = newHistogram(numTop)
  model.participanttype_histogram = newHistogram(numPar)
  model.invcount_histogram = newHistogram(numTop-1)
  model.word_eventtype_histogram = make(map[string]Histogram)
  model.word_participanttype_histogram = make(map[string]Histogram)
  model.participanttype_eventtype_histogram = make(map[int]Histogram, numPar)
  return model
}

func CreateModel (corpus *Corpus, topics int) *Model {
  model := NewModel()
  model.Initialize(corpus)
  return model
}

func (model *Model) Initialize(corpus *Corpus) {
  for i:=0 ; i<numPar ; i++ {
    model.participanttype_eventtype_histogram[i]=newHistogram(numTop)
  }
  //Initialize histograms from corpus
  model.numESDs = len(*corpus)
  // *initialize eventtype counts*
  for _,esd := range(*corpus) {
    // get eventtypes & words from keys
    for eID,event := range(esd.Label) {
      model.eventtype_histogram[eID]++
      for _,term := range(event.Words) {
	if _,ok := model.word_eventtype_histogram[term]; !ok {
	  model.word_eventtype_histogram[term]=newHistogram(numTop)
	}
	model.word_eventtype_histogram[term][eID]++
      }
      // *initialize participanttype & word counts from event->map*
      for pID, words := range(event.Participants) {
	model.participanttype_histogram[pID]++
	model.participanttype_eventtype_histogram[pID][eID]++
	for _,term := range(words) {
	  if _,ok := model.word_participanttype_histogram[term]; !ok {
	    model.word_participanttype_histogram[term]=newHistogram(numTop)
	  }
	  model.word_participanttype_histogram[term][pID]++
	}
      }
    }
    // *initialize inversion counts*
    for event,icount := range esd.V {
      if event<numTop {model.invcount_histogram[event]+=icount}
    }
  }
  model.eventVocabulary = len(model.word_eventtype_histogram)
  model.participantVocabulary = len(model.word_participanttype_histogram)
}



func (model *Model) UpdateEventWordCounts(label Label, count int) {
   for eID, event := range(label) {
     for _, word := range(event.Words) {
       model.word_eventtype_histogram[word][eID]+=count
     }
   }
}

func (model *Model) UpdateParticipantWordCounts(target int, words []string, count int) {
  for _, word := range(words) {
    model.word_participanttype_histogram[word][target]+=count
  }
}

func (model *Model) UpdateEventParticipantCountsAll(label Label, count int) {
  for eID, event := range(label) {
    for pID,_ := range(event.Participants) {
      model.participanttype_eventtype_histogram[pID][eID]+=count
    }
  }
}


func (model Model) Print() {
  fmt.Println("Eventtype Hist: ", model.eventtype_histogram)
  fmt.Println("Participa Hist: ", model.participanttype_histogram)
  fmt.Println("Inversion Cnts: ", model.invcount_histogram)
  fmt.Println("Event Vocabula: ", model.eventVocabulary)
  fmt.Println("Parti Vocabula: ", model.participantVocabulary, "\n")
  fmt.Println("Event-Par Hist:")
  for p,hist := range(model.participanttype_eventtype_histogram) {
    fmt.Println(p, hist)
  }
  fmt.Println("\n", "Word-Event Hist:")
  for wd, e := range(model.word_eventtype_histogram) {
    fmt.Println(wd, e)
  }
  fmt.Println("\n", "Word-Participant Hist:")
  for wd, e := range(model.word_participanttype_histogram) {
    fmt.Println(wd, e)
  }
}