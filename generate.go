package scriptModeling

import "math"
import "math/rand"
// import "fmt"

func (model *Model) Generate(jPrior, lmPrior float64) *ESD {
  var wList []string
  modelTop := 3
  modelPar := 4
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
    }
    esd.ComputePi()
    esd.ComputeZ()
  }
  //Generate words
  for eID, event := range(esd.Label) {
    wList = []string{}
    numEWords := rand.Intn(2)+1
    for i:=0 ; i<numEWords ; i++ {
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
    }
    esd.Label[eID] = Content{wList, esd.Label[eID].Participants}
    for pID, _ := range(event.Participants) {
      idx:=0
      wList := []string{}
      wDist := make([]float64, model.participantVocabulary)
      words := make([]string, model.participantVocabulary)
      for term,dist := range(model.word_participanttype_histogram) {
	words[idx]=term
	wDist[idx]=(float64(dist[pID])+lmPrior)
	idx++
      }
      pp := getAccumulativeSample(wDist)
      wList = append(wList, words[pp])
      esd.Label[eID].Participants[pID]=wList
    }
  }
  return esd
}

func Randomize(esd ESD) (newESD ESD) {
  var oldZPos int
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
  newESD.EventLabel = make([]int, len(newESD.Label))
  for nidx, nval := range(newESD.Label) {
    for idx, val := range(esd.Label) {
      if nval.Words[0] == val.Words[0] {
	for zIdx,zL := range(esd.EventLabel) {
	  if zL==idx {
	    oldZPos = zIdx
	  }
	}
	newESD.EventLabel[oldZPos]=nidx
      }
    }
  }
  if len(newESD.EventLabel) == numTop {
    // if all types are realized: pi==z
    for idx,val := range(newESD.EventLabel) {
      newESD.Pi[idx] = val
    }
  } else {
    newPi := make([]int, len(newESD.EventLabel))
    for idx,val := range(newESD.EventLabel) {
      newPi[idx] = val
    }
    // get diff
    var others []int
    for ii:=0 ; ii<numTop ; ii++ {
      found := false
      for _,v := range(newESD.EventLabel) {
        if ii==v {
	  found = true
	}
      }
      if found==false {
	others = append(others,ii)
      }
    }
    for _, el := range(others) {
      insert := rand.Intn(len(newPi))
      newPi = append(newPi[:insert], append([]int{el},newPi[insert:]...)...)
    }
    for idx,el := range(newPi) {
      newESD.Pi[idx]=el
    }
  }
  newESD.Length = len(newESD.Label)
//   newESD.Init()
  return
}

func GetModel() *Model {
  // 0:boil,heat	1:add		2:serve
  // 0:pasta		1:salt		2:water
  model := new(Model)
  model.numESDs = 20
  model.eventVocabulary = 7
  model.participantVocabulary = 5
  model.eventtype_histogram = Histogram{15,20,19}
  model.participanttype_histogram = Histogram{19,20,15,10}
  model.participanttype_eventtype_histogram = map[int]Histogram{0:Histogram{0,0,19}, 1:Histogram{0,20,0}, 2:Histogram{15,0,0}, 3:Histogram{10,0,0}}
  model.word_eventtype_histogram = map[string]Histogram{"add":Histogram{0,12,0}, "serve":Histogram{0,0,19}, "hot":Histogram{0,0,5}, "boil":Histogram{10,0,0}, "heat":Histogram{5,0,0}, "quickly":Histogram{5,0,0}, "put":Histogram{0,8,0}}
  model.word_participanttype_histogram = map[string]Histogram{"pot":Histogram{0,0,0,10}, "pasta":Histogram{10,0,0,0}, "noodles":Histogram{9,0,0,0}, "water":Histogram{0,0,15,0},"salt":Histogram{0,20,0,0}}
  model.invcount_histogram= Histogram{0,0}
  model.rho = []float64{0.0,0.0}
  return model
}





//     for idx,val := range(newESD.EventLabel) {
//       newESD.Pi[idx] = val
//     }
//     for i:=0 ; i<numTop ; i++ {
//       found := false
//       for _,el := range newESD.EventLabel {
// 	if i==el {
// 	  found = true
// 	}
//       }
//       if found == false {
// 	newESD.Pi[numTop-1]=i
//       }
//     }
//   }


//   lIdx := 0
//   for idx, label := range(newESD.EventLabel) {
//     lIdx = rand.Intn(numTop-lIdx)+lIdx 
//     newESD.Pi[lIdx]=label
//   }
//   var others []int
//   for ii:=0 ; ii<numTop ; ii++ {
//     found := false
//     for _,v := range(newESD.Pi) {
//       if ii == v {
// 	found=true
//       }
//     }
//     if found==false {
//       others = append(others,ii)
//     }
//   }
//   perm := rand.Perm(others)
//   for _,v := range(perm) {
//     for idx,_ := range(newESD.Pi) {
//       if newESD.Pi[idx]==0 {
// 	newESD.Pi[idx]=v
//       }
//     }
//   }
//   for idx,id := range(newESD.Pi) {
//     if id < numTop-1{
//       for tidx:=0 ; tidx<idx ; tidx++ {
// 	if newESD.Pi[tidx]>id {
// 	  newESD.V[id]++
// 	}
//       }
//     }
//   }