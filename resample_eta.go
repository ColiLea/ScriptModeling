package scriptModeling

import "strconv"
import "bytes"
import "strings"
import "fmt"

var octaveIdx int

func (sampler *Sampler)Resample_eta(eta []float64, i int, docLikelihood float64) (newEta float64) {
  
  octaveIdx++
  if octaveIdx == 500 {
    cmdIn.Write([]byte("exit"))
    StartOctave()
    octaveIdx = 0
  }
  
  slicesampler := getSliceSampler([]string{"1", "3", "@normalposterior", strconv.FormatFloat(eta[i], 'f', -1 , 64), "2", "false", String(eta), strconv.Itoa(i+1), sampler.covariances.matrix.InverseStr, strconv.FormatFloat(docLikelihood, 'f', -1, 64)})

  //try up to 10 times
    
  for i:= 0 ; i<10 ; i++ {
    cmdIn.Write(slicesampler)
    out, err := cmdOut.ReadString('\n')
    if err == nil {
      newEta,_ = strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
      break
    } else if i==9 {
        fmt.Println(err)
	fmt.Println(string(out))
	newEta = eta[i]
	err = nil
	break
    } else {
      fmt.Println(err)
      fmt.Println(string(out))
      err = nil
    }
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

func (sampler *Sampler) updateEta(diff map[int][]int, newDL float64, diff2 map[int][]int, oldDL float64,  mode string) {
  // new class
  fmt.Println(diff, diff2)
  fmt.Println("new")
  for class,words := range(diff) {
    sampler.updateEtaClass(class, words, newDL, mode)
    sampler.updatePrior(class, mode)
  }
  // old class   
//   fmt.Println("diff2 (should DECREASE): ")
  fmt.Println("old")
  for class,words := range(diff2) {
    sampler.updateEtaClass(class, words, oldDL, mode)
    sampler.updatePrior(class, mode)
  }
}

func (sampler *Sampler) updateEtaClass(class int, words []int, wordLikelihood float64, mode string) {
//   var wordLikelihood float64
  if mode == "event" {
//     wordLikelihood = sampler.wordLikelihood(class, mode)
    fmt.Println("WordLikelihood", wordLikelihood)
    for _,target := range(words) {
//       fmt.Println(">>>",vocabulary.Itov[target],"<<<")
      for _,word := range(sampler.covariances.equivalenceClasses[target]) {
// 	fmt.Println("  >>",vocabulary.Itov[word],"<<", sampler.EventlmPriors[class][word], sampler.eventEtas[class][word], wordLikelihood)
	sampler.eventEtas[class][word] = sampler.Resample_eta(sampler.eventEtas[class], word, wordLikelihood)
// 	fmt.Println(sampler.eventEtas[class][word], "\n---------\n")
      }
    }
  } else {
//     wordLikelihood = sampler.wordLikelihood(class, mode)
    fmt.Println("WordLikelihood", wordLikelihood)
    for _,target := range(words) {
//       fmt.Println(">>>",vocabulary.Itov[target],"<<<")
      for _,word := range(sampler.covariances.equivalenceClasses[target]) {
// 	fmt.Println("  >>",vocabulary.Itov[word],"<<", sampler.ParticipantlmPriors[class][word], sampler.participantEtas[class][word], wordLikelihood)
	sampler.participantEtas[class][word] = sampler.Resample_eta(sampler.participantEtas[class], word, wordLikelihood)
// 	fmt.Println(sampler.participantEtas[class][word], "\n---------\n")
      }
    }
  }
}