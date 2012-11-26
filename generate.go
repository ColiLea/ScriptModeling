package scriptModeling

import "math"
import "math/rand"
// import "fmt"

func (model *Model) Generate(jPrior, lmPrior float64) *ESD {
//   var numPWords int
  var wList []int
  const modelTop int = 3
  const modelPar int = 3
  vocabulary.Dictionary = Dictionary{map[string]int{"add":1, "serve":2, "boil":0, "water":5, "salt":6, "pasta":3, "noodle":4, "cook":7}, map[int]string{7:"cook", 0:"boil", 1:"add", 2:"serve", 3:"pasta", 4:"noodle", 5:"water", 6:"salt"}, []string{"boil","add","serve","pasta","noodle","water","salt", "cook"}, []string{"v", "v", "v", "n", "n", "n", "n", "v"}}
  vocabulary.Covariances = new(CovarianceStruct)
  vocabulary.Covariances.Matrix = [][]float64{[]float64{1, 0.030303030303, 0, 0, 0, 0, 0, 0.030303030303}, []float64{0.030303030303, 1, 0, 0, 0, 0, 0, 0.0166666666667}, []float64{0, 0, 1, 0, 0, 0, 0, 0}, []float64{0, 0, 0, 1, 0.114285714286, 0, 0, 0}, []float64{0, 0, 0, 0.114285714286, 1, 0, 0, 0}, []float64{0, 0, 0, 0, 0, 1, 0, 0}, []float64{0, 0, 0, 0, 0, 0, 1, 0}, []float64{0.030303030303, 0.0166666666667, 0, 0, 0, 0, 0, 1}}
  vocabulary.Covariances.Inverse = [][]float64{[]float64{1.00180970909, -0.0298601999729, 0, 0, 0, 0, 0, -0.0298601999729}, []float64{-0.0298601999729, 1.00116787582, 0, 0, 0, 0, 0, -0.0157812767191}, []float64{0, 0, 1, 0, 0, 0, 0, 0}, []float64{0, 0, 0, 1.01323407775, -0.115798180314, 0, 0, 0}, []float64{0, 0, 0, -0.115798180314, 1.01323407775, 0, 0, 0}, []float64{0, 0, 0, 0, 0, 1, 0, 0}, []float64{0, 0, 0, 0, 0, 0, 1, 0}, []float64{-0.0298601999729, -0.0157812767191, 0, 0, 0, 0, 0, 1.00116787582}}
  vocabulary.Equivalences=map[int][]int{5:[]int{5}, 0:[]int{0, 1, 7}, 1:[]int{0, 1}, 7:[]int{0, 7}, 2:[]int{2}, 6:[]int{6}, 4:[]int{3, 4}, 3:[]int{3, 4}}
  rho := [2]float64{5.9, 5.9}
  tmpPtau := [numPar]int{}
  esd := new(ESD)
  esd.Label = Label{}
  //Generate Eventtypes
  for jj := 0 ; jj<modelTop ; jj++ {
    jPos := (float64(model.Eventtype_histogram[jj])+jPrior)/(float64(model.NumESDs)+jPrior)
    jNeg := (float64(model.NumESDs-model.Eventtype_histogram[jj])+jPrior)/(float64(model.NumESDs)+jPrior)
    esd.Tau[jj]=sample([]float64{jNeg,jPos})
    if esd.Tau[jj] ==1 {
      esd.Length++
      esd.Label[jj]=Content{[]int{}, map[int][]int{}, tmpPtau}
      //Generate Participants
      for ii:=0 ; ii<modelPar ; ii++ {
	pPos := (float64(model.Participanttype_eventtype_histogram[ii][jj])+jPrior)/(float64(model.Eventtype_histogram[jj])+jPrior)
	pNeg := (float64(model.Eventtype_histogram[jj]-model.Participanttype_eventtype_histogram[ii][jj])+jPrior)/(float64(model.Eventtype_histogram[jj])+jPrior)
	pp := sample([]float64{pNeg,pPos})
	if pp == 1 {
	  esd.Label[jj].Participants[ii] = []int{}
	}
      }
      for pIdx,_ := range(esd.Label[jj].Participants) {
	tmpPtau[pIdx]=1
      }
    }
    //Generate ordering
    vDist := make([]float64, modelTop-jj)
    for vv:=0 ; vv<modelTop-jj && jj<modelTop-1 ; vv++ {
      vDist[vv] = math.Exp(-rho[jj]*float64(vv+1))
    }
    if jj<modelTop-1 {
      esd.V[jj]= getAccumulativeSample(vDist)
    }
    esd.ComputePi()
    esd.ComputeZ()
  }
  //Generate words
  for eID, event := range(esd.Label) {
    wList = []int{}
//     numEWords := rand.Intn(2)+1
//     for i:=0 ; i<numEWords ; i++ {
      wDist := make([]float64, model.EventVocabulary)
      words := make([]int, model.EventVocabulary)
      idx:=0
      for term,dist := range(model.Word_eventtype_histogram) {
	words[idx]=term
	wDist[idx]=(float64(dist[eID])+lmPrior)
	idx++
      }
      ww := getAccumulativeSample(wDist)
      wList = append(wList, words[ww])
//     }
    esd.Label[eID] = Content{wList, esd.Label[eID].Participants, esd.Label[eID].Tau}
    for pID, _ := range(event.Participants) {
      wList = []int{}
//       if pID == 2 {
// 	numPWords = rand.Intn(2)+1
//       } else {
// 	numPWords =1
//       }
//       for j:=0 ; j<numPWords ; j++ {
	idx:=0
	wDist := make([]float64, model.ParticipantVocabulary)
	words := make([]int, model.ParticipantVocabulary)
	for term,dist := range(model.Word_participanttype_histogram) {
	  words[idx]=term
	  wDist[idx]=(float64(dist[pID])+lmPrior)
	  idx++
	}
	pp := getAccumulativeSample(wDist)
	wList = append(wList, words[pp])
//       }
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
    tmpPtau := [numPar]int{}
    pIDs := rand.Perm(numPar)[:len(val.Participants)]
    pIdx:=0
    content := Content{}
    content.Words = val.Words
    content.Participants = make(map[int][]int)
    for _,part := range(val.Participants) {
      content.Participants[pIDs[pIdx]]=part
      tmpPtau[pIDs[pIdx]]=1
      pIdx++
    }
    content.Tau = tmpPtau
    newESD.Label[eIDs[idx]]=content
    newESD.Tau[eIDs[idx]]=1
    idx++
  }
  newESD.EventLabel = make([]int, len(newESD.Label))
  for oIdx, oID := range(esd.EventLabel) {
    oldW := esd.Label[oID].Words
    for nID, newE := range(newESD.Label) {
      if Compare(oldW, newE.Words) == true && isIn(nID, newESD.EventLabel[:oIdx]) == false {
	newESD.EventLabel[oIdx]=nID
      }
    }
  }
  if len(newESD.EventLabel) == numTop {
    // if all types are realized: pi==z
    for idx,val := range(newESD.EventLabel) {
      newESD.Pi[idx] = val
    }
  } else {
    newPi := createOrdering(newESD.EventLabel)
    for idx,el := range(newPi) {
      newESD.Pi[idx]=el
    }
  }
  newESD.ComputeV()
  newESD.Length = len(newESD.Label)
  return
}



func createOrdering(label []int) []int {
  others := getDiffList(label)
  pi := make([]int, len(label))
    for idx,val := range(label) {
      pi[idx] = val
    }
    for _, el := range(others) {
      insert := rand.Intn(len(pi))
      pi = append(pi[:insert], append([]int{el},pi[insert:]...)...)
    }
  return pi
}

func getDiffList(eLabel []int) (dList []int) {
  for ii:=0 ; ii<numTop ; ii++ {
    found := false
    for _,v := range(eLabel) {
      if ii==v {
	found = true
      }
    }
    if found==false {
      dList = append(dList,ii)
    }
  }
  return
}


func GetModel() *Model {
  // 0:boil,heat	1:add		2:serve
  // 0:pasta		1:salt		2:water
  // "boil":0, "add":1, "serve":2, "pasta":3, "noodle":4, "water":5, "salt":6, "cook":7
  model := new(Model)
  model.NumESDs = 20
  model.EventVocabulary = 4
  model.ParticipantVocabulary = 4
  model.Eventtype_histogram = Histogram{20,15,20}
  model.Participanttype_histogram = Histogram{20,15,20}
  model.Participanttype_eventtype_histogram = map[int]Histogram{0:Histogram{0,0,20}, 1:Histogram{0,15,0}, 2:Histogram{20,0,0}}
  model.Word_eventtype_histogram = map[int]Histogram{1:Histogram{0,15,0}, 2:Histogram{0,0,20}, 0:Histogram{12,0,0}, 7:Histogram{8,0,0}}
  model.Word_participanttype_histogram = map[int]Histogram{3:Histogram{10,0,0}, 4:Histogram{10,0,0}, 5:Histogram{0,0,20},6:Histogram{0,15,0}}
  model.Invcount_histogram= Histogram{0,0}
  model.Rho = []float64{0.0,0.0}
  return model
}