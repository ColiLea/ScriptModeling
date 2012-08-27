package scriptModeling

import "fmt"
import "math/rand"

const numTop int = 7

func pick_event(tao [numTop]int) int {
  //randomly select the event we want to resample
  var el int
  for alt:=0 ; alt!=1; {
    el = rand.Intn(len(tao))
    alt = tao[el]
  }
  fmt.Println("Resampling t=", tao , " for eventtype", el)
  return el
}

func newTargets(esd ESD, target int) []int {
  newTargets := make([]int, numTop)
  newTargets[0] = target
  idx := 1
  for eventtype,realized := range(esd.Tao) {
    if realized == 0 {
      newTargets[idx]=eventtype
      idx++
    }
  }
  return newTargets[:idx]
}
