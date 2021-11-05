package main

import (
   "fmt"
   "log"
	"net/http"
	"os"
   "mesa-state-monitor/pkg"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello World!</h1>"))
}

func main() {

   // check if this is a binary or single star evolution
   if len(os.Args) < 2 {log.Fatal("need star filename (and binary if present)")}

   // get star history fname as first argument & binary as second
   var is_binary_evolution bool
   starfilePath := os.Args[1]
   if len(os.Args) == 3 {
      is_binary_evolution = true
      binaryfilePath := os.Args[2]
      fmt.Println(binaryfilePath)
   }

   // create struct of MESAstar_info
   info := new(read_file.MESAstar_info)
   info.History_name = starfilePath
   read_file.Grab_star_header(starfilePath, info)
   read_file.Grab_star_run_info(starfilePath, info)

   fmt.Printf("MESAstar_info: %+v\n", *info)

   if (is_binary_evolution) {
      BInfo := new(read_file.MESAbinary_info)
      BInfo.History_name = ""
   }

   // set port number
   port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

   // open http server
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)

}
