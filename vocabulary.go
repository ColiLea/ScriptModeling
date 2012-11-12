package scriptModeling

import "fmt"
import "os/exec"
import "strings"
import "leaMatrix"

type vocabMap struct {
  vtoi map[string]int
  itov map[int]string
  VList []string
}

var vocabulary vocabMap
var vocabIdx int

func (voc vocabMap) add(words []string) {
  for _,word := range(words) {
    if _,ok := voc.vtoi[word]; !ok {
      voc.VList[vocabIdx] = word
      voc.vtoi[word]=vocabIdx
      voc.itov[vocabIdx]=word
      vocabIdx++
    }
  }
}

func GetVocabulary() []string {
  return vocabulary.VList
}

func GetCovarianceMatrix(vocab []string, outFile string) (matrix *leaMatrix.Matrix) {
  args := []string{"wnCovariance.py", "--vocabulary"}
  fmt.Println(vocab)
  args = append(args, vocab...)
  cmd := new(exec.Cmd)
  cmd.Args = args
//   cmd.Path = "/home/lea/Code/Python/wnCovariance.py"
  cmd.Path = "/local/lea/thesis/python/wnCovariance.py"
  out,err := cmd.Output()
  fmt.Println(string(out), err)
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