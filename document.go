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
  Tau [numPar]int
}

type Corpus []*ESD

func (corpus *Corpus) Init() {
  for _,esd := range(*corpus) {
    esd.Init()
  }
}

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

func (esd *ESD) flip(oldEvent int, newEvent int) {
//   fmt.Println("$", oldEvent, newEvent)
//   fmt.Println("$", esd.Tau)
  esd.Tau[oldEvent]=0
  esd.Tau[newEvent]=1
//   fmt.Println("$", esd.Tau)
}

func (esd *ESD) flipp(oldPart int, newPart int, eID int ) {
  tmpTau := esd.Label[eID].Tau
  tmpTau[oldPart]=0
  tmpTau[newPart]=1
  content := Content{esd.Label[eID].Words, esd.Label[eID].Participants, tmpTau}
  esd.Label[eID] = content
}


func (esd *ESD) UpdateLabelingT() {
  oldZ := esd.EventLabel
  contents := make([]Content, len(oldZ))
  esd.ComputeZ()
  for idx, val := range(oldZ) {
    contents[idx]=esd.Label[val]
  }
  esd.Label = Label{}
  for idx, val := range(esd.EventLabel) {
    esd.Label[val]=contents[idx]
  }
}

func (esd *ESD) UpdateLabelingP(eventIdx int, oldVal int, newVal int) {
  if oldVal != newVal {
    esd.Label[eventIdx].Participants[newVal]=esd.Label[eventIdx].Participants[oldVal]
    delete(esd.Label[eventIdx].Participants, oldVal)
  }
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

func (esd *ESD) Copy() (newESD ESD) {
  newESD.Label = Label{}
  newESD.Length = esd.Length
  newESD.Tau = esd.Tau
  newESD.Pi = esd.Pi
  newESD.V = esd.V
  newESD.EventLabel = esd.EventLabel
  for key,_ := range(esd.Label) {
    tmpP := make(map[int][]string, len(esd.Label[key].Participants))
    for pkey, pVal := range(esd.Label[key].Participants) {
      tmpP[pkey] = pVal
    }
    newESD.Label[key]=Content{esd.Label[key].Words, tmpP, esd.Label[key].Tau}
  }
  return
}


func (esd *ESD) Print() {
  fmt.Println("Labeling")
  for eID,ev := range(esd.Label) {
    fmt.Println(eID, ev.Words)
    for pID, w := range(ev.Participants) {
      fmt.Println("    ", pID, w)
    }
    fmt.Println("    ", ev.Tau)
  }
  fmt.Println("\nTau : ", esd.Tau)
  fmt.Println("V   : ", esd.V)
  fmt.Println("Pi  : ", esd.Pi)
  fmt.Println("eLab: ", esd.EventLabel)
}