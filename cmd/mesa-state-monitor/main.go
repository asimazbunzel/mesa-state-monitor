package main

import (
   // "fmt"
   "html/template"
   "log"
	"net/http"
	"os"
   "mesa-state-monitor/pkg"
)

var tpl_mesastar = template.Must(template.ParseFiles("web/index-mesastar.html"))
var tpl_mesabinary = template.Must(template.ParseFiles("web/index-mesabinary.html"))


type wrapperStruct struct {
   filename string
}


func (ws wrapperStruct) indexHandler_mesastar(w http.ResponseWriter, r *http.Request) {

   starfilePath := ws.filename

   // create struct of MESAstar_info
   info := new(read_file.MESAstar_info)
   info.History_name = starfilePath
   read_file.Grab_star_header(starfilePath, info)
   read_file.Grab_star_run_info(starfilePath, info)

   tpl_mesastar.Execute(w, info)
}

func (ws wrapperStruct) indexHandler_mesabinary(w http.ResponseWriter, r *http.Request) {

   binaryfilePath := ws.filename

   // create struct of MESAbinary_info   
   binfo := new(read_file.MESAbinary_info)
   binfo.History_name = ""
   binfo.MT_case = "none"
   read_file.Grab_binary_header(binaryfilePath, binfo)
   read_file.Grab_binary_run_info(binaryfilePath, binfo)

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
