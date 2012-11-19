package scriptModeling

import "scriptIO"
import "bytes"
import "os"
import "io/ioutil"
import "fmt"
import "encoding/gob"
import "encoding/xml"
import "path"
import "strings"
import "math/rand"

func GetCorpus (xmlDir, scenario string) (corpus Corpus) {
  
  VList := make([]string, 1000)
  POSList := make([]string, 1000)
  vocabulary = VocabMap{map[string]int{}, map[int]string{}, VList, POSList}
  contents,_ := ioutil.ReadDir(xmlDir)
//   corpus := Corpus{}
  for _, file := range(contents) {
    scenarios := ReadScenarios(path.Join(xmlDir, file.Name()))
    for _,scenario := range(scenarios.Script) {
      esd := createESD(scenario)
      corpus = append(corpus, &esd)
    }
  }
  vocabulary.VList = vocabulary.VList[:vocabIdx]
  vocabulary.POSList = vocabulary.POSList[:vocabIdx]
  fmt.Println(vocabulary.VList, len(vocabulary.VList), "\n==================================\n")
  fmt.Println(vocabulary.POSList, len(vocabulary.POSList), "\n==================================\n")
  fmt.Println(vocabulary.Vtoi)
  vocabulary.Store("/local/lea/thesis/data/corpus/vocab/"+scenario+".bin")
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
      eWordIDs := vocabulary.add(eWords, "v")
      esd.EventLabel[eventLabelIdx]=eIDs[eventLabelIdx]
      esd.Tau[eIDs[eventLabelIdx]]=1
      // generate participant labels
      tmpPtao = [numPar]int{}
      pIDs := rand.Perm(numPar)[:len(event.Participants)]
      esd.Label[eIDs[eventLabelIdx]] = Content{eWordIDs, map[int][]int{}, tmpPtao}
      for pIdx, part := range(event.Participants) {
	pWords := preProcess(strings.Split(part.Text, " "))
	if len(pWords) > 0  && pWords[0] != "imp_protagonist" {
	  pWordIDs := vocabulary.add(pWords, "n")
	  esd.Label[eIDs[eventLabelIdx]].Participants[pIDs[pIdx]] = pWordIDs
	}
      }
      for key, _ := range(esd.Label[eIDs[eventLabelIdx]].Participants) {
	tmpPtao[key]=1
      }
      esd.Label[eIDs[eventLabelIdx]] = Content{esd.Label[eIDs[eventLabelIdx]].Words, esd.Label[eIDs[eventLabelIdx]].Participants, tmpPtao}
    }
    eventLabelIdx++
  }
  esd.EventLabel = esd.EventLabel[:eventLabelIdx]
  // generate ordering under word-order constraint
  fmt.Println(esd.EventLabel, esd.EventLabel)
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
  stopWordList := []string{".", "a", "able", "about", "across", "after", "all", "almost", "also", "am", "among", "an", "and", "any", "are", "as", "at", "be", "because", "been", "but", "by", "can", "cannot", "could", "dear", "did", "do", "does", "either", "else", "ever", "every", "for", "from", "got", "had", "has", "have", "he", "her", "hers", "him", "his", "how", "however", "i", "if", "in", "into", "is", "it", "its", "just", "least", "let", "like", "likely", "may", "me", "might", "most", "must", "my", "neither", "no", "nor", "not", "of", "off", "often", "on", "only", "or", "other", "our", "own", "rather", "she", "should", "since", "so", "some", "than", "that", "the", "their", "them", "then", "there", "these", "they", "this", "tis", "to", "too", "twas", "us", "wants", "was", "we", "were", "what", "when", "where", "which", "while", "who", "whom", "why", "will", "with", "would", "yet", "you", "your", "s", "."}
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

func (corpus Corpus) Store (fname string) {
        b := new(bytes.Buffer)
        enc := gob.NewEncoder(b)
        err := enc.Encode(corpus)
        if err != nil {
                fmt.Println(err)
        }

        fh, eopen := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0666)
        defer fh.Close()
        if eopen != nil {
                fmt.Println(eopen)
        }
        n,e := fh.Write(b.Bytes())
        if e != nil {
                fmt.Println(e)
        }
        fmt.Fprintf(os.Stderr, "%d bytes successfully written to file\n", n)
}

func LoadCorpus (fname string) (Corpus) {
        fh, err := os.Open(fname)
        if err != nil {
                fmt.Println(err)
        }
        corpus := Corpus{}
        dec := gob.NewDecoder(fh)
        err = dec.Decode(&corpus)
        if err != nil {
                fmt.Println(err)
        }
        return corpus
}

func (model Model) Store (fname string) {
        b := new(bytes.Buffer)
        enc := gob.NewEncoder(b)
        err := enc.Encode(model)
        if err != nil {
                fmt.Println(err)
        }

        fh, eopen := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0666)
        defer fh.Close()
        if eopen != nil {
                fmt.Println(eopen)
        }
        n,e := fh.Write(b.Bytes())
        if e != nil {
                fmt.Println(e)
        }
        fmt.Fprintf(os.Stderr, "%d bytes successfully written to file\n", n)
}

func LoadModel (fname string) (Model) {
        fh, err := os.Open(fname)
        if err != nil {
                fmt.Println(err)
        }
        model := Model{}
        dec := gob.NewDecoder(fh)
        err = dec.Decode(&model)
        if err != nil {
                fmt.Println(err)
        }
        return model
}

func (vocabulary VocabMap) Store (fname string) {
        b := new(bytes.Buffer)
        enc := gob.NewEncoder(b)
        err := enc.Encode(vocabulary)
        if err != nil {
                fmt.Println(err)
        }

        fh, eopen := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0666)
        defer fh.Close()
        if eopen != nil {
                fmt.Println(eopen)
        }
        n,e := fh.Write(b.Bytes())
        if e != nil {
                fmt.Println(e)
        }
        fmt.Fprintf(os.Stderr, "%d bytes successfully written to file\n", n)
}

func LoadVocabulary (fname string) (VocabMap) {
        fh, err := os.Open(fname)
        if err != nil {
                fmt.Println(err)
        }
        vocabulary := VocabMap{}
        dec := gob.NewDecoder(fh)
        err = dec.Decode(&vocabulary)
        if err != nil {
                fmt.Println(err)
        }
        return vocabulary
}