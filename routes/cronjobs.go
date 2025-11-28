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
	"Hrmodule/controllers/cronjob"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func InitCronforofficeorder() *cron.Cron {
	c := cron.New()
	_, err := c.AddFunc("@every 1h", cronjob.SyncData)
	if err != nil {
		log.Println("❌ Failed to start hourly cron:", err)
		return nil
	}
	c.Start()
	log.Println("✅ Hourly cron started for SyncData()")
	go cronjob.SyncData()
	return c
}

func Cronjobs(router *gin.Engine) {

	// --- OFFICE ORDER CRON JOBS ---
	InitCronforofficeorder()

	router.Any("/OfficeOrder_Sync", gin.WrapH(http.HandlerFunc(cronjob.HandleSync)))
}
