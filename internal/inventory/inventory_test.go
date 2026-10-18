package inventory

import (
	"errors"
	"testing"
)

func TestInventory_AddDubs(t *testing.T) {
	inv := NewInventory()

	err1 := inv.Add(Equipment{ID: 1, Name: "шлем1"})
	if err1 != nil {
		t.Errorf("Ожидалось успешное добавление первого элемента, получена ошибка: %v", err1)
	}

	_ = inv.Add(Equipment{ID: 1, Name: "шлем2"})
	_ = inv.Add(Equipment{ID: 1, Name: "шлем3"})

	actualLen := inv.Len()
	if actualLen != 1 {
		t.Errorf("Ожидалось количество элементов 1, получено: %d", actualLen)
	}

	name := inv.list[0].Name
	if name != "шлем1" {
		t.Errorf("Ожидалось, что останется шлем1, получен: %s", name)
	}
}

func TestInventory_ForEach(t *testing.T) {
	expectedErr := errors.New("Ошибка во время обхода инвентаря")

	setupInvent := func() *Inventory {
		inv := NewInventory()
		_ = inv.Add(Equipment{ID: 1, Name: "шлем1"})
		_ = inv.Add(Equipment{ID: 2, Name: "шлем2"})
		_ = inv.Add(Equipment{ID: 3, Name: "шлем3"})
		return inv
	}

	tests := []struct {
		name        string
		setupInvent func() *Inventory
		idOnFail    int
		wantErr     error
		countCalls  int
	}{
		{
			name:        "Полный обход инвентаря",
			setupInvent: setupInvent,
			idOnFail:    0,
			wantErr:     nil,
			countCalls:  3,
		},
		{
			name:        "Прерывание на 2 элеменете",
			setupInvent: setupInvent,
			idOnFail:    2,
			wantErr:     expectedErr,
			countCalls:  2,
		},
		{
			name: "Пустой массив",
			setupInvent: func() *Inventory {
				return NewInventory()
			},
			idOnFail:   0,
			wantErr:    nil,
			countCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := tt.setupInvent()
			var calls int

			err := inv.ForEach(func(id int, name string) error {
				calls++
				if tt.idOnFail != 0 && id == tt.idOnFail {
					return tt.wantErr
				}
				return nil
			})

			if err != tt.wantErr {
				t.Errorf("Ожидалась ошибка %v, получена: %v", tt.wantErr, err)
			}
			if calls != tt.countCalls {
				t.Errorf("Ожидалось %d вызовов, получено: %d", tt.countCalls, calls)
			}
		})
	}
}

func TestInventory_ListOrder(t *testing.T) {
	inv := NewInventory()
	_ = inv.Add(Equipment{ID: 1, Name: "шлем1"})
	_ = inv.Add(Equipment{ID: 2, Name: "шлем2"})
	_ = inv.Add(Equipment{ID: 3, Name: "шлем3"})

	if !(inv.list[0].ID == 1 && inv.list[1].ID == 2 && inv.list[2].ID == 3) {
		t.Errorf("Оборудование сохранено в неверном порядке")
	}
}
