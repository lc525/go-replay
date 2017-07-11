package main

import (
  "flag"
  "fmt"
  "io"
  "log"
  "time"

  "github.com/lc525/go-cadets"
  "gopkg.in/cheggaaa/pb.v1"
)

var in = flag.String("in", "", "input trace file to process")
var speedf = flag.Float64("speed", 1.0, "replay speed")

/*func refreshBar(tr cadets.TraceOps, bar *pb.ProgressBar) {
 *  current, _ := tr.GetCurrentOffset()
 *  bar.Set64(current)
 *}*/
func refreshBar(bar *pb.ProgressBar, count *int) {
  bar.Set64(int64(*count))
}

func main() {
  flag.Parse()
  var count int

  tr, err := cadets.OpenTraceFile(*in)

  if err != nil {
    log.Fatal("cadets trace:", err)
    return
  }
  defer tr.Close()


  var rcfg = cadets.ReplayConfig {
    SpeedFactor: float32(*speedf),
    BufCtlType: cadets.Manual,
    BufSizeOrder: 22,
    MaxConsecutiveDelta: 1e9,
    ChunkDeltaThresh: 3e9,  // 3 s in ns
  }

  fmt.Println("Buffering...")
  trr := cadets.NewTraceReplay(tr, rcfg)
  trr.Start()

  fmt.Println("Replaying trace...", *in)
  //traceSize, _ := tr.GetTraceSize()
  bar := pb.New64(8383490).SetUnits(pb.U_NO).SetRefreshRate(time.Second * 1)
  bar.ShowSpeed = true
  bar.Start()

  ticker := time.NewTicker(time.Second)
  go func() {
    for _ = range ticker.C {
      refreshBar(bar, &count)
    }
  }()

  for {
    _, err := trr.Read()
    if err != nil {
      if err != io.EOF {
        log.Println("trace error:", err)
      }
      break;
    }

    count++
  }
  bar.Finish()
  ticker.Stop()
  fmt.Println("Trace events:", count)

  return
}
