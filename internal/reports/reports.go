package reports

import (
	"io"
	"fmt"
)

type VehicleReporter interface {
	Len() int
	TotalEnergy() int
}

type EquipmentCounter interface {
	Len() int
}

type ReportData struct {
	vehicles VehicleReporter
	equips EquipmentCounter
}

func NewReportData(v VehicleReporter, e EquipmentCounter) *ReportData {
	return &ReportData{vehicles: v, equips: e}
}

func (r *ReportData) GetReport(w io.Writer) error {
	countVehs := r.vehicles.Len()
	countEquip := r.equips.Len()
	totalEnergy := r.vehicles.TotalEnergy()

	_, err := fmt.Fprintf(w, `There are a total of %d vehicles in the fleet.
There are %d items in the warehouse that are subject to inventory.
The total daily energy consumption of the fleet is %d.`, countVehs, countEquip, totalEnergy)
	return err
}
