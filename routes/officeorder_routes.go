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
	controllersofficeorder "Hrmodule/controllers/officeorder"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Officeorder(router *gin.Engine) {
	router.Any("/OfficeOrder_module", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.OrderSubModule))))
	router.Any("/OfficeOrder_visitdetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.Ordervisitdetails))))
	router.Any("/OfficeOrder_InsertOfficedetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.PCRInsert))))
	router.Any("/OfficeOrder_Count", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.NeedGenerateHandler))))
	router.Any("/OfficeOrder_DropdownValuesHandler", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.DropdownValuesHandler))))
	router.Any("/OfficeOrder_statusupdate", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.OfficeOrderUpdateTaskStatus))))
	router.Any("/OfficeOrder_approval_remarks", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.OfficeComments))))
	router.Any("/OfficeOrder_datatemplate", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.GetOfficeOrderDetailsfortemplate))))
	router.Any("/OfficeOrder_ReturnDropdown", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.ReturnDropdown))))
	router.Any("/OfficeOrder_taskvisitdetails", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.OrderTaskVisitDetails))))
	router.Any("/OfficeOrder_History", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.OfficeOrderHistory))))
	router.Any("/OfficeOrder_Historypdf", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.FetchOrderHistoryPDF))))
	router.Any("/OfficeOrder_CcRoles", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.CcRoles))))
	router.Any("/OfficeOrder_Tasksummary", gin.WrapH(auth.JwtMiddleware(http.HandlerFunc(controllersofficeorder.TaskDetailsTaskSummary))))
}
