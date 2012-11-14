package scriptModeling

import "fmt"
import "os/exec"
import "strings"
import "leaMatrix"

type vocabMap struct {
  vtoi map[string]int
  itov map[int]string
  VList []string
  POSList []string
}

type similarities struct {
  matrix leaMatrix.Matrix
  equivalenceClasses map[int][]int
}

var vocabulary vocabMap
var vocabIdx int

func (voc vocabMap) add(words []string, mode string) {
  for _,word := range(words) {
    if _,ok := voc.vtoi[word]; !ok {
      voc.VList[vocabIdx] = word
      voc.vtoi[word]=vocabIdx
      voc.itov[vocabIdx]=word
      voc.POSList[vocabIdx]=mode
      vocabIdx++
    }
  }
}

func GetVocabulary() []string {
  return vocabulary.VList
}

func getCovarianceMatrix(outFile string) (matrix *leaMatrix.Matrix) {
  args := []string{"wnCovarianceT1.py", "--vocabulary"}
  args = append(args, vocabulary.VList...)
  args = append(args, "--pos")
  args = append(args, vocabulary.POSList...)
  cmd := new(exec.Cmd)
  cmd.Args = args
//   cmd.Path = "/home/lea/Code/Python/wnCovariance.py"
  cmd.Path = "/local/lea/thesis/python/wnCovarianceNoPOS.py"
  out,err := cmd.Output()
  fmt.Println(string(out))
  if err != nil {
    fmt.Println(err)
  }
  matrix = parse(string(out))
  matrix.Store(outFile)
  return matrix
}

func parse(input string) (matrix *leaMatrix.Matrix) {
  matrices := strings.Split(input, "\n")
  matrix = leaMatrix.NewMatrix(matrices[0])
  matrix.SetInverse(matrices[1])
  return
}

func (sim *similarities) getEquivalenceClasses() {
  var simIdx int
  threshold := 0.03
  sim.equivalenceClasses = make(map[int][]int, len(vocabulary.VList))
  for wordID,_ := range(sim.matrix.Data) {
    simIdx=0
    sim.equivalenceClasses[wordID]=make([]int, len(vocabulary.VList))
    for w2ID, value := range(sim.matrix.Data[wordID]) {
      if value > threshold {
	sim.equivalenceClasses[wordID][simIdx]=w2ID
	fmt.Println(vocabulary.itov[wordID], vocabulary.itov[w2ID], value)
	simIdx++
      }
    }
    sim.equivalenceClasses[wordID] = sim.equivalenceClasses[wordID][:simIdx]
  }
}

func (sim *similarities) Print() {
  for key, val := range(sim.equivalenceClasses) {
    fmt.Println(vocabulary.itov[key])
    for _, wID := range(val) {
      fmt.Println("	", vocabulary.itov[wID])
    }
  }
}