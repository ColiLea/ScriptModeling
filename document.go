package scriptModeling

import "fmt"


// Elemantary-Sequence-Description
type ESD struct {
  Label Label
  Length int
  Tau [numTop]int
  V [numTop-1]int
  Pi [numTop]int
  EventLabel []int
}

type Label map[int]Content

type Content struct {
  Words []int
  Participants map[int][]int
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


// Compute global labeling from V (inversion count vector)
func (esd *ESD) ComputePi() {
  esd.Pi[0] = numTop-1
  for j:=numTop-2; j>=0; j-- {
    for i:=numTop-2; i>=esd.V[j]; i-- {
      esd.Pi[i+1]=esd.Pi[i]
    }
    esd.Pi[esd.V[j]]=j
  }
}


// Compute the ESD labeling from Tao (realization vector) and Pi (global labeling)
func (esd *ESD) ComputeZ() {
  esd.EventLabel=make([]int,len(esd.Label))
  event:=0
  for _,el := range(esd.Pi) {
    if esd.Tau[el] == 1 {
      esd.EventLabel[event]=el
      event++
    }
  }
}


// Compute the inversion counts V from Pi (global labeling)
func (esd *ESD) ComputeV() {
  for idx,el := range(esd.Pi) {
    if el < numTop-1 {
      invCount := 0
      for vIdx:=0 ; vIdx<idx ; vIdx++ {
	if esd.Pi[vIdx]>el {
	  invCount++
	}
      }
      esd.V[el]=invCount
    }
  }
}

func (esd *ESD) flip(oldEvent int, newEvent int) {
  esd.Tau[oldEvent]=0
  esd.Tau[newEvent]=1
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
    tmpP := make(map[int][]int, len(esd.Label[key].Participants))
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