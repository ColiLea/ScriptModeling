package scriptModeling

// import "math"
import "math/rand"
// import "fmt"

// func (model *Model) Generate(jPrior, lmPrior float64) *ESD {
//   var numPWords int
//   var wList []string
//   const modelTop int = 3
//   const modelPar int = 3
//   rho := [2]float64{5.9, 5.9}
//   tmpPtau := [numPar]int{}
//   esd := new(ESD)
//   esd.Label = Label{}
//   //Generate Eventtypes
//   for jj := 0 ; jj<modelTop ; jj++ {
//     jPos := (float64(model.eventtype_histogram[jj])+jPrior)/(float64(model.numESDs)+jPrior)
//     jNeg := (float64(model.numESDs-model.eventtype_histogram[jj])+jPrior)/(float64(model.numESDs)+jPrior)
//     esd.Tau[jj]=sample([]float64{jNeg,jPos})
//     if esd.Tau[jj] ==1 {
//       esd.Length++
//       esd.Label[jj]=Content{[]string{}, map[int][]string{}, tmpPtau}
//       //Generate Participants
//       for ii:=0 ; ii<modelPar ; ii++ {
// 	pPos := (float64(model.participanttype_eventtype_histogram[ii][jj])+jPrior)/(float64(model.participanttype_histogram[ii])+jPrior)
// 	pNeg := (float64(model.participanttype_histogram[ii]-model.participanttype_eventtype_histogram[ii][jj])+jPrior)/(float64(model.participanttype_histogram[ii])+jPrior)
// 	pp := sample([]float64{pNeg,pPos})
// 	if pp == 1 {
// 	  esd.Label[jj].Participants[ii] = []string{}
// 	}
//       }
//       for pIdx,_ := range(esd.Label[jj].Participants) {
// 	tmpPtau[pIdx]=1
//       }
//     }
//     //Generate ordering
//     vDist := make([]float64, modelTop-jj)
//     for vv:=0 ; vv<modelTop-jj && jj<modelTop-1 ; vv++ {
//       vDist[vv] = math.Exp(-rho[jj]*float64(vv+1))
//     }
//     if jj<modelTop-1 {
//       esd.V[jj]= getAccumulativeSample(vDist)
//     }
//     esd.ComputePi()
//     esd.ComputeZ()
//   }
//   //Generate words
//   for eID, event := range(esd.Label) {
//     wList = []string{}
//     numEWords := rand.Intn(2)+1
//     for i:=0 ; i<numEWords ; i++ {
//       wDist := make([]float64, model.eventVocabulary)
//       words := make([]string, model.eventVocabulary)
//       idx:=0
//       for term,dist := range(model.word_eventtype_histogram) {
// 	words[idx]=term
// 	wDist[idx]=(float64(dist[eID])+lmPrior)
// 	idx++
//       }
//       ww := getAccumulativeSample(wDist)
//       wList = append(wList, words[ww])
//     }
//     esd.Label[eID] = Content{wList, esd.Label[eID].Participants, esd.Label[eID].Tau}
//     for pID, _ := range(event.Participants) {
//       wList := []string{}
//       if pID == 2 {
// 	numPWords = rand.Intn(2)+1
//       } else {
// 	numPWords =1
//       }
//       for j:=0 ; j<numPWords ; j++ {
// 	idx:=0
// 	wDist := make([]float64, model.participantVocabulary)
// 	words := make([]string, model.participantVocabulary)
// 	for term,dist := range(model.word_participanttype_histogram) {
// 	  words[idx]=term
// 	  wDist[idx]=(float64(dist[pID])+lmPrior)
// 	  idx++
// 	}
// 	pp := getAccumulativeSample(wDist)
// 	wList = append(wList, words[pp])
//       }
//       esd.Label[eID].Participants[pID]=wList
//     }
//   }
//   return esd
// }
// 
// func Randomize(esd ESD) (newESD ESD) {
//   var oldZPos int
//   idx:=0
//   newESD.Label = make(map[int]Content, len(esd.Label))
//   eIDs := rand.Perm(numTop)[:len(esd.Label)]
//   for _, val := range(esd.Label) {
//     tmpPtau := [numPar]int{}
//     pIDs := rand.Perm(numPar)[:len(val.Participants)]
//     pIdx:=0
//     content := Content{}
//     content.Words = val.Words
//     content.Participants = make(map[int][]string)
//     for _,part := range(val.Participants) {
//       content.Participants[pIDs[pIdx]]=part
//       tmpPtau[pIDs[pIdx]]=1
//       pIdx++
//     }
//     content.Tau = tmpPtau
//     newESD.Label[eIDs[idx]]=content
//     newESD.Tau[eIDs[idx]]=1
//     idx++
//   }
//   newESD.EventLabel = make([]int, len(newESD.Label))
//   for nidx, nval := range(newESD.Label) {
//     for idx, val := range(esd.Label) {
//       if nval.Words[0] == val.Words[0] {
// 	for zIdx,zL := range(esd.EventLabel) {
// 	  if zL==idx {
// 	    oldZPos = zIdx
// 	  }
// 	}
// 	newESD.EventLabel[oldZPos]=nidx
//       }
//     }
//   }
//   if len(newESD.EventLabel) == numTop {
//     // if all types are realized: pi==z
//     for idx,val := range(newESD.EventLabel) {
//       newESD.Pi[idx] = val
//     }
//   } else {
//     newPi := createOrdering(newESD.EventLabel)
//     for idx,el := range(newPi) {
//       newESD.Pi[idx]=el
//     }
//   }
//   newESD.ComputeV()
//   newESD.Length = len(newESD.Label)
//   return
// }
// 
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
  model := new(Model)
  model.numESDs = 20
  model.eventVocabulary = 3
  model.participantVocabulary = 4
  model.eventtype_histogram = Histogram{20,15,20}
  model.participanttype_histogram = Histogram{20,20,20}
  model.participanttype_eventtype_histogram = map[int]Histogram{0:Histogram{0,0,20}, 1:Histogram{0,15,0}, 2:Histogram{20,0,0}}
  model.word_eventtype_histogram = map[int]Histogram{0:Histogram{0,15,0}, 1:Histogram{0,0,20}, 2:Histogram{20,0,0}}
  model.word_participanttype_histogram = map[int]Histogram{3:Histogram{10,0,0}, 4:Histogram{10,0,0}, 5:Histogram{0,0,15},6:Histogram{0,20,0}}
  model.invcount_histogram= Histogram{0,0}
  model.rho = []float64{0.0,0.0}
  return model
}