// Package modelsefile contains structs and queries for ALLMenus.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 24-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all ALLMenus Master.

package modelsefile

// ALLMenusMasterQuery - Query to fetch ALLMenus values for specific role
const ALLMenusMasterQuery = `SELECT 
    mcv.module_name,
    CASE 
        WHEN mcr.module_id IS NOT NULL THEN 'YES'
        ELSE 'NO'
    END AS module_status, 
    mcv.id as module_id
FROM meivan.category_visibility mcv
LEFT JOIN meivan.category_role_map mcr 
    ON mcv.id = mcr.module_id  
    AND mcr.status = '1'
    AND mcr.role_name = $1 
WHERE mcv.status = '1'
ORDER BY mcv.id`

// ALLMenusMasterGroupedByRoleQuery - Query to fetch all modules grouped by role_name
const ALLMenusMasterGroupedByRoleQuery = `SELECT 
    DISTINCT mcr.role_name,
    mcv.id as module_id,
    mcv.module_name,
    CASE 
        WHEN mcr.module_id IS NOT NULL AND mcr.status = '1' THEN 'YES'
        ELSE 'NO'
    END AS module_status
FROM meivan.category_visibility mcv
LEFT JOIN meivan.category_role_map mcr 
    ON mcv.id = mcr.module_id
WHERE mcv.status = '1'
ORDER BY mcr.role_name NULLS LAST, mcv.id`