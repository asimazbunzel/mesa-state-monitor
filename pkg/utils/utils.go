package utils


// define phase of evolution for a star based on abundances and central temperature
func Set_evolutionary_stage(mass float64, center_h1 float64, center_he4 float64, log_T_cntr float64) string {

   // H threshold
   h1_threshold := 1e-5

   // He$ threshold
   he4_threshold := 1e-4

   // log T for He ignition
   log_he_temp := 7.8

   // chandra mass
   chandra_mass := 1.4

   // Main sequence star
   if (center_h1 > h1_threshold) {return "MS star"}

   // Hertzsprung gap star or WD
   if (log_T_cntr < log_he_temp) {

      // depending on the center He4 value, either HG star of becoming WD or ECSN
      if (center_he4 < he4_threshold) {

         // chandra mass separates between WD & ECSN
         if (mass < chandra_mass) {
            return "WD"
         } else {
            return "He depleted star, possible EC SN"
         }

      } else {

         // HG star
         return "HG star"
      }
   }

   // core He burning star
   if (center_he4 > he4_threshold) {

      return "CHeB star"

   } else {

      // already past core He depletion
      // chandra mass separates between WD & ECSN
      if (mass < chandra_mass) {
         return "WD"
      } else {
         return "He depleted star"
      }
   }

   return "none"
}
