// Package routes defines  HTTP routes for the Office order.
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
	controllerscommon "Hrmodule/controllers/common"

	"net/http"

	"github.com/gin-gonic/gin"
)

func Common(router *gin.Engine) {

	router.Any("/Defaultrole", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.DefaultRoleName))))
	router.Any("/TaskInbox", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.InboxTasksRole))))
	router.Any("/TaskSummary", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.TaskSummary))))
	router.Any("/Statusmaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.StatusMaster))))
	router.Any("/Inboxactivity", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.TaskUpdateHandler))))
	router.Any("/Statusmasternew", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.StatusMasternew))))
	router.Any("/download-signature", gin.WrapH(http.HandlerFunc(controllerscommon.DownloadSignatureHandler))) //without jwt
	router.Any("/ComboValueMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.ComboValueMaster))))
	router.Any("/Religion", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.Religion))))
	router.Any("/BloodGroupMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.BloodGroupMaster))))
	router.Any("/CasteCategory", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.CasteCategory))))
	router.Any("/LanguageMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.LanguageMaster))))
	router.Any("/OfficialLanguage", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.OfficialLanguage))))
	router.Any("/DesignationMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.DesignationMaster))))
	router.Any("/DepartmentMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.DepartmentMaster))))
	router.Any("/Year", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.Year))))
	router.Any("/EmployeePresentScaleMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.EmployeePresentScaleMaster))))
	router.Any("/CountryMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.CountryMaster))))
	router.Any("/StateMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.StateMaster))))
	router.Any("/DistrictMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.DistrictMaster))))
	router.Any("/BankMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.BankMaster))))
	router.Any("/CityMaster", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllerscommon.CityMaster))))
}
