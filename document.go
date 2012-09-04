package scriptModeling

import "fmt"

type ESD struct {
// Elemantary-Sequence-Description
  Label Label
  Length int
  Tau [numTop]int
  V [numTop-1]int
  Pi [numTop]int
  EventLabel []int
}

type Label map[int]Content

type Content struct {
  Words []string
  Participants map[int][]string
}

type Corpus []*ESD

func (esd *ESD) Init() {
  esd.ComputePi()
  esd.ComputeZ()
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
// Compute the ESD labeling from Tao (realization vector) and Pi (global labeling)
  esd.EventLabel=make([]int,len(esd.Label))
  event:=0
  for _,el := range(esd.Pi) {
    if esd.Tau[el] == 1 {
      esd.EventLabel[event]=el
      event++
    }
  }
}

func (esd *ESD) flipEvent(oldEvent int, newEvent int ) {
  esd.Tau[oldEvent]=0
  esd.Tau[newEvent]=1
}


func (esd *ESD) UpdateLabelingT(oldVal int, newVal int) {
    if oldVal != newVal {
      esd.Label[newVal]=esd.Label[oldVal]
      delete(esd.Label, oldVal)
      esd.ComputeZ()
    }
}

func (esd *ESD) UpdateLabelingP(eventIdx int, oldVal int, newVal int) {
    esd.Label[eventIdx].Participants[newVal]=esd.Label[eventIdx].Participants[oldVal]
    delete(esd.Label[eventIdx].Participants, oldVal)
}

func (esd *ESD) UpdateLabelingV() {
  newZ := computeZ(esd.Tau, esd.Pi)
  update:=false
  for idx,_ := range(newZ) {
    if newZ[idx] != esd.EventLabel[idx] {
      update=true
    }
  }
  if update==true {
    contents := make([]Content, len(esd.EventLabel))
    for idx,eID := range(esd.EventLabel) {
      contents[idx]=esd.Label[eID]
    }
    for idx, id := range(newZ) {
      esd.Label[id]=contents[idx]
    }
    esd.EventLabel=newZ
  }
}


func (esd *ESD) Print() { 
  fmt.Println("Labeling")
  for eID,ev := range(esd.Label) {
    fmt.Println(eID, ev.Words)
    for pID, w := range(ev.Participants) {
      fmt.Println("    ", pID, w)
    }
  }
  fmt.Println("\nTau : ", esd.Tau)
  fmt.Println("V   : ", esd.V)
  fmt.Println("Pi  : ", esd.Pi)
  fmt.Println("eLab: ", esd.EventLabel)
}