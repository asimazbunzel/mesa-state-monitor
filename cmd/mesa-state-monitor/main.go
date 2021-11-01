package main

import (
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello World!</h1>"))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}


// func main() {
//
//    // check if this is a binary or single star evolution
//    if len(os.Args) < 2 {log.Fatal("need star filename (and binary if present)")}
//
//    // get star history fname as first argument & binary as second
//    var is_binary_evolution bool
//    starfilePath := os.Args[1]
//    if len(os.Args) == 3 {
//       is_binary_evolution = true
//       binaryfilePath := os.Args[2]
//       fmt.Println(binaryfilePath)
//    }
//
//    // create struct of MESAstar_info
//    Info := new(MESAstar_info)
//    Info.history_name = starfilePath
//    grab_star_header(starfilePath, Info)
//
//    // now get info on the star using the row containing names for data columns
//    // and the last row written by the MESA code
//    grab_star_run_info(starfilePath, Info)
//
//    fmt.Printf("MESAstar_info: %+v\n", *Info)
//
//    if (is_binary_evolution) {
//       BInfo := new(MESAbinary_info)
//       BInfo.history_name = ""
//    }
// }
