// Package modelsefile contains structs and queries for ALLRoles.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 24-11-2025
// Last Modified By:  
// Last Modified Date: 
// This api is to fetch the all ALLRoles Master.

package modelsefile

// ALLRolesMasterQuery - Query to fetch ALLRoles values
const ALLRolesMasterQuery = `SELECT DISTINCT
    B.campuscode || ' ' ||
    CASE
      WHEN A.sectionid IS NOT NULL THEN D.sectioncode
        ELSE A.departmentcode
    END || ' ' || C.rolename AS rolename
   
FROM
    humanresources.employeerolemapping A
JOIN
    humanresources.campus B ON A.campusid = B.id
JOIN
    meivan.rolemaster C ON A.roleid = C.id
LEFT JOIN 
    humanresources.section D ON A.sectionid = D.id
where C.status='1'`