package web

import (
   "html/template"
   "log"
   "net/http"
   "os"

   "mesa-state-monitor/pkg/io"
   "mesa-state-monitor/pkg/utils"
)


type WrapperStruct struct {
   BinaryFilename string
   Star1Filename string
   Star2Filename string
}


var tpl_mesastar = template.Must(template.ParseFiles("web/index-mesastar.html"))
var tpl_mesabinary = template.Must(template.ParseFiles("web/index-mesabinary.html"))


// index html for MESAstar 
func (ws WrapperStruct) indexHandler_mesastar(w http.ResponseWriter, r *http.Request) {

   log.Printf("setting up index handler for MESAstar")

   starfilePath := ws.Star1Filename

   // create struct of MESAstar_info
   log.Printf("loading struct with MESAstar info")
   info := new(io.MESAstar_info)
   info.History_name = starfilePath
   io.Grab_star_header(starfilePath, info)
   io.Grab_star_run_info(starfilePath, info)
   info.Evol_state = utils.SetEvolutionaryStage(info.Mass, info.Center_h1, info.Center_he4, info.Log_T_cntr)

   tpl_mesastar.Execute(w, info)

   log.Printf("refreshing MESAstar info")

}


// index html for MESAbinary
func (ws WrapperStruct) indexHandler_mesabinary(w http.ResponseWriter, r *http.Request) {

   log.Printf("setting up index handler for MESAstar")

   binaryfilePath := ws.BinaryFilename

   // create struct of MESAbinary_info   
   log.Printf("loading struct with MESAbinary info")
   binfo := new(io.MESAbinary_info)
   binfo.History_name = ""
   binfo.MT_case = "none"
   io.Grab_binary_header(binaryfilePath, binfo)
   io.Grab_binary_run_info(binaryfilePath, binfo)
   // binfo.MT_case = utils.SetMTCase(binfo.Radius1, binfo.RLobe1, info.EvolState)

   tpl_mesabinary.Execute(w, nil)
}


func Start(isBinaryEvolution bool, handlers *WrapperStruct) {

   log.Printf("starting web server")

   // set port number for the webpage
   port := os.Getenv("PORT")
	if port == "" {port = "3000"}
   log.Printf("setting PORT to: %s", port)

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
