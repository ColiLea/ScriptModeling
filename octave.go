package scriptModeling

import "os/exec"
import "io"
import "bufio"
// import "fmt"

var cmd *exec.Cmd
var cmdIn io.WriteCloser
var cmdOut bufio.Reader

var pCmd *exec.Cmd
var pCmdIn io.WriteCloser
var pCmdOut bufio.Reader

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