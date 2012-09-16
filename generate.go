package scriptModeling

import "math"
import "math/rand"
// import "fmt"

func (model *Model) Generate(jPrior, lmPrior float64) *ESD {
  var wList []string
  modelTop := 3
  modelPar := 3
  rho := [2]float64{1.9, 1.9}
  esd := new(ESD)
  esd.Label = Label{}
  //Generate Eventtypes
  for jj := 0 ; jj<modelTop ; jj++ {
    jPos := (float64(model.eventtype_histogram[jj])+jPrior)/(float64(model.numESDs)+jPrior)
    jNeg := (float64(model.numESDs-model.eventtype_histogram[jj])+jPrior)/(float64(model.numESDs)+jPrior)
    esd.Tau[jj]=sample([]float64{jNeg,jPos})
    if esd.Tau[jj] ==1 {
      esd.Length++
      esd.Label[jj]=Content{[]string{}, map[int][]string{}}
      //Generate Participants
      for ii:=0 ; ii<modelPar ; ii++ {
	pPos := (float64(model.participanttype_eventtype_histogram[ii][jj])+jPrior)/(float64(model.participanttype_histogram[ii])+jPrior)
	pNeg := (float64(model.participanttype_histogram[ii]-model.participanttype_eventtype_histogram[ii][jj])+jPrior)/(float64(model.participanttype_histogram[ii])+jPrior)
	pp := sample([]float64{pNeg,pPos})
	if pp == 1 {esd.Label[jj].Participants[ii] = []string{}}
      }
    }
    //Generate ordering
    vDist := make([]float64, modelTop-jj)
    for vv:=0 ; vv<modelTop-jj && jj<modelTop-1 ; vv++ {
      vDist[vv] = math.Exp(-rho[jj]*float64(vv+1))
    }
    if jj<modelTop-1 {
      esd.V[jj]= getAccumulativeSample(vDist)
//       fmt.Println(esd.V[jj])
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

func Randomize(esd ESD) (newESD ESD) {
  idx:=0
  newESD.Label = make(map[int]Content, len(esd.Label))
  eIDs := rand.Perm(numTop)[:len(esd.Label)]
  for _, val := range(esd.Label) {
    pIDs := rand.Perm(numPar)[:len(val.Participants)]
    pIdx:=0
    content := Content{}
    content.Words = val.Words
    content.Participants = make(map[int][]string)
    for _,part := range(val.Participants) {
      content.Participants[pIDs[pIdx]]=part
      pIdx++
    }
    newESD.Label[eIDs[idx]]=content
    newESD.Tau[eIDs[idx]]=1
    idx++
  }
  for idx:=0; idx<numTop-1;idx++ {
    newESD.V[idx]=rand.Intn(numTop-idx)
  }
  newESD.Init()
  return
}

func GetModel() *Model {
  // 0:boil		1:add		2:serve
  // 0:pasta		1:salt		2:water
  model := new(Model)
  model.numESDs = 20
  model.eventVocabulary = 5
  model.participantVocabulary = 3
  model.eventtype_histogram = Histogram{15,20,20}
  model.participanttype_histogram = Histogram{20,20,15}
  model.participanttype_eventtype_histogram = map[int]Histogram{0:Histogram{0,5,15}, 1:Histogram{0,15,5}, 2:Histogram{15,0,0}}
  model.word_eventtype_histogram = map[string]Histogram{"add":Histogram{0,20,0}, "serve":Histogram{0,0,20}, "boil":Histogram{15,0,0}}
  model.word_participanttype_histogram = map[string]Histogram{"pasta":Histogram{20,0,0}, "water":Histogram{0,0,15},"salt":Histogram{0,20,0}}
  model.invcount_histogram= Histogram{0,0}
  return model
}