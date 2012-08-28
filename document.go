package scriptModeling

import "fmt"

type ESD struct {
// Elemantary-Sequence-Description, consisting of Lengt=number of events; Events=events(words); Tao=event realizations; V=inversions; Pi=global Ordering; Label=event labels
  Events *Events
  Participants *Participants
  Length int
  V [numTop-1]int
  Pi [numTop]int
}

type Events struct {
  Words [][]string
  Tau [numTop]int
  Label []int
}

type Participants struct {
  Words [][]string
  Label [][]int
}

type Corpus []*ESD

func (esd *ESD) Init() {
  esd.ComputePi()
  esd.ComputeZ()
  if len(esd.Events.Words) != esd.Length || len(esd.Events.Words)!=len(esd.Participants.Words) {
    panic("Event- and Participantlist not of same length!")
  }
}

func (esd *ESD) ComputePi() {
// Compute global labeling from V (inversion count vector)
  esd.Pi[0] = numTop-1
  for j:=numTop-2; j>=0; j-- {
   for i:=numTop-2; i>=esd.V[j]; i-- {
     esd.Pi[i+1]=esd.Pi[i]
   }
   esd.Pi[esd.V[j]]=j
  }
}

func (esd *ESD) ComputeZ() {
  fmt.Println(esd)
// Compute the ESD labeling from Tao (realization vector) and Pi (global labeling)
  esd.Events.Label=make([]int,len(esd.Events.Words))
  event:=0
  for _,el := range(esd.Pi) {
    if esd.Events.Tau[el] == 1 {
      esd.Events.Label[event]=el
      event++
    }
  }
}

func (esd *ESD) flipEvent(oldEvent int, newEvent int ) {
  esd.Events.Tau[oldEvent]=0
  esd.Events.Tau[newEvent]=1
}

func (esd *ESD) Print() { 
  for ev,_ := range(esd.Events.Words) {
    fmt.Println("String: ", esd.Events.Words[ev], "  Label: ", esd.Events.Label[ev])
  }
}