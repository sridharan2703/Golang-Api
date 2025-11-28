package mainroutes

import (
	"Hrmodule/routes"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Registerroutes() {
	// Create Gin router
	router := gin.Default()

	// --- CORS Configuration ---
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Change this in production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	/*************************************************Cron Jobs*****************************************************/
	routes.Cronjobs(router)
	/***************************************************************************************************************/
	/*************************************************Common API's*****************************************************/
	routes.Common(router)
	/***************************************************************************************************************/
	/*************************************************LOGIN  API's*****************************************************/
	routes.Login(router)
	/***************************************************************************************************************/
	/*************************************************OFFICE ORDER API's*****************************************************/
	routes.Officeorder(router)
	/***************************************************************************************************************/
	/*************************************************Staffadditionaldetails API's*****************************************************/
	routes.Staffadditionaldetails(router)
	/***************************************************************************************************************/
	/*************************************************Employee Efile API's*****************************************************/
	routes.EmployeeEfile(router)
	/***************************************************************************************************************/

	// --- HTTPS SERVER START ---
	fmt.Println("🚀 Server starting on port 2703 (HTTPS Enabled)")

	//certFile := "certificate.pem"
	//keyFile := "key.pem"

	if err := router.Run(":2703"); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}

/* with client ip restiction
package mainroutes

import (
	"Hrmodule/routes"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// getServerIP returns the server's local IP address
func getServerIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "unknown"
}

// getClientIP extracts the client's IP address from the request
func getClientIP(c *gin.Context) string {
	// Try common proxy headers first
	headers := []string{"X-Forwarded-For", "X-Real-IP"}
	for _, h := range headers {
		if ip := c.GetHeader(h); ip != "" {
			return strings.Split(ip, ",")[0]
		}
	}

	// Fallback to the direct remote address
	ip := c.ClientIP()
	return ip
}

func Registerroutes() {
	router := gin.Default()

	// --- CORS Configuration ---
	config := cors.Config{
		AllowOrigins:     []string{"*"}, // ✅ Only allow your React app
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))

	// --- Middleware to print IP info and validate origin ---
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		clientIP := getClientIP(c)
		method := c.Request.Method
		path := c.Request.URL.Path

		// --- CORS Enforcement ---
		if origin == "" {
			c.AbortWithStatusJSON(403, gin.H{"error": "Access denied: Missing Origin header"})
			fmt.Printf("❌ Blocked request (no Origin): IP=%s | Path=%s\n", clientIP, path)
			return
		}

		if origin != "https://wftest2.iitm.ac.in:3737" {
			c.AbortWithStatusJSON(403, gin.H{"error": "Access denied: Unauthorized Origin"})
			fmt.Printf("❌ Blocked unauthorized Origin: %s | IP=%s\n", origin, clientIP)
			return
		}

		// ✅ Log valid request details
		fmt.Printf("\n============================\n")
		fmt.Printf("🕒 Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Printf("🌐 Client IP: %s\n", clientIP)
		fmt.Printf("🔗 Origin: %s\n", origin)
		fmt.Printf("➡️  Method: %s | Path: %s\n", method, path)
		fmt.Printf("============================\n\n")

		c.Next()
	})

	// --- ROUTE REGISTRATION ---
	routes.Cronjobs(router)
	routes.Common(router)
	routes.Login(router)
	routes.Officeorder(router)

	// --- HTTPS SERVER START ---
	fmt.Println("🚀 Server starting on port 7007 (HTTPS Enabled)")

	certFile := "certificate.pem"
	keyFile := "key.pem"

	if err := router.RunTLS(":7007", certFile, keyFile); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
*/
