package scriptModeling

import "path"
import "strings"
import "math/rand"
import "io/ioutil"
import "scriptIO"
import "fmt"
import "encoding/xml"
import "os"

func GetCorpus (xmlDir, scenario string) (corpus Corpus) {
  VList := make([]string, 1000)
  POSList := make([]string, 1000)
  vocabulary.Dictionary = Dictionary{map[string]int{}, map[int]string{}, VList, POSList}
  contents,_ := ioutil.ReadDir(xmlDir)
//   corpus := Corpus{}
  for _, file := range(contents) {
    scenarios := ReadScenarios(path.Join(xmlDir, file.Name()))
    for _,scenario := range(scenarios.Script) {
      esd := createESD(scenario)
      corpus = append(corpus, &esd)
    }
  }
  
  vocabulary.Dictionary.VList = vocabulary.Dictionary.VList[:vocabIdx]
  fmt.Println(vocabulary.Dictionary.VList)
  fmt.Fprintf(os.Stderr, "printed corpus at iteration %d\n", 1)
  vocabulary.Dictionary.POSList = vocabulary.Dictionary.POSList[:vocabIdx]
  fmt.Println(vocabulary.Dictionary.POSList)
  fmt.Println(vocabulary.Dictionary.Itov)
  vocabulary.CreateCovarianceMatrix()
  fmt.Println(vocabulary.Covariances)
  fmt.Fprintf(os.Stderr, "printed corpus at iteration %d\n", 2)
  vocabulary.CreateEquivalenceClasses()
  fmt.Println(vocabulary.Equivalences)
  fmt.Fprintf(os.Stderr, "printed corpus at iteration %d\n", 3)
  vocabulary.Store("/local/lea/thesis/data/vocabularies/"+scenario+".bin")
  
  return corpus
}


func ReadScenarios(file string) scriptIO.Scripts {
  var sc scriptIO.Scripts
  xmlFile, err := ioutil.ReadFile(file)
  if err != nil {
    panic("Error opening file:")
    fmt.Println(err)
  }
  xml.Unmarshal(xmlFile, &sc)
  return sc
}


func createESD (scenario scriptIO.Script) ESD {
  var eventLabelIdx int
  var esd ESD
  var tmpPtao [numPar]int
  esd.Label = make(map[int]Content)
  esd.EventLabel = make([]int, len(scenario.Item))
  // generate event labels
  eIDs := rand.Perm(numTop)[:len(scenario.Item)]
  for _, event := range(scenario.Item) {
    eWords := preProcess(strings.Split(event.Text, " "))
    if len(eWords) > 0 && eWords[0] != "."  /*|| len(event.Participants)>0*/ {
      eWordIDs := vocabulary.Dictionary.add(eWords, "v")
      esd.EventLabel[eventLabelIdx]=eIDs[eventLabelIdx]
      esd.Tau[eIDs[eventLabelIdx]]=1
      // generate participant labels
      tmpPtao = [numPar]int{}
      pIDs := rand.Perm(numPar)[:len(event.Participants)]
      esd.Label[eIDs[eventLabelIdx]] = Content{eWordIDs, map[int][]int{}, tmpPtao}
      for pIdx, part := range(event.Participants) {
	pWords := preProcess(strings.Split(part.Text, " "))
	if len(pWords) > 0  && pWords[0] != "imp_protagonist" {
	  pWordIDs := vocabulary.Dictionary.add(pWords, "n")
	  esd.Label[eIDs[eventLabelIdx]].Participants[pIDs[pIdx]] = pWordIDs
	}
      }
      for key, _ := range(esd.Label[eIDs[eventLabelIdx]].Participants) {
	tmpPtao[key]=1
      }
      esd.Label[eIDs[eventLabelIdx]] = Content{esd.Label[eIDs[eventLabelIdx]].Words, esd.Label[eIDs[eventLabelIdx]].Participants, tmpPtao}
      eventLabelIdx++
    }
  }
  esd.EventLabel = esd.EventLabel[:eventLabelIdx]
  // generate ordering under word-order constraint
  newPi := createOrdering(esd.EventLabel)
  for idx,el := range(newPi) {
    esd.Pi[idx]=el
  }
  esd.ComputeV()
  esd.Length=len(esd.Label)
  return esd
}


//function for stopwordremoval, trimming
func preProcess(full []string) []string {
  stopWordList := []string{".", "a", "able", "about", "across", "after", "all", "almost", "also", "am", "among", "an", "and", "any", "are", "as", "at", "be", "because", "been", "but", "by", "can", "cannot", "could", "dear", "did", "do", "does", "either", "else", "ever", "every", "for", "from", "got", "had", "has", "have", "he", "her", "hers", "him", "his", "how", "however", "i", "if", "in", "into", "is", "it", "its", "just", "least", "like", "likely", "may", "me", "might", "most", "must", "my", "neither", "no", "nor", "not", "of", "off", "often", "on", "only", "or", "other", "our", "own", "rather", "she", "should", "since", "so", "some", "than", "that", "the", "their", "them", "then", "there", "these", "they", "this", "tis", "to", "too", "twas", "us", "wants", "was", "we", "were", "what", "when", "where", "which", "while", "who", "whom", "why", "will", "with", "would", "yet", "you", "your", "s", "."}
  clean := make([]string, len(full))
  var idx int
  var add bool
  for _, word := range(full) {
    add=true
    for _, stop := range(stopWordList) {
      if word == stop {
	add=false
      }
    }
    if add==true {
      clean[idx]=word
      idx++
    }
  }
  clean = clean[:idx]
  for idx, _ := range(clean) {
    clean[idx] = strings.Trim(clean[idx], `.,!?'"-:(){}[]#$@%^&*_+=`)
  }
  return clean
}
