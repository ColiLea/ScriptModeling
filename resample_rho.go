 package scriptModeling
 
 
import "fmt"
import "strconv"
import "os/exec"
import "strings"

func (sampler *Sampler)Resample_rho() {
  var v_0, nu_0 float64
  var totalV, numDocs, nminusj int
  var slicesampler string
  for idx,target := range(sampler.rho) {
    lastRho := target
    target := idx
    totalV = sampler.Model.invcount_histogram[target]
    v_0 = sampler.v_0[target]
    nu_0 = sampler.nu_0
    numDocs = sampler.Model.numESDs
    nminusj = numTop-target

    slicesampler = getSliceSampler([]string{"1", "3", "@logposterior", strconv.FormatFloat(lastRho, 'f', -1 , 64), "5", "false", strconv.Itoa(totalV), strconv.FormatFloat(v_0, 'f', -1, 64), strconv.FormatFloat(nu_0, 'f', -1, 64), strconv.Itoa(numDocs), strconv.Itoa(nminusj)})    

    cmd := exec.Command("octave", "-q")
    cmd.Dir = "/home/lea/Code/Octave/"
    cmd.Stdin = strings.NewReader(slicesampler)
    out, err := cmd.Output()
    
    newRho,_ := strconv.ParseFloat(strings.Split(strings.TrimSpace(string(out)), " ")[3], 64)
    sampler.rho[idx]=newRho
    if err != nil {
      fmt.Println(err)
    }
  }
}


func getSliceSampler(args []string) string {
  return "samples = slicesample("+strings.Join(args, ", ")+`)`
}