// Package modelsefile contains structs and queries for InsertModules.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 24-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to insert data into category_role_map.

package modelsefile

// InsertCategoryRoleMapQuery - Query to insert data into category_role_map
const InsertCategoryRoleMapQuery = `
INSERT INTO meivan.category_role_map (module_id, role_name, created_at, updated_at, status)
VALUES ($1, $2, NOW(), NOW(), $3)`