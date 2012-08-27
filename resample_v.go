package scriptModeling

  import "fmt"
  import "math/rand"
  import "math"


func pick_invcount(v [numTop-1]int) int {
  newV := rand.Intn(len(v))
  fmt.Println("Resampling v=", v , " for eventtype", newV)
  return newV
}

func vPrior (rho0 float64) []float64 {
  vPrior := make([]float64, numTop)
  for j:=0 ; j<numTop ; j++ {
    vPrior[j] = (1.0/(math.Exp(rho0)-1.0))-((float64(numTop)-float64(j)+1.0)/(math.Exp(float64(numTop)-float64(j)+1.0)-1.0))
  }
  return vPrior
}