package io


var star_history_name = "history.data"
var binary_history_name = "binary_history.data"

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


func IsBinary (path string) bool {

   return false

}
