 package scriptModeling

import "math/rand"
import "fmt"

type Histogram []int
type Distribution []float64

func newHistogram(topics int) Histogram {
  return make([]int, topics)
}


func getAccumulativeSample(distribution Distribution) int {
  // normalize
  sum := 0.0
  for _,val := range(distribution) {sum += val}
  for dIdx,_ := range(distribution) {distribution[dIdx] = distribution[dIdx]/sum}
  fmt.Println(distribution)
  // get sample
  distribution_sum := 0.0
  for _, v := range distribution {
    distribution_sum += v
  }
  choice := rand.Float64() * float64(distribution_sum)
  sum_so_far := 0.0
  for i, v := range distribution {
    sum_so_far += v
    if sum_so_far >= choice {
      return i;
    }
  }
  return -1;
}

func computePi(v [numTop-1]int) [numTop]int {
// Compute global labeling from V (inversion count vector)
  var pi [numTop]int
  pi[0] = numTop-1
  for j:=numTop-2; j>=0; j-- {
   for i:=numTop-2; i>=v[j]; i-- {
     pi[i+1]=pi[i]
   }
   pi[v[j]]=j
  }
  return pi
}

func computeZ(tao [numTop]int, pi [numTop]int) []int{
// Compute the ESD labeling from Tao (realization vector) and Pi (global labeling)
  label := make([]int,numTop)
  event:=0
  for _,el := range(pi) {
    if tao[el] == 1 {
      label[event]=el
      event++
    }
  }
  return label[:event]
}

func getLabels(esd ESD, oldE int, newEs []int)  [][]int {
  // Get event labels for all suggested flips
  labels := make([][]int, len(newEs))
  for idx,newE := range(newEs) {
    tmpTao := esd.Tau
    tmpTao[oldE]=0
    tmpTao[newE]=1
    labels[idx]=computeZ(tmpTao, esd.Pi)
  }
  return labels
}


func getPLabels(currentLabel [][]int, target int, event int, proposals []int)  [][][]int {
  // get participant labels for all suggested flips
  labels := make([][][]int, len(proposals))
  var tmpLabel [][]int
  for propIndex,pp := range(proposals) {
    tmpLabel = make([][]int, len(currentLabel))
    for idx, el := range(currentLabel) {
      tmpLabel[idx]=make([]int, len(currentLabel[idx]))
      for ii, p := range(el) {
	if ii==target && idx==event {
	  tmpLabel[idx][ii]=pp
	} else {
	  tmpLabel[idx][ii]=p
	}
      }
    }
    labels[propIndex]=tmpLabel
  }
  return labels
}


func updateLabeling(event int, oldVal int, newVal int, label Label, mode string) Label {
  newLabel := make(Label, len(label))
  if oldVal == newVal {
    return label
  }
  if mode == "event" {
    for k, v := range(label) {
      if k==oldVal {
	newLabel[newVal]=v
      }
      newLabel[k]=v
      delete(newLabel, oldVal)
    }
  } else if mode == "participant" {
    for k,v := range(label) {
      newLabel[k]=v
    }
    for k,v := range(label[event].Participants) {
      if k==oldVal {
	newLabel[event].Participants[newVal]=v
      }
      newLabel[event].Participants[k]=v
      delete(newLabel[event].Participants, oldVal)
    }
  }
  return newLabel
}

func UpdateLabelingV(tau [numTop]int, pi [numTop]int, eventLabel []int, label Label) Label {
  newZ := computeZ(tau, pi)
  newLabel := make(Label, len(label))
  update:=false
  for idx,_ := range(newZ) {
    if newZ[idx] != eventLabel[idx] {
      update=true
    }
  }
  if update==false {
    return label
  } else {
    contents := make([]Content, len(label))
    for idx,eID := range(eventLabel) {
      contents[idx]=label[eID]
    }
    for idx, id := range(newZ) {
      newLabel[id]=contents[idx]
    }
  }
  return newLabel
}