package scriptModeling

import "strconv"
import "bytes"
import "strings"
import "fmt"
import "math"

func (sampler *Sampler)Resample_eta(eta []float64, i int, docLikelihood float64) (newEta float64) {
  
  slicesampler := getSliceSampler([]string{"1", "3", "@normalposterior", strconv.FormatFloat(eta[i], 'f', -1 , 64), "2", "false", String(eta), strconv.Itoa(i+1), sampler.covariances.matrix.InverseStr, strconv.FormatFloat(docLikelihood, 'f', -1, 64)})

  cmdIn.Write(slicesampler)
  out, err := cmdOut.ReadString('\n')
  
  newEta,_ = strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
  if err != nil {
    fmt.Println(err)
    fmt.Println(string(slicesampler))
    fmt.Println(string(out))
    newEta=eta[i]
  } else if math.IsNaN(newEta) {
    fmt.Println(string(slicesampler))
    fmt.Println(string(out))
  }
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

func (sampler *Sampler) updateEta(diff, diff2 map[int][]int, mode string) {
  // new class
  fmt.Println("Ordering", diff, diff2)
  for class,words := range(diff) {
    sampler.updateEtaClass(class, words, mode)
  }
  // old class   
  fmt.Println("diff2 (should DECREASE): ")
  for class,words := range(diff2) {
    sampler.updateEtaClass(class, words, mode)
  }
}

func (sampler *Sampler) updateEtaClass(class int, words []int, mode string) {
  var wordLikelihood float64
  if mode == "event" {
    wordLikelihood = sampler.wordLikelihood(class, mode)
    for _,target := range(words) {
      fmt.Println(">>>",vocabulary.itov[target],"<<<")
      for _,word := range(sampler.covariances.equivalenceClasses[target]) {
	fmt.Println("  >>",vocabulary.itov[word],"<<")
	fmt.Println(sampler.EventlmPriors[class][word], wordLikelihood)
	sampler.eventEtas[class][word] = sampler.Resample_eta(sampler.eventEtas[class], word, wordLikelihood)
	sampler.updatePrior(class, mode)
	fmt.Println(sampler.EventlmPriors[class][word], "\n---------\n")
      }
    }
  } else {
    wordLikelihood = sampler.wordLikelihood(class, mode)
    for _,target := range(words) {
      fmt.Println(">>>",vocabulary.itov[target],"<<<")
      for _,word := range(sampler.covariances.equivalenceClasses[target]) {
	fmt.Println("  >>",vocabulary.itov[word],"<<")
	fmt.Println(sampler.ParticipantlmPriors[class][word], wordLikelihood)
	sampler.participantEtas[class][word] = sampler.Resample_eta(sampler.participantEtas[class], word, wordLikelihood)
	sampler.updatePrior(class, mode)
	fmt.Println(sampler.ParticipantlmPriors[class][word], "\n---------\n")
      }
    }
  }
}