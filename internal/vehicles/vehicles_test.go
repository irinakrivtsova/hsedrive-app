package vehicles

import "testing"

func TestVehicles_ConstructorsAndGetters(t *testing.T) {
	tests := []struct {
		name           string
		vehicle        Vehicle
		wantID         int
		wantName       string
		wantEnergy     int
		wantComplexity int
	}{
		{"Scooter", NewScooter(1, "Scooter1", 10, 2), 1, "Scooter1", 10, 2},
		{"Bicycle", NewBicycle(2, "Bike1", 0, 5), 2, "Bike1", 0, 5},
		{"ElectricScooter", NewElectricScooter(3, "EScooter1", 20, 3), 3, "EScooter1", 20, 3},
		{"ElectricBicycle", NewElectricBicycle(4, "EBike1", 30, 6), 4, "EBike1", 30, 6},
		{"Unicycle", NewUnicycle(5, "Uni1", 5, 8), 5, "Uni1", 5, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vehicle.GetID(); got != tt.wantID {
				t.Errorf("GetID() = %v, ожидалось %v", got, tt.wantID)
			}
			if got := tt.vehicle.GetName(); got != tt.wantName {
				t.Errorf("GetName() = %v, ожидалось %v", got, tt.wantName)
			}
			if got := tt.vehicle.GetEnergy(); got != tt.wantEnergy {
				t.Errorf("GetEnergy() = %v, ожидалось %v", got, tt.wantEnergy)
			}
			if got := tt.vehicle.GetComplexity(); got != tt.wantComplexity {
				t.Errorf("GetComplexity() = %v, ожидалось %v", got, tt.wantComplexity)
			}
		})
	}
}
