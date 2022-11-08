package models

import "gorm.io/gorm"

// "github.com/jinzhu/gorm"

type Vehicle struct {
	// ID               uint64    `gorm:"primary_index;auto_increment" json:"id"`
	gorm.Model
	Description      string   `json:"description" form:"description"`
	Price            uint     `gorm:"not null" json:"price" form:"price"`
	Mileage          uint     `json:"mileage" form:"mileage"`
	ManufacturedYear uint     `json:"manufactured" form:"manufactured"`
	CarModel         CarModel `json:"carmodel" binding:"required"`
	CarModelID       uint     `json:"-" form:"-"`
}

type CarModel struct {
	// ID               uint64 `gorm:"primary_key;auto_increment" json:"id"`
	// ID               int
	gorm.Model
	Brand            string  `json:"brand" form:"brand"`
	CarType          string  `json:"carmodel" form:"carmodel"`
	FuelTankCapacity uint    `json:"fueltankcapacity" form:"fueltankcapacity"`
	Seats            uint    `json:"seats" form:"seats"`
	EngineVolume     float64 `json:"enginevolume" form:"enginevolume"`
}

type VehicleService interface {
	SaveVehicle(Vehicle) error
	UpdateVehicle(Vehicle) error
	DeleteVehicle(Vehicle) error
	FindAllVehicles() []Vehicle
	FindAllCarModels() []CarModel
	VehicleByID(id uint) Vehicle
}

type vehicleService struct {
	vehicleDB VehicleDB
}

func NewVehicleService(vehicleDB VehicleDB) VehicleService {
	return &vehicleService{
		vehicleDB: vehicleDB,
	}
}

func (service *vehicleService) SaveVehicle(v Vehicle) error {
	return service.vehicleDB.Save(v)
}

func (service *vehicleService) UpdateVehicle(v Vehicle) error {
	return service.vehicleDB.Update(v)
}

func (service *vehicleService) DeleteVehicle(v Vehicle) error {
	return service.vehicleDB.Delete(v)
}

func (service *vehicleService) FindAllVehicles() []Vehicle {
	return service.vehicleDB.FindAllVehicles()
}

func (service *vehicleService) FindAllCarModels() []CarModel {
	return service.vehicleDB.FindAllCarModels()
}

func (service *vehicleService) VehicleByID(id uint) Vehicle {
	return service.vehicleDB.VehicleByID(id)
}
