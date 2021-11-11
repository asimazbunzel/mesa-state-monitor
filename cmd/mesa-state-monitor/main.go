package main

import (
   "fmt"
   "html/template"
   "log"
	"net/http"
	"os"
   "mesa-state-monitor/pkg"
)

var tpl_mesastar = template.Must(template.ParseFiles("web/index-mesastar.html"))
var tpl_mesabinary = template.Must(template.ParseFiles("web/index-mesabinary.html"))


func indexHandler_mesastar(w http.ResponseWriter, r *http.Request) {
   tpl_mesastar.Execute(w, nil)
}

func indexHandler_mesabinary(w http.ResponseWriter, r *http.Request) {
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
      fmt.Println(binaryfilePath)
   }

   // create struct of MESAstar_info
   info := new(read_file.MESAstar_info)
   info.History_name = starfilePath
   read_file.Grab_star_header(starfilePath, info)
   read_file.Grab_star_run_info(starfilePath, info)

   fmt.Printf("MESAstar_info: %+v\n", *info)

   if (is_binary_evolution) {
      binfo := new(read_file.MESAbinary_info)
      binfo.History_name = binaryfilePath
      binfo.MT_case = "none"
      read_file.Grab_binary_header(binaryfilePath, binfo)
      read_file.Grab_binary_run_info(binaryfilePath, binfo)

      fmt.Printf("MESAbinary_info: %+v\n", *binfo)
   }

   // set port number
   port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

   fs := http.FileServer(http.Dir("./web/assets/"))

   // open http server
	mux := http.NewServeMux()
   mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

   if (is_binary_evolution) {
      mux.HandleFunc("/", indexHandler_mesabinary)
   } else {
      mux.HandleFunc("/", indexHandler_mesastar)
   }

	http.ListenAndServe(":"+port, mux)

}
