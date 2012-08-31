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

func updateLabelingP(event int, oldVal int, newVal int, label Label) Label {
  newLabel := make(Label, len(label))
  if oldVal == newVal {
    return label
  }
  for k,v := range(label) {
    newLabel[k]=v
  }
  if _,ok := label[event].Participants[oldVal] ; ok {
    newLabel[event].Participants[newVal]=newLabel[event].Participants[oldVal]
  }
  delete(newLabel[event].Participants, oldVal)
  return newLabel
}

func updateLabelingT(oldVal int, newVal int, label Label) Label {
  newLabel := make(Label, len(label))
  if oldVal == newVal {
    return label
  }
  for k,v := range(label) {
    newLabel[k]=v
    if _,ok := label[oldVal] ; ok {
      newLabel[newVal]=newLabel[oldVal]
    }
  }
  delete(newLabel, oldVal)
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