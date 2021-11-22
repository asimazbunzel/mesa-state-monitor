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
// to check if this is a binary simulation, look for the MESAbinary output
func IsBinary (path string) bool {

   log.Printf("searching for binary evolution")

   // name of the MESAbinary output
   binaryFile := fmt.Sprintf("%s/%s", path, binaryHistoryName)

   log.Printf("looking for file `%s`", binaryFile)
   _, err := os.Stat(binaryFile)
   if err != nil {

      log.Printf("`%s` file not found. now searching inside %s folder", binaryFile, binaryLogDirectory)

      binaryFile := fmt.Sprintf("%s/%s/%s", path, binaryLogDirectory, binaryHistoryName)

      _, err := os.Stat(binaryFile)
      if err != nil {

         log.Printf("binary logs not found. single evolution assumed")
         return false

      } else {

         log.Printf("found binary log: `%s`. binary evolution assumed", binaryFile)
         return true

      }
   } else {

      log.Printf("found binary logs. binary evolution assumed")
      return true

   }
}


// return logs names from MESA folder
func GetLogNames (path string, isBinary bool) (string, string, string) {

   log.Printf("searching for MESA LOGS filename(s)")

   binaryLogName := ""
   star1LogName := ""
   star2LogName := ""

   if (isBinary) {

      // search for binary output
      // use defaults values defined at beginning of module
      binaryLogName = fmt.Sprintf("%s/%s", path, binaryHistoryName)
      _, err := os.Stat(binaryLogName)
      if err != nil {
         binaryLogName = fmt.Sprintf("%s/%s/%s", path, binaryLogDirectory, binaryHistoryName)
         _, err = os.Stat(binaryLogName)
         if err != nil {log.Fatal("cannot find binary LOG output file")}
      }
      log.Printf("found binary output: `%s`", binaryLogName)

      // now look for star 1 data
      star1LogName = fmt.Sprintf("%s/%s/%s", path, starLogDirectory, starHistoryName)
      _, err = os.Stat(star1LogName)
      if err != nil {
         star1LogName = fmt.Sprintf("%s/%s/%s", path, star1LogDirectory, starHistoryName)
         _, err = os.Stat(star1LogName)
         if err != nil {
            star1LogName = fmt.Sprintf("%s/%s/%s", path, star1LogDirectory, "primary_history.data")
            _, err = os.Stat(star1LogName)
            if err != nil {log.Fatal("cannot find star 1 LOG output file")}
         }
      }
      log.Printf("found star 1 output: `%s`", star1LogName)

      // now look for star 2 data (though not always found if doing star + point-mass)
      star2LogName = fmt.Sprintf("%s/%s/%s", path, star2LogDirectory, starHistoryName)
      _, err = os.Stat(star2LogName)
      if err != nil {
         star2LogName = fmt.Sprintf("%s/%s/%s", path, star2LogDirectory, "secondary_history.data")
         _, err = os.Stat(star2LogName)
         if err != nil {
            log.Printf("cannot find star 2 LOG output file. maybe doing star + point-mass evolution")
            star2LogName = ""
         } else {
            log.Printf("found star 2 output: `%s`", star2LogName)
         }
      } else {
         log.Printf("found star 2 output: `%s`", star2LogName)
      }

   } else {

      // only need to search for star1LogName
      star1LogName = fmt.Sprintf("%s/%s/%s", path, starLogDirectory, starHistoryName)

      _, err := os.Stat(star1LogName)
      if err != nil {log.Fatal("cannot find star LOG output file of single evolution")}
      log.Printf("found single evolution output: `%s`", star1LogName)
   }

   return binaryLogName, star1LogName, star2LogName
}
