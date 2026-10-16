package cli

import (
	"bytes"
	"strings"
	"testing"

	"hsedrive-app/internal/inventory"
	"hsedrive-app/internal/park"
	"hsedrive-app/internal/print"
	"hsedrive-app/internal/reports"
	"hsedrive-app/internal/vehicles"
)

type StubInspector struct {
	result bool
}

func (s *StubInspector) ApplyPolicy(v vehicles.Vehicle) bool {
	return s.result
}

func TestCLI_Commands(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		inspResult bool
		wantOut    string
		wantErrOut string
		wantPark   int
		wantInv    int
	}{
		{
			name:       "Exit",
			input:      "exit\n",
			inspResult: true,
			wantOut:    "Завершение работы.",
		},
		{
			name:       "Help",
			input:      "help\nexit\n",
			inspResult: true,
			wantOut:    "Команды:",
		},
		{
			name:       "Add Vehicle Success",
			input:      "add vehicle 1 scooter test 10 5\nexit\n",
			inspResult: true,
			wantOut:    "успешно прошел техосмотр",
			wantPark:   1,
		},
		{
			name:       "Add Vehicle Failed Inspection",
			input:      "add vehicle 1 scooter test 10 5\nexit\n",
			inspResult: false,
			wantOut:    "Ошибка добавления",
			wantPark:   0,
		},
		{
			name:       "Add Vehicle Missing Args",
			input:      "add vehicle 1 scooter test\nexit\n",
			inspResult: true,
			wantOut:    "Формат: add vehicle",
		},
		{
			name:       "Add Vehicle Invalid Numbers",
			input:      "add vehicle abc scooter test 10 5\nexit\n",
			inspResult: true,
			wantErrOut: "Ошибка: id, energy и complexity должны быть числами",
		},
		{
			name:       "Add Vehicle Unknown Type",
			input:      "add vehicle 1 ufo test 10 5\nexit\n",
			inspResult: true,
			wantOut:    "Неизвестный тип транспорта.",
		},
		{
			name:       "Add Equip Success",
			input:      "add equip 1 helmet\nexit\n",
			inspResult: true,
			wantOut:    "успешно добавлено",
			wantInv:    1,
		},
		{
			name:       "Add Equip Missing Args",
			input:      "add equip\nexit\n",
			inspResult: true,
			wantOut:    "Формат: add item",
		},
		{
			name:       "Add Equip Invalid ID",
			input:      "add equip abc helmet\nexit\n",
			inspResult: true,
			wantOut:    "Ошибка: id должен быть числом",
		},
		{
			name:       "Print Report",
			input:      "report\nexit\n",
			inspResult: true,
			wantOut:    "There are a total of",
		},
		{
			name:       "Print Energy",
			input:      "energy\nexit\n",
			inspResult: true,
			wantOut:    "Суммарное суточное потребление",
		},
		{
			name:       "Beginners Default",
			input:      "beginners\nexit\n",
			inspResult: true,
			wantOut:    "Подходящего транспорта для новичков не найдено.",
		},
		{
			name:       "Beginners Custom Level",
			input:      "beginners 4\nexit\n",
			inspResult: true,
			wantOut:    "Подходящего транспорта для новичков не найдено.",
		},
		{
			name:       "Print Park",
			input:      "park\nexit\n",
			inspResult: true,
			wantOut:    "", 
		},
		{
			name:       "Print Inventory",
			input:      "inventory\nexit\n",
			inspResult: true,
			wantOut:    "Склад инвентаря пуст.",
		},
		{
			name:       "Print Types",
			input:      "types\nexit\n",
			inspResult: true,
			wantOut:    "Доступные типы транспорта:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}
			errOut := &bytes.Buffer{}

			stubInsp := &StubInspector{result: tt.inspResult}
			p := park.NewPark(stubInsp)
			inv := inventory.NewInventory()

			pPrinter := print.NewPrinter(p)
			iPrinter := print.NewPrinter(inv)
			rep := reports.NewReportData(p, inv)

			c := NewCLI(in, out, errOut, p, inv, pPrinter, iPrinter, rep)
			c.Run()

			if tt.wantOut != "" && !strings.Contains(out.String(), tt.wantOut) {
				t.Errorf("Ожидалось в stdout: %s, получено: %s", tt.wantOut, out.String())
			}
			if tt.wantErrOut != "" && !strings.Contains(errOut.String(), tt.wantErrOut) {
				t.Errorf("Ожидалось в stderr: %s, получено: %s", tt.wantErrOut, errOut.String())
			}
			if p.Len() != tt.wantPark {
				t.Errorf("Ожидалось транспорта: %d, получено: %d", tt.wantPark, p.Len())
			}
			if inv.Len() != tt.wantInv {
				t.Errorf("Ожидалось инвентаря: %d, получено: %d", tt.wantInv, inv.Len())
			}
		})
	}
}
