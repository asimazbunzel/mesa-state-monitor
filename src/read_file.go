package main

import (
   "bufio"
   "fmt"
   "log"
   "strings"
   "strconv"
   "os"
)

// struct holding info on MESAstar
type MESAstar_info struct {
   version int
   date string
   history_name string
}

// struct holding info on MESAbinary
type MESAbinary_info struct {
   period float64
   star_1_mass, star_2_mass float64
   MT_case string
   history_name string
}

func grab_star_header (fname string, star_info *MESAstar_info) {

   // MESA specfic row numbers for header names & values in history output
   nr_header_names := 2
   nr_header_values := 3

   // open star file
   fstar, err := os.Open(fname)
   if err != nil {log.Fatal(err)}
   defer fstar.Close()

   // scan star file
   scanner := bufio.NewScanner(fstar)

   // arrays holding star header names & values
   var header_names []string
   var header_values []string
   var header_names_found, header_values_found bool
   lineCount := 0

   for scanner.Scan() {

      lineCount++

      // get header names
      if lineCount == nr_header_names {
         header_names = strings.Fields(scanner.Text())
         header_names_found = true
      }

      // find header values
      if lineCount == nr_header_values {
         header_values = strings.Fields(scanner.Text())
         header_values_found = true
      }

      if (header_names_found && header_values_found) {
         for k, name := range header_names {
            if name == "version_number" {
               i, err := strconv.Atoi(strings.Split(header_values[k], "\"")[1])
               // handle error
               if err != nil {
                  fmt.Println(err)
                  os.Exit(2)
               }
               star_info.version = i
            }
            if name == "date" {star_info.date = header_values[k]}
         }
      }
   }

   if err := scanner.Err(); err != nil {
      log.Fatal(err)
   }
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
   Info := new(MESAstar_info)
   Info.history_name = starfilePath
   grab_star_header(starfilePath, Info)
   fmt.Println(Info)

   // now get info on the star using the row containing names for data columns
   // and the last row written by the MESA code

   if (is_binary_evolution) {
      BInfo := new(MESAbinary_info)
      BInfo.history_name = ""
   }
}
