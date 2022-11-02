package models

type Vehicle struct {
	// gorm.Model
	ID               uint64 `gorm:"primary_index;auto_increment" json:"id"`
	Description      string `gorm:"not null" json:"description"`
	Price            uint   `gorm:"not null" json:"price"`
	Mileage          uint   `json:"mileage"`
	ManufacturedYear uint   `json:"manufactured"`
}

type VehicleService interface {
	Save(Vehicle) error
	Update(Vehicle) error
	Delete(Vehicle) error
	FindAll() []Vehicle
}

type vehicleService struct {
	vehicleDB VehicleDB
}

func NewVehicleService(vehicleDB VehicleDB) VehicleService {
	return &vehicleService{
		vehicleDB: vehicleDB,
	}
}

func (service *vehicleService) Save(v Vehicle) error {
	return service.vehicleDB.Save(v)
}

func (service *vehicleService) Update(v Vehicle) error {
	return service.vehicleDB.Update(v)
}

func (service *vehicleService) Delete(v Vehicle) error {
	return service.vehicleDB.Delete(v)
}

func (service *vehicleService) FindAll() []Vehicle {
	return service.vehicleDB.FindAll()
}
