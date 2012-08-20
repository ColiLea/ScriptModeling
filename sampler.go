package scriptModeling

import "math/rand"

func PickVariable(esd *ESD) {
  //select which random variable to resample; 0:t  1:v  2:rho
  rr := rand.Intn(3)
  if rr==0{
    resample_t(esd)
  } else if rr==1{
    resample_v(esd)
  } else {
    resample_rho()
  }
}
