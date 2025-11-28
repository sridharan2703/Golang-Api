// Package modelsefile contains structs and queries for UpdateModules.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 24-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to update status in category_role_map.

package modelsefile

// UpdateCategoryRoleMapQuery - Query to update status in category_role_map for single module_id
const UpdateCategoryRoleMapQuery = `
UPDATE meivan.category_role_map 
SET status = $1, updated_at = NOW() 
WHERE module_id = $2 AND role_name = $3`

// UpdateCategoryRoleMapMultipleQuery - Query template for multiple module_ids (used with fmt.Sprintf)
const UpdateCategoryRoleMapMultipleQuery = `
UPDATE meivan.category_role_map 
SET status = $1, updated_at = NOW() 
WHERE module_id IN (%s) AND role_name = $%d`