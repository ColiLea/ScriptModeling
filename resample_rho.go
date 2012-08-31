 package scriptModeling

import "fmt"
import "math"

func Resample_rho() {
  fmt.Println("resampling rho...")
}

func (sampler *Sampler)rhoPosterior(rho float64, target int) (prob float64) {
  prob = -rho * (float64(sampler.Model.invcount_histogram[target]) + sampler.v_0[target]*sampler.nu_0) - (sampler.nu_0 + float64(sampler.Model.numESDs)) * (math.Log1p(-math.Exp((float64(numTop-target)) * rho)) - math.Log1p(-math.Exp(-rho)))
  return
}