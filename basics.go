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
  labels := make([][]int, len(newEs))
  for idx,newE := range(newEs) {
    tmpTao := esd.Tao
    tmpTao[oldE]=0
    tmpTao[newE]=1
    labels[idx]=computeZ(tmpTao, esd.Pi)
  }
  return labels
}
