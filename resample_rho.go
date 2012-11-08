 package scriptModeling
 
 
import "fmt"
import "strconv"
import "strings"

func (sampler *Sampler)Resample_rho() {
  fmt.Println("Resampling Rho")
  var v_0, nu_0 float64
  var totalV, numDocs, nminusj int
  var slicesampler []byte
  for idx,target := range(sampler.Model.rho) {
    lastRho := target
    target := idx
    totalV = sampler.Model.invcount_histogram[target]
    v_0 = sampler.v_0[target]
    nu_0 = sampler.nu_0
    numDocs = sampler.Model.numESDs
    nminusj = numTop-target

    slicesampler = getSliceSampler([]string{"1", "3", "@logposterior", strconv.FormatFloat(lastRho, 'f', -1 , 64), "5", "false", strconv.Itoa(totalV), strconv.FormatFloat(v_0, 'f', -1, 64), strconv.FormatFloat(nu_0, 'f', -1, 64), strconv.Itoa(numDocs), strconv.Itoa(nminusj)})
    fmt.Println(string(slicesampler))
    cmdIn.Write(slicesampler)
    out, err := cmdOut.ReadString('\n')
    newRho,_ := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
    sampler.Model.rho[idx]=newRho
    if err != nil {
      fmt.Println(err)
    }
  }
}