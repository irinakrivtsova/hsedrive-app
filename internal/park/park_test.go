package park

import (
	"errors"
	"hsedrive-app/internal/vehicles"
	"testing"
)

type SpyInspector struct {
	callsCount int
	stubResult bool
}

func (s *SpyInspector) ApplyPolicy(v vehicles.Vehicle) bool {
	s.callsCount++
	return s.stubResult
}

type StubInspector struct {
	stubResult bool
}

func (s *StubInspector) ApplyPolicy(v vehicles.Vehicle) bool {
	return s.stubResult
}

type MockInspection struct {
	policyFunc func(v vehicles.Vehicle) bool
}

func (m *MockInspection) ApplyPolicy(v vehicles.Vehicle) bool {
	return m.policyFunc(v)
}

func TestPark_AddVehicle_CallsInspection(t *testing.T) {
	spy := &SpyInspector{stubResult: true}
	p := NewPark(spy)
	v := vehicles.NewScooter(1, "Scooter", 10, 3)

	err := p.AddVehicle(v)

	if err != nil {
		t.Errorf("Не ожидалась ошибка, получено: %v", err)
	}
	if spy.callsCount != 1 {
		t.Errorf("Ожидался 1 вызов техосмотра, получено: %d", spy.callsCount)
	}
}

func TestPark_AddVehicle_RejectsFailed(t *testing.T) {
	spy := &SpyInspector{stubResult: false}
	p := NewPark(spy)
	v := vehicles.NewScooter(1, "Scooter", 10, 3)

	err := p.AddVehicle(v)

	if err == nil {
		t.Errorf("Ожидалась ошибка при провале техосмотра, но ее нет")
	}
	if len(p.list) != 0 {
		t.Errorf("Транспорт не должен был добавиться в массив парка")
	}
}

func TestPark_Filter(t *testing.T) {
	spyAllow := &SpyInspector{stubResult: true}
	spyReject := &SpyInspector{stubResult: false}

	tests := []struct {
		name       string
		policy     *SpyInspector
		setupCount int
		wantCount  int
	}{
		{
			name:       "Политика пропускает все",
			policy:     spyAllow,
			setupCount: 3,
			wantCount:  3,
		},
		{
			name:       "Политика блокирует все",
			policy:     spyReject,
			setupCount: 3,
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPark(spyAllow)
			for i := 0; i < tt.setupCount; i++ {
				_ = p.AddVehicle(vehicles.NewScooter(i, "Test", 10, 3))
			}

			filtered := p.Filter(tt.policy)

			if len(filtered) != tt.wantCount {
				t.Errorf("Ожидалось %d элементов после фильтрации, получено %d", tt.wantCount, len(filtered))
			}
		})
	}
}

func TestPark_Filter_Extended(t *testing.T) {
	v1 := vehicles.NewScooter(1, "Scooter 1", 10, 3)
	v2 := vehicles.NewBicycle(2, "Bike 2", 0, 1)
	v3 := vehicles.NewElectricScooter(3, "El-Scooter 3", 50, 8)
	stub := &StubInspector{stubResult: true}

	tests := []struct {
		name       string
		setupPark  func() *Park
		inspection *MockInspection
		wantIDs    []int
	}{
		{
			name: "Частичная выборка (пропускаем только ID 2)",
			setupPark: func() *Park {
				p := NewPark(stub)
				_ = p.AddVehicle(v1)
				_ = p.AddVehicle(v2)
				_ = p.AddVehicle(v3)
				return p
			},
			inspection: &MockInspection{
				policyFunc: func(v vehicles.Vehicle) bool {
					return v.GetID() == 2
				},
			},
			wantIDs: []int{2},
		},
		{
			name: "Никто не прошел",
			setupPark: func() *Park {
				p := NewPark(stub)
				_ = p.AddVehicle(v1)
				return p
			},
			inspection: &MockInspection{
				policyFunc: func(v vehicles.Vehicle) bool { return false },
			},
			wantIDs: []int{},
		},
		{
			name: "Пустой парк",
			setupPark: func() *Park {
				return NewPark(stub)
			},
			inspection: &MockInspection{
				policyFunc: func(v vehicles.Vehicle) bool { return true },
			},
			wantIDs: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.setupPark()

			got := p.Filter(tt.inspection)

			if len(got) != len(tt.wantIDs) {
				t.Fatalf("Ожидалось элементов: %d, получено: %d", len(tt.wantIDs), len(got))
			}

			for i, v := range got {
				if v.GetID() != tt.wantIDs[i] {
					t.Errorf("На позиции %d ожидался ID %d, получен ID %d", i, tt.wantIDs[i], v.GetID())
				}
			}
		})
	}
}

