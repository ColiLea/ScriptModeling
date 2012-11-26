package scriptModeling

import "fmt"

// noodle/telepnone T:30 P:40
// toy.xml P:15 T:15

const numTop int = 10
const numPar int = 15

// Model type, Contains all priors; eventhistogram=counts of events; eventInvcounthistogram:eventspecific inversioncounts; wordEventhistogram:wordspecific eventcounts
type Model struct {
  Eventtype_histogram Histogram
  Participanttype_histogram Histogram
  Participanttype_eventtype_histogram map[int]Histogram	//key=participanttype
  Word_eventtype_histogram map[int]Histogram		//key=word mapped to int
  Word_participanttype_histogram map[int]Histogram 	//key=word mapped to int
  Invcount_histogram Histogram
  NumESDs int
  NumEvents int
  EventVocabulary int
  ParticipantVocabulary int
  Rho []float64
}

func NewModel() *Model {
  model := new(Model)
  model.Eventtype_histogram = newHistogram(numTop)
  model.Participanttype_histogram = newHistogram(numPar)
  model.Invcount_histogram = newHistogram(numTop-1)
  model.Rho = make([]float64, numTop-1)
  model.Word_eventtype_histogram = make(map[int]Histogram)
  model.Word_participanttype_histogram = make(map[int]Histogram)
  model.Participanttype_eventtype_histogram = make(map[int]Histogram, numPar)
  return model
}


func CreateModel (corpus *Corpus) *Model {
  model := NewModel()
  model.Initialize(corpus)
  return model
}

func (model *Model) Initialize(corpus *Corpus) {
  for i:=0 ; i<numPar ; i++ {
    model.Participanttype_eventtype_histogram[i]=newHistogram(numTop)
  }
  //Initialize histograms from corpus
  model.NumESDs = len(*corpus)
  // *initialize eventtype counts*
  for _,esd := range(*corpus) {
    model.NumEvents += len(esd.EventLabel)
    // get eventtypes & words from keys
    for eID,event := range(esd.Label) {
      model.Eventtype_histogram[eID]++
      for _,term := range(event.Words) {
	if _,ok := model.Word_eventtype_histogram[term]; !ok {
	  model.Word_eventtype_histogram[term]=newHistogram(numTop)
	}
	model.Word_eventtype_histogram[term][eID]++
      }
      // *initialize participanttype & word counts from event->map*
      for pID, words := range(event.Participants) {
	model.Participanttype_histogram[pID]++
	model.Participanttype_eventtype_histogram[pID][eID]++
	for _,term := range(words) {
	  if _,ok := model.Word_participanttype_histogram[term]; !ok {
	    model.Word_participanttype_histogram[term]=newHistogram(numPar)
	  }
	  model.Word_participanttype_histogram[term][pID]++
	}
      }
    }
    // *initialize inversion counts*
    for event,icount := range esd.V {
      if event<numTop {model.Invcount_histogram[event]+=icount}
    }
  }
  for i:=0 ; i<numTop-1 ; i++ {
    model.Rho[i]=0.0
  }
  model.EventVocabulary = len(model.Word_eventtype_histogram)
  model.ParticipantVocabulary = len(model.Word_participanttype_histogram)
}

func (model *Model) UpdateEventWordCounts(label Label, count int) {
   for eID, event := range(label) {
     for _, word := range(event.Words) {
       model.Word_eventtype_histogram[word][eID]+=count
       if model.Word_eventtype_histogram[word][eID] < 0 {
	 panic("Negative EventWord Count")
       }
     }
   }
}

func (model *Model) UpdateParticipantWordCounts(target int, words []int, count int) {
  for _, word := range(words) {
    model.Word_participanttype_histogram[word][target]+=count
     if model.Word_participanttype_histogram[word][target] < 0 {
	panic("Negative Participant Count")
     }
  }
}

func (model *Model) UpdateEventParticipantCounts(label Label, count int) {
  for eID, event := range(label) {
    for pID,_ := range(event.Participants) {
      model.Participanttype_eventtype_histogram[pID][eID]+=count
        if model.Participanttype_eventtype_histogram[pID][eID] < 0 {
	 panic("Negative EventParticipant Count")
       }
    }
  }
}


func (model Model) Print() {
  fmt.Println("Eventtype Hist: ", model.Eventtype_histogram)
  fmt.Println("Participa Hist: ", model.Participanttype_histogram)
  fmt.Println("Rho           : ", model.Rho)
  fmt.Println("Inversion Cnts: ", model.Invcount_histogram)
  fmt.Println("Event Vocabula: ", model.EventVocabulary)
  fmt.Println("Parti Vocabula: ", model.ParticipantVocabulary, "\n")
  fmt.Println("Event-Par Hist:")
  for p,hist := range(model.Participanttype_eventtype_histogram) {
    fmt.Println(p, hist)
  }
  fmt.Println("\n", "Word-Event Hist:")
  for wd, e := range(model.Word_eventtype_histogram) {
    fmt.Println(vocabulary.Dictionary.Itov[wd], e)
  }
  fmt.Println("\n", "Word-Participant Hist:")
  for wd, e := range(model.Word_participanttype_histogram) {
    fmt.Println(vocabulary.Dictionary.Itov[wd], e)
  }
}
