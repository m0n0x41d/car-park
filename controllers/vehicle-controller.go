package controllers

import (
	"car-park/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type VehicleController interface {
	FindAll() []models.Vehicle
	Save(ctx *gin.Context) error
	Update(ctx *gin.Context) error
	Delete(ctx *gin.Context) error
	ShowAll(ctx *gin.Context)
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

func (c *vehicleController) FindAll() []models.Vehicle {
	return c.service.FindAll()
}

func (c *vehicleController) Save(ctx *gin.Context) error {
	var vehicle models.Vehicle
	err := ctx.ShouldBindJSON(&vehicle)
	if err != nil {
		return err
	}

	err = validate.Struct(vehicle)
	if err != nil {
		return err
	}

	return c.service.Save(vehicle)
}

func (c *vehicleController) Update(ctx *gin.Context) error {
	var vehicle models.Vehicle
	err := ctx.ShouldBindJSON(&vehicle)
	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		return err
	}

	vehicle.ID = id
	err = validate.Struct(vehicle)
	if err != nil {
		return err
	}

	c.service.Update(vehicle)
	return nil
}

func (c *vehicleController) Delete(ctx *gin.Context) error {
	var vehicle models.Vehicle
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		return err
	}
	vehicle.ID = id
	c.service.Delete(vehicle)
	return nil
}

func (c *vehicleController) ShowAll(ctx *gin.Context) {
	vehicles := c.service.FindAll()
	data := gin.H{
		"title":    "Vihecles Page",
		"vehicles": vehicles,
	}

	ctx.HTML(http.StatusOK, "index.html", data)
}
