package scriptModeling

import "math/rand"
// import "fmt"

func (esd *ESD) hasParticipants() bool {
  for _,event := range(esd.Participants.Label) {
    if len(event) > 0 {
      return true
    }
  }
  return false
}

func pick_participant(esd ESD) [2]int {
  //randomly select the participant we want to resample
  event := rand.Intn(len(esd.Participants.Label))
  // pick an event (which has participants)
  for len(esd.Participants.Label[event]) == 0 {
    event = rand.Intn(len(esd.Participants.Label))
  }
  // pick a participant
  participant := rand.Intn(len(esd.Participants.Label[event]))
  return [2]int{event, esd.Participants.Label[event][participant]}
}
  
func getAlternatives(participant int, label []int) []int {
  // Get alternative participant types ; TODO: ugly function!!
  var add bool
  var idx int
  alts := make([]int, numPar-len(label))
  alts[0] = participant
  for ii:=0 ; ii<numPar ; ii++ {
    if ii!= alts[0] {
      add=true
      for _,el := range(label) {
	if el==ii {
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