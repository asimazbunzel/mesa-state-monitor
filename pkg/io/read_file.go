package io

import (
   "bufio"
   "fmt"
   "io"
   "log"
   "strings"
   "strconv"
   "os"
)


// function retrieve from this post:
// https://stackoverflow.com/questions/17863821/how-to-read-last-lines-from-a-big-file-with-go-every-10-secs
func getLastLineWithSeek (filepath string) string {
    fileHandle, err := os.Open(filepath)

    if err != nil {log.Fatal("cannot open file")}
    defer fileHandle.Close()

    line := ""
    var cursor int64 = 0
    stat, _ := fileHandle.Stat()
    filesize := stat.Size()
    for {
        cursor -= 1
        fileHandle.Seek(cursor, io.SeekEnd)

        char := make([]byte, 1)
        fileHandle.Read(char)

        if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
            break
        }

        line = fmt.Sprintf("%s%s", string(char), line) // there is more efficient way

        if cursor == -filesize { // stop if we are at the begining
            break
        }
    }

    return line
}


func GrabStarHeader (fname string, star_info *MESAstar_info) {

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
               star_info.Version = i
            }
            if name == "date" {star_info.Date = header_values[k]}
         }
      }
   }

   if err := scanner.Err(); err != nil {
      log.Fatal(err)
   }
}


func GrabStarRunInfo (fname string, star_info *MESAstar_info) {

   nr_column_names := 6

   // open star file
   fstar, err := os.Open(fname)
   if err != nil {log.Fatal(err)}
   defer fstar.Close()

   // scan star file
   scanner := bufio.NewScanner(fstar)

   // arrays holding star header names & values
   var column_names []string
   var column_values []string
   var column_names_found bool
   lineCount := 0

   for scanner.Scan() {

      lineCount++

      // get header names
      if lineCount == nr_column_names {
         column_names = strings.Fields(scanner.Text())
         column_names_found = true
      }

      if column_names_found {break}
   }

   if (column_names_found) {
      column_values = strings.Fields(getLastLineWithSeek(fname))
   }

   if (column_names_found) {
      for k, name := range column_names {

         val := column_values[k]

         if name == "model_number" {
            i, err := strconv.Atoi(val)
            // handle error
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            star_info.ModelNumber = i
         }
         if name == "num_zones" {
            i, err := strconv.Atoi(val)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            star_info.NumZones = i
         }
         if name == "star_mass" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            star_info.Mass = i
         }
         if name == "log_abs_mdot" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            star_info.LogMdot = i
         }
         if name == "star_age" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            star_info.Age = i
         }
         if name == "center_h1" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            star_info.CenterH1 = i
         }
         if name == "center_he4" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            star_info.CenterHe4 = i
         }
         if name == "log_center_T" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            star_info.LogTcntr = i
         }
         if name == "num_retries" {
            i, err := strconv.Atoi(val)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            star_info.NumRetries = i
         }
         if name == "num_iters" {
            i, err := strconv.Atoi(val)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            star_info.NumIters = i
         }
         if name == "elapsed_time" {
            i, err := strconv.ParseFloat(val, 64)  // i is in sec
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            star_info.ElapsedTime = i / 60 // from sec to min
         }
      }
   }

   if err := scanner.Err(); err != nil {
      log.Fatal(err)
   }
}


func GrabBinaryHeader (fname string, binary_info *MESAbinary_info) {

   // MESA specfic row numbers for header names & values in history output
   nr_header_names := 2
   nr_header_values := 3

   // open binary file
   fbinary, err := os.Open(fname)
   if err != nil {log.Fatal(err)}
   defer fbinary.Close()

   // scan star file
   scanner := bufio.NewScanner(fbinary)

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
            val := header_values[k]
            if name == "initial_don_mass" {
               i, err := strconv.ParseFloat(val, 64)  // i is in Msun
               if err != nil {
                  fmt.Println(err)
                  os.Exit(2)
               }
               binary_info.InitialDonorMass = i
            }
            if name == "initial_acc_mass" {
               i, err := strconv.ParseFloat(val, 64)  // i is in Msun
               if err != nil {
                  fmt.Println(err)
                  os.Exit(2)
               }
               binary_info.InitialAccretorMass = i
            }
            if name == "initial_period_days" {
               i, err := strconv.ParseFloat(val, 64)  // i is in days
               if err != nil {
                  fmt.Println(err)
                  os.Exit(2)
               }
               binary_info.InitialPeriod = i
            }
         }
      }
   }

   if err := scanner.Err(); err != nil {
      log.Fatal(err)
   }
}


func GrabBinaryRunInfo (fname string, binary_info *MESAbinary_info) {

   nr_column_names := 6

   // open star file
   fbinary, err := os.Open(fname)
   if err != nil {log.Fatal(err)}
   defer fbinary.Close()

   // scan star file
   scanner := bufio.NewScanner(fbinary)

   // arrays holding star header names & values
   var column_names []string
   var column_values []string
   var column_names_found bool
   lineCount := 0

   for scanner.Scan() {

      lineCount++

      // get header names
      if lineCount == nr_column_names {
         column_names = strings.Fields(scanner.Text())
         column_names_found = true
      }

      if column_names_found {break}
   }

   if (column_names_found) {
      column_values = strings.Fields(getLastLineWithSeek(fname))
   }

   if (column_names_found) {
      for k, name := range column_names {

         val := column_values[k]

         if name == "model_number" {
            i, err := strconv.Atoi(val)
            // handle error
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            binary_info.ModelNumber = i
         }
         if name == "age" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            binary_info.Age = i
         }
         if name == "period_days" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            binary_info.Period = i
         }
         if name == "star_1_mass" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            binary_info.Star1Mass = i
         }
         if name == "star_2_mass" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            binary_info.Star2Mass = i
         }
         if name == "donor_index" {
            i, err := strconv.Atoi(val)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            binary_info.DonorIndex = i
         }
         if name == "point_mass_index" {
            i, err := strconv.Atoi(val)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            binary_info.PointMassIndex = i
         }
         if name == "rl_relative_overflow_1" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            binary_info.RelRLOF1 = i
         }
         if name == "rl_relative_overflow_2" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               os.Exit(2)
            }
            binary_info.RelRLOF2 = i
         }
      }
   }
}
