package main

import (
   "fmt"
   "html/template"
   "log"
	"net/http"
	"os"

   "mesa-state-monitor/pkg/io"
   "mesa-state-monitor/pkg/utils"
)

var tpl_mesastar = template.Must(template.ParseFiles("web/index-mesastar.html"))
var tpl_mesabinary = template.Must(template.ParseFiles("web/index-mesabinary.html"))


type wrapperStruct struct {
   binaryFilename string
   star1Filename string
   star2Filename string
}


// index html for MESAstar 
func (ws wrapperStruct) indexHandler_mesastar(w http.ResponseWriter, r *http.Request) {

   starfilePath := ws.star1Filename

   // create struct of MESAstar_info
   info := new(io.MESAstar_info)
   info.History_name = starfilePath
   io.Grab_star_header(starfilePath, info)
   io.Grab_star_run_info(starfilePath, info)
   info.Evol_state = utils.SetEvolutionaryStage(info.Mass, info.Center_h1, info.Center_he4, info.Log_T_cntr)

   tpl_mesastar.Execute(w, info)
}


// index html for MESAbinary
func (ws wrapperStruct) indexHandler_mesabinary(w http.ResponseWriter, r *http.Request) {

   binaryfilePath := ws.binaryFilename

   // create struct of MESAbinary_info   
   binfo := new(io.MESAbinary_info)
   binfo.History_name = ""
   binfo.MT_case = "none"
   io.Grab_binary_header(binaryfilePath, binfo)
   io.Grab_binary_run_info(binaryfilePath, binfo)
   // binfo.MT_case = utils.SetMTCase(binfo.Radius1, binfo.RLobe1, info.EvolState)

   tpl_mesabinary.Execute(w, nil)
}


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
   handlers := wrapperStruct{
         binaryFilename: binaryLogName,
         star1Filename: star1LogName,
         star2Filename: star2LogName}


   // set port number for the webpage
   port := os.Getenv("PORT")
	if port == "" {port = "3000"}

   // configs
	mux := http.NewServeMux()
   fs := http.FileServer(http.Dir("./web/assets/"))
   mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

   // serve
   if (isBinaryEvolution) {

      // open http server with template of MESAstar (single evolution)
      mux.HandleFunc("/", handlers.indexHandler_mesabinary)

   } else {

      // open http server with template of MESAbinary (binary evolution)
      mux.HandleFunc("/", handlers.indexHandler_mesastar)

   }

   http.ListenAndServe(":"+port, mux)
}
