// Package routes defines  HTTP routes for the Office order.

// --- Creator's Info ---
// Creator: Vaishnavi
// Created On: 04-11-2025
// Last Modified By:  Rovita
// Last Modified Date: 12-11-2025


package routes

import (
	"Hrmodule/auth"
	controllerssad "Hrmodule/controllers/Staffadditionaldetails"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Staffadditionaldetails(router *gin.Engine) {

	router.Any("/EmployeeDropdown", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerssad.EmployeeDropdown))))
	router.Any("/EmployeeEfileDetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerssad.EmployeeEfileDetails))))
	router.Any("/RoleBasedModulesHandler", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerssad.RoleBasedModulesHandler))))
	router.Any("/SadBasicDetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerssad.SadBasicDetails))))
	router.Any("/EmployeeBasicDetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerssad.EmployeePersonalDetailsHandler))))
	router.Any("/EmployeeAppointmentDetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerssad.EmployeeAppointmentDetailsHandler))))
	router.Any("/EmployeeEducationDetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerssad.EmployeeEducationDetailsHandler))))
	router.Any("/EmployeeExperienceDetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerssad.EmployeeExperienceDetailsHandler))))
	router.Any("/EmployeeDependentDetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerssad.EmployeeDependentDetailsHandler))))
	router.Any("/EmployeeLanguageDetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerssad.EmployeeLanguageDetailsHandler))))
	router.Any("/EmployeeDocumentDetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerssad.EmployeeDocumentDetailsHandler))))
	router.Any("/EmployeeNomineeDetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerssad.EmployeeNomineeDetailsHandler))))
	router.Any("/EmployeeContactDetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerssad.EmployeeContactDetailsHandler))))

	
}
