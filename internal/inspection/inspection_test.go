package inspection

import (
	"hsedrive-app/internal/vehicles"
	"testing"
)

type stubVehicle struct {
	energy     int
	complexity int
}

func (s *stubVehicle) GetID() int {
	return 1
}

func (s *stubVehicle) GetName() string {
	return "test vehicle"
}

func (s *stubVehicle) GetEnergy() int {
	return s.energy
}

func (s *stubVehicle) GetComplexity() int {
	return s.complexity
}

type spyPolicy struct {
	callsCount int
	stubResult bool
}

func (s *spyPolicy) Check(v vehicles.Vehicle) bool {
	s.callsCount++
	return s.stubResult
}

func TestTechInspection_Check(t *testing.T) {
	ti := NewTechInspection(100, 5)

	tests := []struct {
		name string
		v    *stubVehicle
		want bool
	}{
		{
			name: "Подходит по всем параметрам",
			v:    &stubVehicle{energy: 50, complexity: 6},
			want: true,
		},
		{
			name: "Точные границы",
			v:    &stubVehicle{energy: 100, complexity: 5},
			want: true,
		},
		{
			name: "Превышена энергия",
			v:    &stubVehicle{energy: 101, complexity: 6},
			want: false,
		},
		{
			name: "Слишком низкая сложность",
			v:    &stubVehicle{energy: 50, complexity: 4},
			want: false,
		},
		{
			name: "Не подходит по обоим параметрам",
			v:    &stubVehicle{energy: 120, complexity: 2},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ti.Check(tt.v); got != tt.want {
				t.Errorf("Ожидалось %v, получено %v", tt.want, got)
			}
		})
	}
}

func TestBeginnersInspection_Check(t *testing.T) {
	bi := NewBeginnersInspection(3)

	tests := []struct {
		name string
		v    *stubVehicle
		want bool
	}{
		{
			name: "Сложность выше нормы",
			v:    &stubVehicle{complexity: 5},
			want: true,
		},
		{
			name: "Сложность на нижней границе",
			v:    &stubVehicle{complexity: 3},
			want: true,
		},
		{
			name: "Сложность ниже нормы",
			v:    &stubVehicle{complexity: 2},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bi.Check(tt.v); got != tt.want {
				t.Errorf("Ожидалось %v, получено %v", tt.want, got)
			}
		})
	}
}

func TestServiceCenter_ApplyPolicy(t *testing.T) {
	tests := []struct {
		name       string
		stubResult bool
	}{
		{
			name:       "Политика разрешает",
			stubResult: true,
		},
		{
			name:       "Политика отклоняет",
			stubResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy := &spyPolicy{stubResult: tt.stubResult}
			sc := NewServiceCenter(spy)
			v := &stubVehicle{}

			got := sc.ApplyPolicy(v)

			if got != tt.stubResult {
				t.Errorf("Ожидалось %v, получено %v", tt.stubResult, got)
			}
			if spy.callsCount != 1 {
				t.Errorf("Ожидался 1 вызов политики, получено %d", spy.callsCount)
			}
		})
	}
}
