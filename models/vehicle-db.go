package models

//TODO: Candidate on refactoring...

import (
	// "github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type VehicleDB interface {
	Save(vehicle Vehicle) error
	Update(vehicle Vehicle) error
	Delete(vehicle Vehicle) error
	FindAllVehicles() []Vehicle
	VehicleByID(id uint) Vehicle
	FindAllCarModels() []CarModel
	CloseConnection()
	DestructiveReset() error
}

type dbConn struct {
	connection *gorm.DB
}

//TODO: Make it parameterizable and secure.
var connectionString = "host=localhost port=5432 user=admin password=qwerty dbname=car_park_dev sslmode=disable"

func NewVehicleDB() VehicleDB {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to connect to database")
	}
	// db.LogMode(true)

	db.AutoMigrate(&CarModel{}, &Vehicle{})
	return &dbConn{
		connection: db,
	}
}

func (db *dbConn) CloseConnection() {
	sqlDb, err := db.connection.DB()
	if err != nil {
		panic("Failed to close database connection (fetching interface)")
	}

	err = fmt.Errorf(sqlDb.Close().Error())

	if err != nil {
		panic("Failed to close database connection")
	}
}

func (db *dbConn) Save(v Vehicle) error {
	return db.connection.Create(&v).Error
}

func (db *dbConn) Update(v Vehicle) error {
	return db.connection.Save(&v).Error
}

func (db *dbConn) Delete(v Vehicle) error {
	var findV Vehicle
	findV.ID = v.ID
	db.connection.Preload("CarModel").Find(&findV)
	fmt.Println(findV)
	var carM CarModel
	carM.ID = findV.CarModel.ID
	fmt.Println(carM)
	err := db.connection.Delete(&carM).Error
	if err != nil {
		return err
	}
	return db.connection.Delete(&findV).Error
}

func (db *dbConn) VehicleByID(id uint) Vehicle {
	var vehicle Vehicle
	vehicle.ID = id
	db.connection.Preload("CarModel").Find(&vehicle)
	return vehicle

}

func (db *dbConn) FindAllVehicles() []Vehicle {
	var vehicles []Vehicle
	// in case of appearence foreign keys (nested structs in Vehicle table)
	// change below to db.connection.Set("gorm:auto_preload", true).Find(&vehicles)
	db.connection.Preload("CarModel").Find(&vehicles)
	// db.connection.Find(&vehicles)
	return vehicles
}

func (db *dbConn) FindAllCarModels() []CarModel {
	var carModels []CarModel
	db.connection.Find(&carModels)
	return carModels
}

// DROPDATABASE!
func (db *dbConn) DestructiveReset() error {
	err := db.connection.Delete(&Vehicle{}, &CarModel{}).Error
	if err != nil {
		return err
	}

	//TODO: This need to be refactored
	db.connection.AutoMigrate(&Vehicle{})
	return nil
}
