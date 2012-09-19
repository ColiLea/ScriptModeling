 package scriptModeling

import "math/rand"
import "math"
import "fmt"

func (esd *ESD) hasParticipants() bool {
  // check whether there are any participants in the esd
  for _,event := range(esd.Label) {
    if len(event.Participants) > 0 {
      return true
    }
  }
  return false
}

func Pick_participant(label *Label) [2]int {
  // Pick a random participant type to resample from the esd labeling
  participants := make([][2]int, numTop*numPar)
  var idx int
  for eID, event := range(*label) {
    for pID,_ := range(event.Participants) {
      participants[idx]=[2]int{eID, pID}
      idx++
    }
  }
  target := rand.Intn(len(participants[:idx]))
  fmt.Println("Resampling p=", participants[:idx] , " for participanttype", participants[target])
  return participants[target]
}
  
func getAlternatives(participant int, label map[int][]string) []int {
  // Get alternative participant types ; TODO: ugly function!!
  var add bool
  idx := 1
  alts := make([]int, numPar-len(label)+1)
  alts[0] = participant
  for ii:=0 ; ii<numPar ; ii++ {
    add=false
    if ii!= alts[0] {
      add=true
      for pID,_ := range(label) {
	if pID==ii {
	  add=false
	  break
	}
      }
    }
    if add == true {
      alts[idx]=ii
      idx++
    }
  }
  return alts
}

func (sampler *Sampler) Resample_p(esd *ESD, targets [2]int) {
  var lgamma, distTotal, totalgamma, documentLikelihood, totaldoclikelihood, update, pPositive, pNegative, pNormalize float64
  var distribution []float64
  var docLikelihoods []float64
  var newV int
  event := targets[0]
  target := targets[1]
  // Get alternative participant types
  alternatives := getAlternatives(target, esd.Label[event].Participants) 
  proposedLabels := make([]Label, len(alternatives))  
  // Decrement Counts
  sampler.Model.participanttype_histogram[target]--
  sampler.Model.participanttype_eventtype_histogram[target][event]--
  if sampler.Model.participanttype_histogram[target]<0 {
    panic("Negative Participant Count")
  }
  sampler.Model.UpdateParticipantWordCounts(target, esd.Label[event].Participants[target], -1)
  if sampler.Model.participanttype_eventtype_histogram[target][event]<0 {
    panic("Negative Event Participant Count in resample_p")
  }
  // Compute likelihood for every type
  distribution = make([]float64, len(alternatives))
  docLikelihoods = make([]float64, len(alternatives))
  for idx, proposedP := range(alternatives) {
    if idx==0 {
      proposedLabels[idx]=esd.Label
    } else {
      esd.UpdateLabelingP(event, alternatives[idx-1], proposedP)
      proposedLabels[idx]=esd.Label
    }
    target=alternatives[idx]
    lgamma = 1.0
    for i:=0 ; i<numPar ; i++ {
      update = 0.0
      if i==proposedP {update = 1.0}
      pPositive = math.Gamma(float64(sampler.Model.participanttype_eventtype_histogram[i][event]) + sampler.participantPosPrior + update)
      pNegative = math.Gamma(float64(sampler.Model.eventtype_histogram[event]-sampler.Model.participanttype_eventtype_histogram[i][event]) + sampler.participantNegPrior - update)
      pNormalize = math.Gamma(float64(sampler.Model.participanttype_histogram[i])+sampler.participantPosPrior+sampler.participantNegPrior+update)
      lgamma *= ((pPositive*pNegative)/pNormalize)
    }
    documentLikelihood = sampler.documentLikelihoodP(event, target, proposedLabels[idx])
    distribution[idx]=lgamma
    docLikelihoods[idx]=documentLikelihood
    totaldoclikelihood += documentLikelihood
    totalgamma += lgamma
  }
  
  //compute document likelihood
  for idx,_ := range(distribution) {
     fmt.Println(distribution[idx]/totalgamma , docLikelihoods[idx]/totaldoclikelihood)
    distribution[idx] = (distribution[idx]/totalgamma) * (docLikelihoods[idx]/totaldoclikelihood)
    distTotal += distribution[idx]
  }
  for idx,_ := range(distribution) {
    distribution[idx]=distribution[idx]/distTotal
  }
  fmt.Println(distribution)
  newV = sample(distribution)
  fmt.Println(newV, "  = participanttype", alternatives[newV])
  //update esd and model
  esd.UpdateLabelingP(event, alternatives[len(alternatives)-1], alternatives[newV])
  sampler.Model.participanttype_histogram[alternatives[newV]]++
  sampler.Model.participanttype_eventtype_histogram[alternatives[newV]][event]++
  sampler.Model.UpdateParticipantWordCounts(alternatives[newV], esd.Label[event].Participants[alternatives[newV]], 1)
}