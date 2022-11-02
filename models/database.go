package models

//TODO: Candidate on refactoring...

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type VehicleDB interface {
	Save(vehicle Vehicle) error
	Update(vehicle Vehicle) error
	Delete(vehicle Vehicle) error
	FindAll() []Vehicle
	CloseConnection()
	DestructiveReset() error
}

type dbConn struct {
	connection *gorm.DB
}

//TODO: Make it parameterizable and secure.
var connectionString = "host=localhost port=5432 user=admin password=qwerty dbname=car_park_dev sslmode=disable"

func NewVehicleDB() VehicleDB {
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		panic("Failed to connect to database")
	}
	db.LogMode(true)
	db.AutoMigrate(&Vehicle{})
	return &dbConn{
		connection: db,
	}
}

func (db *dbConn) CloseConnection() {
	err := db.connection.Close()
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
	return db.connection.Delete(&v).Error
}

func (db *dbConn) FindAll() []Vehicle {
	var vehicles []Vehicle
	// in case of appearence foreign keys (nested structs in Vehicle table)
	// change below to db.connection.Set("gorm:auto_preload", true).Find(&vehicles)
	db.connection.Find(&vehicles)
	return vehicles
}

// DROPDATABASE!
func (db *dbConn) DestructiveReset() error {
	err := db.connection.DropTableIfExists(&Vehicle{}).Error
	if err != nil {
		return err
	}

	//TODO: This need to be refactored
	db.connection.AutoMigrate(&Vehicle{})
	return nil
}
