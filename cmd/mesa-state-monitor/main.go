package main


import (
   "fmt"
   "log"
	"os"

   "mesa-state-monitor/pkg/io"
   "mesa-state-monitor/pkg/web"
)


func main() {

   // logging output to a file
   // so first, create folder in $HOME/.local/share/mesa-state-monitor if don't exist,
   // then configure logging to output to a file inside that folder
   logFolder := fmt.Sprintf("%s/.local/share/%s", os.Getenv("HOME"), "mesa-state-monitor")
   if err := os.MkdirAll(logFolder, os.ModePerm); err != nil {
        log.Fatal(err)
    }

    logFilename := fmt.Sprintf("%s/mesa-state-monitor.log", logFolder)
    f, err := os.OpenFile(logFilename, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
    if err != nil {log.Fatal(err)}

    // remember to close it
    defer f.Close()

    // output to a file.
    log.SetOutput(f)

   // get path with MESA simulation
   log.Printf("MESA root folder with simulation in `%s`", os.Args[1])
   mesaRoot := os.Args[1]

   // find out if this is a binary evolution (using MESAbinary)
   isBinaryEvolution := io.IsBinary(mesaRoot)

   // get LOG names for either single or binary evolutions.
   // in the case of a single evolution, only star1LogName should not be empty
   // for a binary, binaryLogName and star1LogName will not be empty; star2LogName might, if its a
   // star + point-mass simulations
   binaryLogName, star1LogName, star2LogName := io.GetLogNames(mesaRoot, isBinaryEvolution)

   // create the wrapper struct from which the server will be launched
   handlers := web.WrapperStruct{
         BinaryFilename: binaryLogName,
         Star1Filename: star1LogName,
         Star2Filename: star2LogName}

   web.Start(isBinaryEvolution, &handlers)

}
