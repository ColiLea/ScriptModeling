package scriptModeling

import "fmt"
import "math/rand"

const numTop int = 7

func resample_t(esd *ESD) {
  p := pick_event(esd.Tao)
  fmt.Println("resampling t=", esd.Tao, "...for event", p)
  alternatives := flip_event(*esd, p)
  fmt.Println("Alternatives: ", alternatives, "\n")
}

func pick_event(tao [numTop]int) int {
  //randomly select the event we want to resample
  var el int
  for alt:=0 ; alt!=1; {
    el = rand.Intn(len(tao))
    alt = tao[el]
  }
  return el
}

func flip_event(esd ESD, p int) map[[numTop]int][]int{
  //compute all possible tao vectors with the selected event p flipped to an alternative position and the corresponding ESD labeling
  //return: map[tao]label
  alts := make(map[[numTop]int][]int)
  tmp := new([numTop]int)
  alts[esd.Tao]=esd.Label
  for idx,_:=range(esd.Tao){
    *tmp=esd.Tao
    if esd.Tao[idx]==0{
      tmp[idx]=1
      tmp[p]=0
      alts[*tmp]=computeZ(*tmp, esd.Pi, numTop)
    }
  }
  return alts
}
