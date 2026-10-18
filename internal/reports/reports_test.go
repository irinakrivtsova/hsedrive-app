package reports

import (
	"bytes"
	"errors"
	"testing"
)

type stubVehicleReporter struct {
	len    int
	energy int
}

func (s *stubVehicleReporter) Len() int {
	return s.len
}

func (s *stubVehicleReporter) TotalEnergy() int {
	return s.energy
}

type stubEquipmentCounter struct {
	len int
}

func (s *stubEquipmentCounter) Len() int {
	return s.len
}

type stubErrWriter struct{}

func (w stubErrWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("Ошибка записи")
}

func TestReportData_GetReport(t *testing.T) {
	tests := []struct {
		name      string
		vehRep    *stubVehicleReporter
		eqCount   *stubEquipmentCounter
		useErrW   bool
		wantErr   bool
		wantPrint string
	}{
		{
			name: "Успешная генерация отчета",
			vehRep: &stubVehicleReporter{
				len:    10,
				energy: 500,
			},
			eqCount: &stubEquipmentCounter{
				len: 42,
			},
			useErrW:   false,
			wantErr:   false,
			wantPrint: "В автопарке сейчас 10 транспортных средств.\nНа складе находится 42 единицы оборудования.\nОбщее ежедневное потребление энергии автопарком составляет 500.",
		},
		{
			name: "Ошибка записи (io.Writer)",
			vehRep: &stubVehicleReporter{
				len:    5,
				energy: 100,
			},
			eqCount: &stubEquipmentCounter{
				len: 10,
			},
			useErrW: true,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd := NewReportData(tt.vehRep, tt.eqCount)

			var buf bytes.Buffer
			var w interface {
				Write(p []byte) (n int, err error)
			} = &buf
			if tt.useErrW {
				w = stubErrWriter{}
			}

			err := rd.GetReport(w)

			if (err != nil) != tt.wantErr {
				t.Errorf("Ожидалось наличие ошибки: %v, получена: %v", tt.wantErr, err)
			}

			if !tt.wantErr && !tt.useErrW {
				if got := buf.String(); got != tt.wantPrint {
					t.Errorf("Ожидался вывод:\n%q\nПолучено:\n%q", tt.wantPrint, got)
				}
			}
		})
	}
}
