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
	setupLogFile()
	defer vehicleDB.CloseConnection()

	server := gin.New()
	server.Use(gin.Recovery(), middlewares.Logger())

	server.Static("/css", "./views/templates/css")
	server.StaticFile("/logo.svg", "./views/assets/logo-1-white.svg")
	server.LoadHTMLGlob("./views/templates/*.html")

	server.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": "404", "message": "PAGE NOT FOUND"})
	})

	viewRoutes := server.Group("/view")
	{
		viewRoutes.GET("/vehicles", vehicleController.ShowAllVehicles)
		viewRoutes.GET("/vehicles/create", vehicleController.ShowCreateVehicle)
		viewRoutes.GET("/vehicles/edit/:id", vehicleController.ShowEditVehicle)
		viewRoutes.GET("/carmodels", vehicleController.ShowAllCarModels)
	}

	apiRoutes := server.Group("/api")
	{
		// apiRoutes.GET("/vehicles", func(ctx *gin.Context) {
		// 	ctx.JSON(200, vehicleController.FindAllVehicles())
		// })

		// POST new vehicle
		apiRoutes.POST("/vehicles", func(ctx *gin.Context) {
			err := vehicleController.SaveVehicle(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				// ctx.JSON(http.StatusOK, gin.H{"message": "Success!"})
				ctx.Redirect(http.StatusFound, "/view/vehicles")
			}

		})

		// PUT is update existing vehicle
		apiRoutes.POST("/update/vehicle/:id", func(ctx *gin.Context) {
			vehicleController.UpdateVehicle(ctx)
			err := vehicleController.UpdateVehicle(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.Redirect(http.StatusFound, "/view/vehicles")
			}

		})

		apiRoutes.POST("/delete/vehicle/:id", func(ctx *gin.Context) {
			err := vehicleController.DeleteVehicle(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.Redirect(http.StatusFound, "/view/vehicles")
			}

		})

	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	// vehicleDB.DestructiveReset()
	server.Run(":" + port)
}
