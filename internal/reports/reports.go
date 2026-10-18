package reports

import (
	"fmt"
	"io"
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
	equips   EquipmentCounter
}

func NewReportData(v VehicleReporter, e EquipmentCounter) *ReportData {
	return &ReportData{vehicles: v, equips: e}
}

func (r *ReportData) GetReport(w io.Writer) error {
	countVehs := r.vehicles.Len()
	countEquip := r.equips.Len()
	totalEnergy := r.vehicles.TotalEnergy()

	_, err := fmt.Fprintf(w, `В автопарке сейчас %d транспортных средств.
На складе находится %d единиц оборудования.
Общее ежедневное потребление энергии автопарком составляет %d.`, countVehs, countEquip, totalEnergy)
	return err
}
