package scriptModeling
 
// import "fmt"
import "sliceSampler"

func (sampler *Sampler)Resample_rho() {
//   fmt.Println("Resampling Rho")
  var v_0, nu_0, totalV, numDocs, nminusj float64
  for idx,target := range(sampler.Model.Rho) {
    lastRho := target
    target := idx
    totalV = float64(sampler.Model.Invcount_histogram[target])
    v_0 = sampler.v_0[target]
    nu_0 = sampler.nu_0
    numDocs = float64(sampler.Model.NumESDs)
    nminusj = float64(numTop-target)

    newRho := sliceSampler.SampleRho(3, 5.0, lastRho, false, totalV, v_0, nu_0, numDocs, nminusj)
    sampler.Model.Rho[idx]=newRho
  }
}