// Package routes defines  HTTP routes for the Office order.

// --- Creator's Info ---
// Creator: Rovita
// Created On: 24-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 24-11-2025


package routes

import (
	"Hrmodule/auth"
	controllersefile "Hrmodule/controllers/Efile"
	"net/http"

	"github.com/gin-gonic/gin"
)

func EmployeeEfile(router *gin.Engine) {

	router.Any("/ALLRoles", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersefile.ALLRoles))))
	router.Any("/ALLModules", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersefile.ALLModules))))
	router.Any("/ALLMenus", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersefile.ALLMenus))))
	router.Any("/InsertModules", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersefile.InsertModules))))
	router.Any("/UpdateModules", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersefile.UpdateModules))))


	
}
