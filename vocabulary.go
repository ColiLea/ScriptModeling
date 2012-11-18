package scriptModeling

import "fmt"
import "os/exec"
import "strings"
import "leaMatrix"


type VocabMap struct {
  Vtoi map[string]int
  Itov map[int]string
  VList []string
  POSList []string
}

type similarities struct {
  matrix leaMatrix.Matrix
  equivalenceClasses map[int][]int
}

var vocabulary VocabMap
var vocabIdx int

// stem words and add them to vocabulary
func (vocabulary VocabMap) add(words []string, mode string) (wordIDs []int) {
  var lemma string
  wordIDs = make([]int, len(words))
  for idx, word := range(words) {
    if _, ok := vocabulary.Vtoi[word]; ok {
      wordIDs[idx]=vocabulary.Vtoi[word]
    } else {
      python := new(exec.Cmd)
      args := []string{"lemmatizer.py", "--vocabulary"}
      args = append(args, word)
      args = append(args,"--pos")
      args = append(args, mode)
      python.Args = args
//       python.Path = "/home/lea/Code/Python/lemmatizer.py"
      python.Path = "/local/lea/thesis/python/lemmatizer.py"
      out, _ := python.Output()
      lemma = strings.Trim(string(out), "\n")
      if _, ok := vocabulary.Vtoi[lemma]; !ok {
	fmt.Println(vocabIdx)
	vocabulary.VList[vocabIdx] = lemma
	vocabulary.Vtoi[lemma]=vocabIdx
	vocabulary.Itov[vocabIdx]=lemma
	vocabulary.POSList[vocabIdx]=mode
	vocabIdx++
      }
      wordIDs[idx]=vocabulary.Vtoi[lemma]
    }
  }
  return wordIDs
}

func GetVocabulary() []string {
  return vocabulary.VList
}

func getCovarianceMatrix(outFile string) (matrix *leaMatrix.Matrix) {
  fmt.Println(vocabulary.VList)
  fmt.Println(vocabulary.POSList)
  args := []string{"wnCovarianceSvsH.py", "--vocabulary"}
  args = append(args, vocabulary.VList...)
  args = append(args, "--pos")
  args = append(args, vocabulary.POSList...)
  cmd := new(exec.Cmd)
  cmd.Args = args
//   cmd.Path = "/home/lea/Code/Python/wnCovariance.py"
  cmd.Path = "/local/lea/thesis/python/wnCovarianceSvsH.py"
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
	fmt.Println(vocabulary.Itov[wordID], vocabulary.Itov[w2ID], value)
	simIdx++
      }
    }
    sim.equivalenceClasses[wordID] = sim.equivalenceClasses[wordID][:simIdx]
  }
}

func (sim *similarities) Print() {
  for key, val := range(sim.equivalenceClasses) {
    fmt.Println(vocabulary.Itov[key])
    for _, wID := range(val) {
      fmt.Println("	", vocabulary.Itov[wID])
    }
  }
}