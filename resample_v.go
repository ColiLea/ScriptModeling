 package scriptModeling

  import "fmt"
  import "math/rand"
  import "math"


func pick_invcount(v [numTop-1]int) int {
  newV := rand.Intn(len(v))
  fmt.Println("Resampling v=", v , " for eventtype", newV)
  return newV
}

func vPrior (rho0 float64) [numTop-1]float64 {
  var vPrior [numTop-1]float64
  for j:=0 ; j<numTop-1 ; j++ {
    vPrior[j] = (1.0/(math.Exp(rho0)-1.0))-((float64(numTop)-float64(j)+1.0)/(math.Exp((float64(numTop)-float64(j)+1.0)*rho0)-1.0))
  }
  return vPrior
}