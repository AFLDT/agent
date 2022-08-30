package http

import (
	"image"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/kerberos-io/agent/machinery/src/cloud"
	"github.com/kerberos-io/agent/machinery/src/components"
	"github.com/kerberos-io/agent/machinery/src/log"
	"github.com/kerberos-io/agent/machinery/src/models"
	"github.com/kerberos-io/agent/machinery/src/utils"
)

func AddRoutes(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware, configuration *models.Configuration, communication *models.Communication) *gin.RouterGroup {

	// This is legacy should be removed in future! Now everything
	// lives under the /api prefix.
	r.GET("/config", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"config":   configuration.Config,
			"custom":   configuration.CustomConfig,
			"global":   configuration.GlobalConfig,
			"snapshot": components.GetSnapshot(),
		})
	})

	// This is legacy should be removed in future! Now everything
	// lives under the /api prefix.
	r.POST("/config", func(c *gin.Context) {
		var config models.Config
		err := c.BindJSON(&config)
		if err == nil {
			err := components.SaveConfig(config, configuration, communication)
			if err == nil {
				c.JSON(200, gin.H{
					"data": "☄ Reconfiguring",
				})
			} else {
				c.JSON(400, gin.H{
					"data": "Something went wrong: " + err.Error(),
				})
			}
		} else {
			c.JSON(400, gin.H{
				"data": "Something went wrong: " + err.Error(),
			})
		}
	})

	api := r.Group("/api")
	{
		api.POST("/login", authMiddleware.LoginHandler)

		api.GET("/dashboard", func(c *gin.Context) {

			// This will return the timestamp when the last packet was correctyl received
			// this is to calculate if the camera connection is still working.
			lastPacketReceived := int64(0)
			if communication.LastPacketTimer != nil {
				lastPacketReceived = communication.LastPacketTimer.Load().(int64)
			}

			// The total number of recordings stored in the directory
			numberOfRecordings := utils.NumberOfFilesInDirectory("./data/recordings")

			c.JSON(200, gin.H{
				"cameraOnline":       lastPacketReceived,
				"numberOfRecordings": numberOfRecordings,
			})
		})

		// Streaming handler
		api.GET("/stream", func(c *gin.Context) {
			// TODO add a token validation!
			imageFunction := func() (image.Image, error) {
				// We will only send an image once per second.
				time.Sleep(time.Second * 1)
				log.Log.Info("AddRoutes (/stream): reading from MJPEG stream")
				img, err := components.GetImageFromFilePath()
				return img, err
			}
			h := components.StartMotionJPEG(imageFunction, 80)
			h.ServeHTTP(c.Writer, c.Request)
		})

		api.GET("/config", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"config":   configuration.Config,
				"custom":   configuration.CustomConfig,
				"global":   configuration.GlobalConfig,
				"snapshot": components.GetSnapshot(),
			})
		})

		api.POST("/config", func(c *gin.Context) {
			var config models.Config
			err := c.BindJSON(&config)
			if err == nil {
				err := components.SaveConfig(config, configuration, communication)
				if err == nil {
					c.JSON(200, gin.H{
						"data": "☄ Reconfiguring",
					})
				} else {
					c.JSON(200, gin.H{
						"data": "☄ Reconfiguring",
					})
				}
			} else {
				c.JSON(400, gin.H{
					"data": "Something went wrong: " + err.Error(),
				})
			}
		})

		api.GET("/restart", func(c *gin.Context) {
			communication.HandleBootstrap <- "restart"
			c.JSON(200, gin.H{
				"restarted": true,
			})
		})

		api.GET("/stop", func(c *gin.Context) {
			communication.HandleBootstrap <- "stop"
			c.JSON(200, gin.H{
				"stopped": true,
			})
		})

		api.POST("/hub/verify", func(c *gin.Context) {
			cloud.VerifyHub(c)
		})

		api.POST("/persistence/verify", func(c *gin.Context) {
			cloud.VerifyPersistence(c)
		})

		api.Use(authMiddleware.MiddlewareFunc())
		{
			// Secured endpoints..

		}
	}
	return api
}
