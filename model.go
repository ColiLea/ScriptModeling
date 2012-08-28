package scriptModeling

import "fmt"

const vocabsize float64 = 5.0
const numTop int = 7
const numPar int = 5

type Model struct {
// Model type, Contains all priors; eventhistogram=counts of events; eventInvcounthistogram:eventspecific inversioncounts; wordEventhistogram:wordspecific eventcounts
  eventtype_histogram Histogram
  participanttype_histogram Histogram
  participanttype_eventtype_histogram map[int]Histogram	//key=participanttype
  word_eventtype_histogram map[string]Histogram		//key=word
  word_participanttype_histogram map[string]Histogram	//key=word
  invcount_histogram Histogram
  numESDs int
  vocabSize int
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
  //Initialize histograms from corpus
  model.numESDs = len(*corpus)
  // *initialize eventtype counts*
  for _,esd := range(*corpus) {
    for eIdx,eLabel := range(esd.Events.Label) {
      model.eventtype_histogram[eLabel]++
      // *initialize participanttype counts*
      for _, pLabel := range(esd.Participants.Label[eIdx]) {
	model.participanttype_histogram[pLabel]++
	if _, ok := model.participanttype_eventtype_histogram[pLabel] ; !ok {
	  model.participanttype_eventtype_histogram[pLabel]=newHistogram(numTop)
	}
	model.participanttype_eventtype_histogram[pLabel][eLabel]++
      }
    }
    // *initialize inversion counts*
    for event,icount := range esd.V {
      if event<numTop {model.invcount_histogram[event]+=icount}
    }
    for event,_ := range esd.Events.Words {
      // *initialize eventtype language models*
      for _,word := range(esd.Events.Words[event]) {
	if _,ok := model.word_eventtype_histogram[word]; !ok {
	  model.word_eventtype_histogram[word]=newHistogram(numTop)
	}
	model.word_eventtype_histogram[word][esd.Events.Label[event]]++
      }
      // *initialize participanttype language models*
      for participant,word := range(esd.Participants.Words[event]){
	if _,ok := model.word_participanttype_histogram[word]; !ok {
	  model.word_participanttype_histogram[word]=newHistogram(numPar)
	}
// 	fmt.Println(participant, word, esd.Participants.Label)
	model.word_participanttype_histogram[word][esd.Participants.Label[event][participant]]++
      }
    }
  }
  model.vocabSize = len(model.word_eventtype_histogram)
}

func (model *Model) IncrementWordCount(label int, word string, count int) {
   model.word_eventtype_histogram[word][label]+=count
}


func (model *Model) IncrementEventWordCount(events [][]string, label[]int, count int) {
  for idx, event := range(events) {
    for _, word := range(event) {
      model.word_eventtype_histogram[word][label[idx]]+=count
    }
  }
}

func (model *Model) IncrementParticipantWordCount(participants []string, target int, count int) {
  for _, word := range(participants) {
    model.word_eventtype_histogram[word][target]+=count
  }
}

func (model *Model) IncrementEventCount(event int, count int) {
   model.eventtype_histogram[event]+=count
}

func (model *Model) IncrementParticipantCount(participant int, count int) {
   model.participanttype_histogram[participant]+=count
}

func (model *Model) IncrementEventParticipantCount(participants [][]int, events []int, count int) {
  for idx, label := range(events) {
    for _, participant := range(participants[idx]) {
      if _, ok := model.participanttype_eventtype_histogram[participant] ; !ok {
	model.participanttype_eventtype_histogram[participant]=newHistogram(numTop)
      }
      model.participanttype_eventtype_histogram[participant][label]+=count
    }
  }
}

func (model *Model) IncrementInversionCount(target int, count int) {
  model.invcount_histogram[target] += count
}

func (model *Model) ReassignCounts(esd ESD, target int, newValue int, mode string){
  if mode=="event" {
    model.IncrementEventCount(newValue, 1)
    model.IncrementEventParticipantCount(esd.Participants.Label, esd.Events.Label, 1)
  } else if mode=="inversion" {
    model.IncrementInversionCount(target, newValue)
  }
  model.IncrementEventWordCount(esd.Events.Words, esd.Events.Label, 1)
}


func (model *Model) DecrementCounts(eTarget int, pTarget int, eventLabel []int, participantLabel[][]int, invCount int, eventDescriptions [][]string, participantDescriptions [][]string, mode string) {
  if mode=="event" {
    model.IncrementEventCount(eTarget, -1)
    model.IncrementEventWordCount(eventDescriptions, eventLabel, -1)
    model.IncrementEventParticipantCount(participantLabel, eventLabel, -1)
  } else if mode == "participant" {
    model.IncrementParticipantCount(pTarget, -1)
    model.IncrementParticipantWordCount(participantDescriptions[eTarget], pTarget, -1)
    model.IncrementEventParticipantCount(participantLabel, eventLabel, -1)
  } else if mode == "inversion" {
    model.IncrementInversionCount(eTarget, -invCount)
    model.IncrementEventWordCount(eventDescriptions, eventLabel, -1)
    model.IncrementEventParticipantCount(participantLabel, eventLabel, -1)
  } else {
    panic("Invalid resampling mode!!")
  }
}

// func (model *Model) DecrementCounts(target int, label []int, esd [][]string){
//   model.IncrementEventCount(target, -1)
//   for idx,event := range(esd) {
//     for _, word := range(event) {
//       model.IncrementWordCount(label[idx], word, -1)
//     }
//   }
// }