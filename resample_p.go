 package scriptModeling

import "math/rand"
// import "fmt"

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