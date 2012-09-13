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


func GetCorpus (xmlDir string) (Corpus) {
  contents,_ := ioutil.ReadDir(xmlDir)
  corpus := Corpus{}
  for _, file := range(contents) {
    scenarios := ReadScenarios(path.Join(xmlDir, file.Name()))
    for _,scenario := range(scenarios.Script) {
      esd := createESD(scenario)
      corpus = append(corpus, &esd)
    }
  }
  return corpus
}


func ReadScenarios(file string) scriptIO.Scripts {
  var sc scriptIO.Scripts
  xmlFile, err := ioutil.ReadFile(file)
  if err != nil {
    fmt.Println(err)
    panic("Error opening file:")
  }
  xml.Unmarshal(xmlFile, &sc)
  return sc
}


func createESD (scenario scriptIO.Script) ESD {
  var esd ESD
  esd.Label = make(map[int]Content)
  eIDs := rand.Perm(numTop-1)[:len(scenario.Item)]
  for idx, event := range(scenario.Item) {
    esd.Label[eIDs[idx]] = Content{strings.Split(event.Text, " "), map[int][]string{}}
    esd.Tau[eIDs[idx]]=1
    pIDs := rand.Perm(numPar-1)[:len(event.Participants)]
    for pIdx, part := range(event.Participants) {
      esd.Label[eIDs[idx]].Participants[pIDs[pIdx]] = strings.Split(part.Text, " ")
    }
  }
  for idx:=0; idx<numTop-1;idx++ {
    esd.V[idx]=rand.Intn(numTop-idx)
  }
  esd.Length=len(esd.Label)
  return esd
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