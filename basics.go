 package scriptModeling

import "math"
import "math/rand"
// import "fmt"

type Histogram []int
type Distribution []float64

func newHistogram(topics int) Histogram {
  return make([]int, topics)
}

// samples from a normalized distribution; returns the INDEX of the value sampled
func sample(distribution []float64) int {
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


// samples from an unnormalized distribution; (1) normalizes (2)returns the INDEX of the value sampled
func getAccumulativeSample(distribution Distribution) int {
  // normalize
  sum := 0.0
  for _,val := range(distribution) {sum += val}
  for dIdx,_ := range(distribution) {distribution[dIdx] = distribution[dIdx]/sum}
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


func max(dist []float64) int {
  var pos int
  max := 0.0
  for idx, val := range(dist) {
    if val > max {
      pos = idx
    }
  }
  return pos
}

func computeNorm(dist []float64) (max, norm float64) {
  max=math.Inf(-1)
  for _,v := range(dist) {
    if v > max {
      max=v
    }
  }
  for _,v := range(dist) {
    norm += math.Exp(v-max)
  }
  return
}

// Compute global labeling from V (inversion count vector)
func computePi(v [numTop-1]int) [numTop]int {
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


// Compute the ESD labeling from Tao (realization vector) and Pi (global labeling)
func computeZ(tao [numTop]int, pi [numTop]int) []int{
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
  if Compare(newZ, eventLabel) == true {
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

func Compare(list1, list2 []int) bool {
  if len(list1) != len(list2) {
    return false
  } else {
    for idx, _ := range(list1) {
      if list1[idx] != list2[idx] {
	return false
      }
    }
  }
  return true
}

func isIn(el int, list []int) bool {
  for _,val := range(list) {
    if el==val {
      return true
    }
  }
  return false
}

func sum(numbers []float64) float64 {
  sum1 := 0.0
  for _, number := range(numbers) {
    sum1 += number
  }
  return sum1
}

func expSum(priors []float64) float64 {
  sum := 0.0
  for _, value := range(priors) {
    sum += math.Exp(value)
  }
  return sum
}

func normalized(dist []float64) []float64 {
  norm := sum(dist)
  for idx,_ := range(dist) {
    dist[idx]=dist[idx]/norm
  }
  return dist
}

func expNormalized(dist []float64) []float64 {
  max, norm := computeNorm(dist)
  for idx,_ := range(dist) {
    dist[idx] = math.Exp(dist[idx]-max)/norm
  }
  return dist
}