package scriptModeling

import "bytes"
import "os"
import "fmt"
import "encoding/gob"

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

func (vocabulary Vocabulary) Store (fname string) {
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

func LoadVocabulary (fname string) (Vocabulary) {
        fh, err := os.Open(fname)
        if err != nil {
                fmt.Println(err)
        }
        vocabulary := Vocabulary{}
        dec := gob.NewDecoder(fh)
        err = dec.Decode(&vocabulary)
        if err != nil {
                fmt.Println(err)
        }
        return vocabulary
}