func TestPark_TotalEnergy(t *testing.T) {
	tests := []struct {
		name       string
		setupPark  func() *Park
		wantEnergy int
	}{
		{
			name: "Парк с транспортом",
			setupPark: func() *Park {
				stub := &StubInspector{stubResult: true}
				p := NewPark(stub)

				energy := 5
				_ = p.AddVehicle(vehicles.NewBicycle(1, "test 1", energy, 5))
				_ = p.AddVehicle(vehicles.NewBicycle(2, "test 2", energy, 5))
				_ = p.AddVehicle(vehicles.NewBicycle(3, "test 3", energy, 5))

				return p
			},
			wantEnergy: 15,
		},
		{
			name: "Пустой парк",
			setupPark: func() *Park {
				stub := &StubInspector{stubResult: true}
				return NewPark(stub)
			},
			wantEnergy: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.setupPark()

			res := p.TotalEnergy()

			if res != tt.wantEnergy {
				t.Errorf("Ожидалось суммарное потребление энергии %d, получено %d", tt.wantEnergy, res)
			}
		})
	}
}

func TestLen(t *testing.T) {
	stub := &StubInspector{stubResult: true}
	p := NewPark(stub)

	n := 3
	for i := 1; i <= n; i++ {
		_ = p.AddVehicle(vehicles.NewBicycle(i, "test", 5, 5))
	}

	res := p.Len()
	if res != n {
		t.Errorf("Ожидалось количество элементов в парке %d, получено %d", n, res)
	}
}

func TestPark_ForEach(t *testing.T) {
	expectedErr := errors.New("Ошибка базы данных")

	makeFilledPark := func() *Park {
		stub := &StubInspector{stubResult: true}
		p := NewPark(stub)
		for i := 1; i <= 3; i++ {
			_ = p.AddVehicle(vehicles.NewBicycle(i, "test", 5, 5))
		}
		return p
	}

	tests := []struct {
		name      string
		setupPark func() *Park
		failOnID  int
		wantErr   error
		wantCalls int
	}{
		{
			name:      "Полный обход без ошибок",
			setupPark: makeFilledPark,
			failOnID:  0,
			wantErr:   nil,
			wantCalls: 3,
		},
		{
			name:      "Прерывание при ошибке",
			setupPark: makeFilledPark,
			failOnID:  2,
			wantErr:   expectedErr,
			wantCalls: 2,
		},
		{
			name: "Пустой парк",
			setupPark: func() *Park {
				return NewPark(&StubInspector{stubResult: true})
			},
			failOnID:  0,
			wantErr:   nil,
			wantCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.setupPark()
			var calls int

			err := p.ForEach(func(id int, name string) error {
				calls++
				if tt.failOnID != 0 && id == tt.failOnID {
					return tt.wantErr
				}
				return nil
			})

			if err != tt.wantErr {
				t.Errorf("Ожидалась ошибка %v, получена %v", tt.wantErr, err)
			}
			if calls != tt.wantCalls {
				t.Errorf("Ожидалось %d вызовов, получено %d", tt.wantCalls, calls)
			}
		})
	}
}
