package vehicles

type Vehicle interface {
	GetID() int
	GetName() string
	GetEnergy() int
	GetComplexity() int
}

type BaseVehicle struct {
	ID         int
	Name       string
	Energy     int
	Complexity int
}

func (b BaseVehicle) GetID() int         { return b.ID }
func (b BaseVehicle) GetName() string    { return b.Name }
func (b BaseVehicle) GetEnergy() int     { return b.Energy }
func (b BaseVehicle) GetComplexity() int { return b.Complexity }

type Scooter struct {
	BaseVehicle
}

func NewScooter(id int, name string, energy, complexity int) *Scooter {
	return &Scooter{
		BaseVehicle: BaseVehicle{Name: name, ID: id, Energy: energy, Complexity: complexity},
	}
}

type Bicycle struct {
	BaseVehicle
}

func NewBicycle(id int, name string, energy, complexity int) *Bicycle {
	return &Bicycle{
		BaseVehicle: BaseVehicle{Name: name, ID: id, Energy: energy, Complexity: complexity},
	}
}

type ElectricScooter struct {
	BaseVehicle
}

func NewElectricScooter(id int, name string, energy, complexity int) *ElectricScooter {
	return &ElectricScooter{
		BaseVehicle: BaseVehicle{Name: name, ID: id, Energy: energy, Complexity: complexity},
	}
}

type ElectricBicycle struct {
	BaseVehicle
}

func NewElectricBicycle(id int, name string, energy, complexity int) *ElectricBicycle {
	return &ElectricBicycle{
		BaseVehicle: BaseVehicle{Name: name, ID: id, Energy: energy, Complexity: complexity},
	}
}

type Unicycle struct {
	BaseVehicle
}

func NewUnicycle(id int, name string, energy, complexity int) *Unicycle {
	return &Unicycle{
		BaseVehicle: BaseVehicle{Name: name, ID: id, Energy: energy, Complexity: complexity},
	}
}
