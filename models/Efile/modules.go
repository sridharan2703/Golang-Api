// Package modelsefile contains structs and queries for ALLModules.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 24-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all ALLModules Master.

package modelsefile

// ALLModulesMasterQuery - Query to fetch ALLModules values
const ALLModulesMasterQuery = `SELECT module_name
FROM meivan.category_role_map a 
JOIN meivan.category_visibility b 
    ON b.id = a.module_id
WHERE role_name = $1 and a.status='1'
ORDER BY module_name`