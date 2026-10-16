package park

import (
	"fmt"
	"hsedrive-app/internal/vehicles"	
)

type Park struct {
	service Inspection
	list []vehicles.Vehicle
}

type Inspection interface {
	ApplyPolicy(vehicle vehicles.Vehicle) bool
}

func NewPark(service Inspection) *Park {
	return &Park{
		service: service,
	}
}

func (a *Park) TotalEnergy() int {
	sum := 0
	for _, val := range a.list {
		sum += val.GetEnergy()
	}
	return sum
}

func (a *Park) AddVehicle(newVehicle vehicles.Vehicle) error {
	for _, existing := range a.list {
		if existing.GetID() == newVehicle.GetID() {
			return fmt.Errorf("Транспорт с ID %d уже существует в парке", newVehicle.GetID())
		}
	}

	if !a.service.ApplyPolicy(newVehicle) {
		return fmt.Errorf("Транспорт #%d не прошел техосмотр", newVehicle.GetID())
	}

	a.list = append(a.list, newVehicle)
	return nil
}

func (a *Park) Len() int {
	return len(a.list)
}

func (a *Park) ForEach(fn func(id int, name string) error) error {
	for _, v := range a.list {
		if err := fn(v.GetID(), v.GetName()); err != nil {
			return err
		}
	}
	return nil
}

func (a *Park) Filter(inspection Inspection) []vehicles.Vehicle {
	var available []vehicles.Vehicle
	for _, v := range a.list {
		if inspection.ApplyPolicy(v) {
			available = append(available, v)
		}
	}
	return available
}

