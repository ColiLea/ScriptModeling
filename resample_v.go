package scriptModeling

  import "fmt"
  import "math/rand"

func resample_v(esd *ESD) {
  p := pick_invcount(esd.V)
  fmt.Println("resampling v= ", esd.V, "for event", p)
  alternatives := flip_invcount(*esd, p)
  fmt.Println("Alternatives: ", alternatives, "\n")
}

func pick_invcount(v [numTop-1]int) int {
  return rand.Intn(len(v))
}

func flip_invcount(esd ESD, c int) map[[numTop-1]int]struct{pi [numTop]int; label []int} {
  alternatives := make(map[[numTop-1]int]struct{pi [numTop]int; label []int})
  tmp:=new([numTop-1]int)
  alternatives[esd.V] = struct{pi [numTop]int; label []int}{esd.Pi,esd.Label}
  for ii:=0; ii<=(numTop-1)-c; ii++ {
    *tmp=esd.V
    if ii!= esd.V[c] {
      tmp[c]=ii
      newPi := computePi(*tmp, numTop)
      newZ := computeZ(esd.Tao, newPi, numTop)
      alternatives[*tmp]=struct{pi [numTop]int; label []int}{newPi,newZ}
    }
  }
  return alternatives
}