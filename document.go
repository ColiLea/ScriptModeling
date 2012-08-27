package scriptModeling

import "fmt"

type ESD struct {
// Elemantary-Sequence-Description, consisting of Lengt=number of events; Events=events(words); Tao=event realizations; V=inversions; Pi=global Ordering; Label=event labels
  Length int
  //TODO make Events [][]string (-> multiword events)
  Events [][]string
  Tao [numTop]int
  V [numTop-1]int
  Pi [numTop]int
  Label []int
}

type Corpus []*ESD

func (esd *ESD) ComputePi(K int) {
// Compute global labeling from V (inversion count vector)
  esd.Pi[0] = K-1
  for j:=K-2; j>=0; j-- {
   for i:=K-2; i>=esd.V[j]; i-- {
     esd.Pi[i+1]=esd.Pi[i]
   }
   esd.Pi[esd.V[j]]=j
  }
}

func (esd *ESD) ComputeZ(K int) {
// Compute the ESD labeling from Tao (realization vector) and Pi (global labeling)
  esd.Label=make([]int,len(esd.Events))
  event:=0
  for _,el := range(esd.Pi) {
    if esd.Tao[el] == 1 {
      esd.Label[event]=el
      event++
    }
  }
}

func (esd *ESD) flipEvent(oldEvent int, newEvent int ) {
  esd.Tao[oldEvent]=0
  esd.Tao[newEvent]=1
}

func PrintESD(esd ESD) { 
  for ev,_ := range(esd.Events) {
    fmt.Println("String: ", esd.Events[ev], "  Label: ", esd.Label[ev])
  }
}