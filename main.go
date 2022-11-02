package main

import (
	"car-park/controllers"
	"car-park/middlewares"
	"car-park/models"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	vehicleDB      models.VehicleDB      = models.NewVehicleDB()
	vehicleService models.VehicleService = models.NewVehicleService(vehicleDB)

	vehicleController controllers.VehicleController = controllers.NewVehicleController(vehicleService)
)

func setupLogFile() {
	f, err := os.Create("gin.log")
	if err != nil {
		fmt.Println("[WARNING] Cannot create logfile, will use only stdout.")
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

}

func main() {
	// vehicleDB.DestructiveReset()
	setupLogFile()
	defer vehicleDB.CloseConnection()

	server := gin.New()
	server.Use(gin.Recovery(), middlewares.Logger())

	server.Static("/css", "./views/templates/css")
	server.LoadHTMLGlob("./views/templates/*.html")

	viewRoutes := server.Group("/view")
	{
		viewRoutes.GET("/vehicles", vehicleController.ShowAll)
	}

	apiRoutes := server.Group("/api")
	{
		apiRoutes.GET("/vehicles", func(ctx *gin.Context) {
			ctx.JSON(200, vehicleController.FindAll())
		})

		// POST new vehicle
		apiRoutes.POST("/vehicles", func(ctx *gin.Context) {
			err := vehicleController.Save(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Success!"})
			}

		})

		// PUT is update existing vehicle
		apiRoutes.PUT("/vehicles/:id", func(ctx *gin.Context) {
			err := vehicleController.Update(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Success!"})
			}

		})

		apiRoutes.DELETE("/vehicles/:id", func(ctx *gin.Context) {
			err := vehicleController.Delete(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Success!"})
			}

		})

	}

	// {
	// 	apiRoutes.GET("/vehicles", func(ctx *gin.Context) {
	// 		ctx.JSON(200, vehicleController.FindAll())
	// 	})

	// 	apiRoutes.POST("/vehicles", func(ctx *gin.Context) {
	// 		err := vehicleController.Save(ctx)
	// 		if err != nil {
	// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 			return
	// 		}

	// 		ctx.JSON(http.StatusOK, gin.H{"message": "Saved successfully."})

	// 	})
	// }

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	server.Run(":" + port)
}
