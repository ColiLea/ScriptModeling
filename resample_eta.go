package scriptModeling

import "leaMatrix"
import "math"

func (sampler *Sampler)Resample_eta(eta []float64, i int) (prob float64) {
  // NOTE I'll do all this stuff in Octave actually -.- ...
  // Pass the vectors (=COLUMN vectors) as v= [e1 ; e2 ; ...]
  // Pass the matrix as m = [11 , 12 , 13 ; 21 , 22 , 23 ; 31 , 32 , 33 ; ...]
  // save it like that..??

}


func getEtas(eta []float64, cmp int) (eta_i, eta_not []float64) {
  eta_i = make([]float64, len(eta))
  eta_not = make([]float64, len(eta))
  copy(eta_i, eta)
  eta_i[cmp] = 0.0
  for idx,_ := range(eta_not){
    if idx == cmp {
      eta_not[idx]=eta[idx]
    } else {
      eta_not[idx]=0
    }
  }
  return
}