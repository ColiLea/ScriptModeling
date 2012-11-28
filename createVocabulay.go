package scriptModeling

import "fmt"
import "os/exec"
import "strings"
import "strconv"

type Vocabulary struct {
  Dictionary Dictionary
  Covariances *CovarianceStruct
  Equivalences EquivalenceStruct
}

type CovarianceStruct struct {
  Matrix [][]float64
  Inverse [][]float64
}

type Dictionary struct {
  Vtoi map[string]int
  Itov map[int]string
  VList []string
  POSList []string
}

type EquivalenceStruct map[int][]int

var vocabulary Vocabulary
var vocabIdx int

// stem words and add them to vocabulary
func (vocabulary Dictionary) add(words []string, mode string) (wordIDs []int) {
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


func (vocabulary *Vocabulary)CreateCovarianceMatrix() {
  vocabulary.Covariances = new(CovarianceStruct)
  args := []string{"wnCovarianceSvsH.py", "--vocabulary"}
  args = append(args, vocabulary.Dictionary.VList...)
  args = append(args, "--pos")
  args = append(args, vocabulary.Dictionary.POSList...)
  cmd := new(exec.Cmd)
  cmd.Args = args
//   cmd.Path = "/home/lea/Code/Python/wnCovariance.py"
  cmd.Path = "/local/lea/thesis/python/wnCovarianceSvsH.py"
  out,err := cmd.Output()
  if err != nil {
    fmt.Println(err)
  }
  vocabulary.parse(string(out))
}


func (vocabulary *Vocabulary)CreateEquivalenceClasses() {
  var simIdx int
  threshold := 0.03
  vocabulary.Equivalences = make(map[int][]int, len(vocabulary.Dictionary.VList))
  for w1ID, row := range(vocabulary.Covariances.Matrix) {
    simIdx=0
    vocabulary.Equivalences[w1ID]=make([]int, len(vocabulary.Dictionary.VList))
    for w2ID, similarity := range(row) {
      if similarity > threshold {
	vocabulary.Equivalences[w1ID][simIdx]=w2ID
	simIdx++	
      }
    }
    vocabulary.Equivalences[w1ID] = vocabulary.Equivalences[w1ID][:simIdx]
  }
}


func (vocabulary *Vocabulary) parse(input string) {
  var v float64
  matrices := strings.Split(input, "\n")
  matrix := strings.Split(matrices[0], ";")
  inverse := strings.Split(matrices[1], ";")
  for idx, row := range(matrix) {
    matRow := make([]float64, len(strings.Split(row, " ")))
    invRow := make([]float64, len(strings.Split(row, " ")))
    for mIdx,val := range(strings.Split(row, " ")) {
      v,_ = strconv.ParseFloat(val, 64)
      matRow[mIdx] = v
    }
    for iIdx,val := range(strings.Split(inverse[idx], " ")) {
      v,_ = strconv.ParseFloat(val, 64)
      invRow[iIdx] = v
    }
    vocabulary.Covariances.Matrix = append(vocabulary.Covariances.Matrix, matRow)
    vocabulary.Covariances.Inverse = append(vocabulary.Covariances.Inverse, invRow)
  }
  vocabulary.Covariances.Matrix = vocabulary.Covariances.Matrix[:len(vocabulary.Covariances.Matrix)-1]
  vocabulary.Covariances.Inverse = vocabulary.Covariances.Inverse[:len(vocabulary.Covariances.Inverse)-1]
  if len(vocabulary.Covariances.Matrix) != len(vocabulary.Covariances.Matrix[0]) {
    fmt.Println(len(vocabulary.Covariances.Matrix), len(vocabulary.Covariances.Matrix[0]))
    panic("Matrix isn't square!")
  } 
  if len(vocabulary.Covariances.Inverse) != len(vocabulary.Covariances.Inverse[0]) {
    panic("Inverse isn't square!")
  }
}


func (sim *EquivalenceStruct) Print() {
  for key, val := range(*sim) {
    fmt.Println(vocabulary.Dictionary.Itov[key])
    for _, wID := range(val) {
      fmt.Println("	", vocabulary.Dictionary.Itov[wID])
    }
  }
}