// Package routes defines  HTTP routes for the Login.
//
// --- Creator's Info ---
//
// Creator: Sridharan
// Created On: 31-10-2025

// Last Modified By: Sridharan

// Last Modified Date: 31-10-2025
package routes

import (
	"Hrmodule/auth"
	controllerslogin "Hrmodule/controllers/login"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(router *gin.Engine) {
	router.Any("/HRldap", gin.WrapH(http.HandlerFunc(controllerslogin.HandleLDAPAuth)))
	router.Any("/HRldapfailure", gin.WrapH(http.HandlerFunc(controllerslogin.HandleLDAPAuthf)))
	router.Any("/Loginotp", gin.WrapH(http.HandlerFunc(controllerslogin.InsertOTPHandler)))
	router.Any("/Loginotpupdate", gin.WrapH(http.HandlerFunc(controllerslogin.ValidateOTPHandler)))
	router.Any("/SessionTimeout", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerslogin.SessionTimeoutHandler))))
	router.Any("/Sessiondata", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerslogin.SessionData))))
	router.Any("/InsertUserActivityLog", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerslogin.InsertUserActivityLog))))
	router.Any("/SendOTP", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerslogin.SendOTPHandler))))
	router.Any("/Datadecrypt", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerslogin.DatadecryptHandler))))
	router.Any("/DatadecryptKey", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerslogin.Datadecryptsessionkey))))
}
