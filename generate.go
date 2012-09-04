package scriptModeling

import "math"
// import "fmt"

func (model *Model) Generate(jPrior, lmPrior/* jlmPrior, ilmPrior*/ float64) *ESD {
  var wList []string
  rho := [numTop-1]float64{0.4, 0.2}
  esd := new(ESD)
  esd.Label = Label{}
  //Generate Eventtypes
  for jj := 0 ; jj<numTop ; jj++ {
    jPos := (float64(model.eventtype_histogram[jj])+jPrior)/(float64(model.numESDs)+jPrior)
    jNeg := (float64(model.numESDs-model.eventtype_histogram[jj])+jPrior)/(float64(model.numESDs)+jPrior)
    esd.Tau[jj]=sample([]float64{jNeg,jPos})
    if esd.Tau[jj] ==1 {
      esd.Length++
      esd.Label[jj]=Content{[]string{}, map[int][]string{}}
      //Generate Participants
      for ii:=0 ; ii<numPar ; ii++ {
	pPos := (float64(model.participanttype_eventtype_histogram[ii][jj])+jPrior)/(float64(model.participanttype_histogram[ii])+jPrior)
	pNeg := (float64(model.participanttype_histogram[ii]-model.participanttype_eventtype_histogram[ii][jj])+jPrior)/(float64(model.participanttype_histogram[ii])+jPrior)
	pp := sample([]float64{pNeg,pPos})
	if pp == 1 {esd.Label[jj].Participants[ii] = []string{}}
      }
    }
    //Generate ordering
    vDist := make([]float64, numTop-jj)
    for vv:=0 ; vv<numTop-jj && jj<numTop-1 ; vv++ {
      vDist[vv] = math.Abs(math.Exp(-rho[jj]*float64(vv))/((1.0-math.Exp(-(float64(numTop-jj+1))*rho[jj]))/1.0-math.Exp(-rho[jj])))
    }
    if jj<numTop-1 {
      esd.V[jj]= getAccumulativeSample(vDist)
    }
    esd.ComputePi()
    esd.ComputeZ()
  }
  //Generate words
  for eID, event := range(esd.Label) {
    wList = []string{}
    wDist := make([]float64, model.eventVocabulary)
    words := make([]string, model.eventVocabulary)
    idx:=0
    for term,dist := range(model.word_eventtype_histogram) {
      words[idx]=term
      wDist[idx]=(float64(dist[eID])+lmPrior)
      idx++
    }
    ww := getAccumulativeSample(wDist)
    wList = append(wList, words[ww])
    esd.Label[eID] = Content{wList, esd.Label[eID].Participants}
    for pID, _ := range(event.Participants) {
      idx=0
      wList = []string{}
      wDist = make([]float64, model.participantVocabulary)
      words = make([]string, model.participantVocabulary)
      for term,dist := range(model.word_participanttype_histogram) {
	words[idx]=term
	wDist[idx]=(float64(dist[pID])+lmPrior)
	idx++
      }
      ww = getAccumulativeSample(wDist)
      wList = append(wList, words[ww])
      esd.Label[eID].Participants[pID]=wList
    }
  }
  return esd
}

func GetModel() *Model {
  model := new(Model)
  model.numESDs = 5
  model.eventVocabulary = 3
  model.participantVocabulary = 2
  model.eventtype_histogram = Histogram{1,4,2}
  model.participanttype_histogram = Histogram{4,2}
  model.participanttype_eventtype_histogram = map[int]Histogram{0:Histogram{1,2,1}, 1:Histogram{0,1,1}}
  model.word_eventtype_histogram = map[string]Histogram{"go":Histogram{0,2,1}, "read":Histogram{1,0,1}, "sing":Histogram{0,2,0}}
  model.word_participanttype_histogram = map[string]Histogram{"lea":Histogram{3,1}, "dom":Histogram{1,1}}
  model.invcount_histogram= Histogram{8,3}
  return model
}