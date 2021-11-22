package io


import (
   "fmt"
   "log"
   "os"
)


var starHistoryName = "history.data"
var binaryHistoryName = "binary_history.data"

var binaryLogDirectory = "LOGS_binary"
var starLogDirectory = "LOGS"
var star1LogDirectory = "LOGS1"
var star2LogDirectory = "LOGS2"


// struct holding info on MESAstar
type MESAstar_info struct {
   Version int
   Date string
   History_name string
   Model_number int
   Num_zones int
   Mass float64
   Log_mdot float64
   Age float64
   Center_h1, Center_he4 float64
   Log_T_cntr float64
   Num_retries, Num_iters int
   Elapsed_time float64
   Evol_state string
}


// struct holding info on MESAbinary
type MESAbinary_info struct {
   model_number int
   initial_don_mass, initial_acc_mass float64
   initial_period float64
   age float64
   star_1_mass, star_2_mass float64
   period float64
   MT_case string
   History_name string
   donor_index, point_mass_index int
   rel_rl_1, rel_rl_2 float64
}


// bool function to find out if path contains a single or binary evolution
func IsBinary (path string) bool {

   log.Print("searching for binary evolution")

   binaryFile := fmt.Sprintf("%s/%s", path, binaryHistoryName)

   log.Print("looking for file ", binaryFile)

   _, err := os.Open(binaryFile)
   if err != nil {

      log.Print(binaryFile, " not found. now searching inside ", binaryLogDirectory, " folder")

      binaryFile := fmt.Sprintf("%s/%s/%s", path, binaryLogDirectory, binaryHistoryName)

      _, err := os.Open(binaryFile)
      if err != nil {

         log.Print("binary logs not found. single evolution assumed")
         return false

      } else {

         log.Print("found binary logs. binary evolution assumed")
         return true

      }
   } else {

      log.Print("found binary logs. binary evolution assumed")
      return true

   }
}
