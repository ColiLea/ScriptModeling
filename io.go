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
// import "stemmer"

func GetCorpus (xmlDir string) (Corpus) {
  VList := make([]string, 100000000)
  vocabulary = vocabMap{map[string]int{}, map[int]string{}, VList}
  contents,_ := ioutil.ReadDir(xmlDir)
  corpus := Corpus{}
  for _, file := range(contents) {
    scenarios := ReadScenarios(path.Join(xmlDir, file.Name()))
    for _,scenario := range(scenarios.Script) {
      esd := createESD(scenario)
      corpus = append(corpus, &esd)
    }
  }
  vocabulary.VList = vocabulary.VList[:vocabIdx]
  fmt.Println(vocabulary.VList, len(vocabulary.VList), "\n==================================\n")
  fmt.Println(vocabulary.vtoi, len(vocabulary.vtoi), "\n==================================\n")
  fmt.Println(vocabulary.itov, len(vocabulary.itov), "\n==================================\n")
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
  var esd ESD
  var tmpPtao [numPar]int
  esd.Label = make(map[int]Content)
  esd.EventLabel = make([]int, len(scenario.Item))
  // generate event labels
  eIDs := rand.Perm(numTop)[:len(scenario.Item)]
  for idx, event := range(scenario.Item) {
    eWords := preProcess(strings.Split(event.Text, " "))
    vocabulary.add(eWords)
    if len(eWords) > 0 || len(event.Participants)>0 {
      esd.EventLabel[idx]=eIDs[idx]
      esd.Tau[eIDs[idx]]=1
      // generate participant labels
      tmpPtao = [numPar]int{}
      pIDs := rand.Perm(numPar)[:len(event.Participants)]
      esd.Label[eIDs[idx]] = Content{getIntRep(eWords), map[int][]int{}, tmpPtao}
      for pIdx, part := range(event.Participants) {
	pWords := preProcess(strings.Split(part.Text, " "))
	vocabulary.add(pWords)
	if len(pWords) > 0 {
	  esd.Label[eIDs[idx]].Participants[pIDs[pIdx]] = getIntRep(pWords)
	}
      }
      for key, _ := range(esd.Label[eIDs[idx]].Participants) {
	tmpPtao[key]=1
      }
      esd.Label[eIDs[idx]] = Content{esd.Label[eIDs[idx]].Words, esd.Label[eIDs[idx]].Participants, tmpPtao}
    }
  }
  // generate ordering under word-order constraint
  newPi := createOrdering(esd.EventLabel)
  for idx,el := range(newPi) {
    esd.Pi[idx]=el
  }
  esd.ComputeV()
  esd.Length=len(esd.Label)
  return esd
}

func getIntRep(words []string) (intRep []int) {
  intRep = make([]int, len(words))
  for idx, word := range(words) {
    intRep[idx] = vocabulary.vtoi[word]
  }
  return
}


//function for stopwordremoval, trimming and stemming (for the latter using external code by ...)
func preProcess(full []string) []string {
  stopWordList := []string{"a", "able", "about", "across", "after", "all", "almost", "also", "am", "among", "an", "and", "any", "are", "as", "at", "be", "because", "been", "but", "by", "can", "cannot", "could", "dear", "did", "do", "does", "either", "else", "ever", "every", "for", "from", "got", "had", "has", "have", "he", "her", "hers", "him", "his", "how", "however", "i", "if", "in", "into", "is", "it", "its", "just", "least", "let", "like", "likely", "may", "me", "might", "most", "must", "my", "neither", "no", "nor", "not", "of", "off", "often", "on", "only", "or", "other", "our", "own", "rather", "she", "should", "since", "so", "some", "than", "that", "the", "their", "them", "then", "there", "these", "they", "this", "tis", "to", "too", "twas", "us", "wants", "was", "we", "were", "what", "when", "where", "which", "while", "who", "whom", "why", "will", "with", "would", "yet", "you", "your", "s", "."}
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
//     clean[idx] = string(stemmer.Stem([]byte(clean[idx])))
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