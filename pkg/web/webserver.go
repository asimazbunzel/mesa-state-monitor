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




// index html for MESAstar 
func (ws WrapperStruct) indexHandler_mesastar(w http.ResponseWriter, r *http.Request) {

   log.Printf("setting up index handler for MESAstar")

   starfilePath := ws.Star1Filename

   // create struct of MESAstar_info
   log.Printf("(re)loading struct with MESAstar info")
   info := new(io.MESAstar_info)
   info.HistoryName = starfilePath
   io.GrabStarHeader(starfilePath, info)
   io.GrabStarRunInfo(starfilePath, info)
   info.EvolState = utils.SetEvolutionaryStage(info.Mass, info.CenterH1, info.CenterHe4, info.LogTcntr)

   tpl_mesastar := template.Must(template.ParseFiles("web/index-mesastar.html"))
   tpl_mesastar.Execute(w, info)

}


// index html for MESAbinary
func (ws WrapperStruct) indexHandler_mesabinary(w http.ResponseWriter, r *http.Request) {

   log.Printf("setting up index handler for MESAbinary")

   binaryfilePath := ws.BinaryFilename
   star1filePath := ws.Star1Filename
   star2filePath := ws.Star2Filename

   // create struct of MESAbinary_info   
   log.Printf("(re)loading struct with MESAbinary info")
   binfo := new(io.MESAbinary_info)
   binfo.HistoryName = ""
   binfo.MTCase = "none"
   io.GrabBinaryHeader(binaryfilePath, binfo)
   io.GrabBinaryRunInfo(binaryfilePath, binfo)

   // create struct of MESAstar_info for star 1
   log.Printf("(re)loading struct with MESAstar info for star 1")
   info1 := new(io.MESAstar_info)
   info1.HistoryName = star1filePath
   io.GrabStarHeader(star1filePath, info1)
   io.GrabStarRunInfo(star1filePath, info1)
   info1.EvolState = utils.SetEvolutionaryStage(info1.Mass, info1.CenterH1, info1.CenterHe4, info1.LogTcntr)

   // create struct of MESAstar_info for star 2
   info2 := new(io.MESAstar_info)
   have2stars := false
   if (binfo.PointMassIndex == 0) {

      log.Printf("(re)loading struct with MESAstar info for star 2")
      info2.HistoryName = star2filePath
      io.GrabStarHeader(star2filePath, info2)
      io.GrabStarRunInfo(star2filePath, info2)
      info2.EvolState = utils.SetEvolutionaryStage(info2.Mass, info2.CenterH1, info2.CenterHe4, info2.LogTcntr)

      have2stars = true

   }

   // need to set case of MT (if present)
   if (binfo.DonorIndex == 1) {

      binfo.MTCase = utils.SetMTCase(binfo.RelRLOF1, info1.EvolState)

   } else {

      binfo.MTCase = utils.SetMTCase(binfo.RelRLOF2, info2.EvolState)

   }

   // put everything in the same struct
   info := io.CompleteBinaryInfo {

      BinaryInfo: binfo,
      Star1Info: info1,
      Star2Info: info2,
      Have2Stars: have2stars,

   }

   tpl_mesabinary := template.Must(template.ParseFiles("web/index-mesabinary.html"))
   tpl_mesabinary.Execute(w, info)
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
