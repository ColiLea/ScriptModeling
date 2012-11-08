package scriptModeling

import "strconv"
import "bytes"
import "strings"
import "fmt"

func (sampler *Sampler)Resample_eta(eta []float64, i int, docLikelihood float64) (newEta float64) {
  fmt.Println(i, vocabulary.itov[i], eta[i], docLikelihood)
  
  slicesampler := getSliceSampler([]string{"1", "3", "@normalposterior", strconv.FormatFloat(eta[i], 'f', -1 , 64), "2", "false", String(eta), strconv.Itoa(i+1), sampler.covariances.InverseStr, strconv.FormatFloat(docLikelihood, 'f', -1, 64)})
  
//   fmt.Println(string(slicesampler))
  cmdIn.Write(slicesampler)

  out, err := cmdOut.ReadString('\n')
  newEta,_ = strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
  if err != nil {
    newEta=eta[i]
  }
  fmt.Println(string(out), newEta, "\n\n")
  return newEta
}

func String(eta []float64) string {
  var etaS bytes.Buffer
  etaS.WriteString("[")
  for idx,_ := range(eta){
      etaS.WriteString(strconv.FormatFloat(eta[idx], 'f', -1, 64)+";")
  }
  etaS.WriteString("]")
  return etaS.String()
}