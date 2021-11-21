package main

import (
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
   filename string
}


// index html for MESAstar 
func (ws wrapperStruct) indexHandler_mesastar(w http.ResponseWriter, r *http.Request) {

   starfilePath := ws.filename

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

   binaryfilePath := ws.filename

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

   // check if this is a binary or single star evolution
   if len(os.Args) < 2 {log.Fatal("need star filename (and binary if present)")}

   // get star history fname as first argument & binary as second
   var is_binary_evolution bool
   starfilePath := os.Args[1]
   binaryfilePath := ""
   if len(os.Args) == 3 {
      is_binary_evolution = true
      binaryfilePath = os.Args[2]
   }


   // set port number
   port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

   // configs
	mux := http.NewServeMux()
   fs := http.FileServer(http.Dir("./web/assets/"))
   mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

   // serve
   if (is_binary_evolution) {

      handlers := wrapperStruct{filename: binaryfilePath}

      // open http server with template of MESAstar (single evolution)
      mux.HandleFunc("/", handlers.indexHandler_mesastar)

   } else {

      handlers := wrapperStruct{filename: starfilePath}

      // open http server with template of MESAbinary (binary evolution)
      mux.HandleFunc("/", handlers.indexHandler_mesastar)

   }

   http.ListenAndServe(":"+port, mux)
}
