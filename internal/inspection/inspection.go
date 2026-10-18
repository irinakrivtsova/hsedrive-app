package inspection

import "hsedrive-app/internal/vehicles"

type CheckPolicy interface {
	Check(vehicles.Vehicle) bool
}

type TechInspection struct {
	normalCompl  int
	normalEnergy int
}

func NewTechInspection(energy int, compl int) *TechInspection {
	return &TechInspection{normalCompl: compl, normalEnergy: energy}
}

func (i *TechInspection) Check(vehicle vehicles.Vehicle) bool {
	return !(vehicle.GetEnergy() > i.normalEnergy || vehicle.GetComplexity() < i.normalCompl)
}

type BeginnersInspection struct {
	normalCompl int
}

func NewBeginnersInspection(compl int) *BeginnersInspection {
	return &BeginnersInspection{normalCompl: compl}
}

func (i *BeginnersInspection) Check(vehicle vehicles.Vehicle) bool {
	return vehicle.GetComplexity() >= i.normalCompl
}

type ServiceCenter struct {
	policy CheckPolicy
}

func NewServiceCenter(policy CheckPolicy) *ServiceCenter {
	return &ServiceCenter{policy}
}

func (c *ServiceCenter) ApplyPolicy(vehicle vehicles.Vehicle) bool {
	return c.policy.Check(vehicle)
}
