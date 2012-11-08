package scriptModeling

import "os/exec"
import "io"
import "bufio"
import "strings"
// import "fmt"

var cmd *exec.Cmd
var cmdIn io.WriteCloser
var cmdOut bufio.Reader

func StartOctave() {
    cmd = exec.Command("octave", "-q")
    cmd.Dir = "/local/lea/thesis/sliceSample/"
//     cmd.Dir = "/home/lea/Code/Octave/"        
    cmdIn,_ = cmd.StdinPipe()
    outpipe, _ := cmd.StdoutPipe()
    cmdOut = *bufio.NewReader(outpipe)
    cmd.Start()
    return
}

func getSliceSampler(args []string) []byte {
  octPrint := ";"+`printf("%f\n",samples)`+"\n"
  parameter := strings.Join(args, ", ")
  command := "samples = slicesample("+parameter+`)`+octPrint
  return []byte(command)
}