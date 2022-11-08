package controllers

import (
	"car-park/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
)

type VehicleController interface {
	FindAllVehicles() []models.Vehicle
	FindAllCarModels() []models.CarModel
	SaveVehicle(ctx *gin.Context) error
	UpdateVehicle(ctx *gin.Context) error
	DeleteVehicle(ctx *gin.Context) error
	ShowAllVehicles(ctx *gin.Context)
	ShowAllCarModels(ctx *gin.Context)
	ShowCreateVehicle(ctx *gin.Context)
	ShowEditVehicle(ctx *gin.Context)
}

type vehicleController struct {
	service models.VehicleService
}

var validate *validator.Validate

func NewVehicleController(svc models.VehicleService) VehicleController {
	validate = validator.New()
	return &vehicleController{
		service: svc,
	}
}

func (c *vehicleController) FindAllVehicles() []models.Vehicle {
	return c.service.FindAllVehicles()
}

func (c *vehicleController) FindAllCarModels() []models.CarModel {
	return c.service.FindAllCarModels()
}

func (c *vehicleController) SaveVehicle(ctx *gin.Context) error {
	var vehicle models.Vehicle

	fmt.Println("HELLOFROM SAVE CONTROLLER")
	fmt.Println(vehicle)
	err := ctx.Bind(&vehicle)

	fmt.Println(vehicle)
	// err := ctx.ShouldBind(&vehicle)
	// ctx.Request.ParseForm()
	// fmt.Println(ctx.Request.PostForm)

	err = validate.Struct(vehicle)
	if err != nil {
		return err
	}

	err = c.service.SaveVehicle(vehicle)
	return err
}

func (c *vehicleController) UpdateVehicle(ctx *gin.Context) error {
	var vehicleOrig models.Vehicle

	id, _ := strconv.ParseUint(ctx.Param("id"), 0, 0)

	vehicleOrig = c.service.VehicleByID(uint(id))
	validate.Struct(vehicleOrig)
	fmt.Println("HELLO FROM UPDATE CONTROLLER: THIS IS ORIGINAL")
	fmt.Println(vehicleOrig)
	// data := gin.H{
	// 	"title":   "Edit car in stock",
	// 	"vehicle": vehicle,
	// }

	var newVehicle models.Vehicle
	err := ctx.Bind(&newVehicle)
	mergo.Merge(&newVehicle, vehicleOrig)
	fmt.Println("THIS IS NEW ONE")
	fmt.Println(newVehicle)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	err = validate.Struct(newVehicle)
	if err != nil {
		return err
	}

	return c.service.UpdateVehicle(newVehicle)
}

func (c *vehicleController) DeleteVehicle(ctx *gin.Context) error {
	var vehicle models.Vehicle
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		return err
	}
	vehicle.ID = uint(id)
	c.service.DeleteVehicle(vehicle)
	return nil
}

func (c *vehicleController) ShowAllVehicles(ctx *gin.Context) {
	vehicles := c.service.FindAllVehicles()
	data := gin.H{
		"title":    "Vihecles Stock",
		"vehicles": vehicles,
	}

	ctx.HTML(http.StatusOK, "vehicles.html", data)
}

func (c *vehicleController) ShowAllCarModels(ctx *gin.Context) {
	carModels := c.service.FindAllCarModels()
	data := gin.H{
		"title":     "Cars Specifications",
		"carmodels": carModels,
	}

	ctx.HTML(http.StatusOK, "carmodels.html", data)
}

func (c *vehicleController) ShowCreateVehicle(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "create-vehicle.html", gin.H{})
}

func (c *vehicleController) ShowEditVehicle(ctx *gin.Context) {
	var vehicle models.Vehicle
	ctx.ShouldBind(&vehicle)

	id, _ := strconv.ParseUint(ctx.Param("id"), 0, 0)

	vehicle = c.service.VehicleByID(uint(id))
	validate.Struct(vehicle)
	fmt.Println(vehicle)
	data := gin.H{
		"title":   "Edit car in stock",
		"vehicle": vehicle,
	}

	ctx.HTML(http.StatusOK, "edit-vehicle.html", data)
}
