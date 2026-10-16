package inventory

import "fmt"

type Equipment struct {
	ID int
	Name string 
}

func (e Equipment) GetID() int { return e.ID }
func (e Equipment) GetName() string { return e.Name }

type Inventory struct {
	list []Equipment
}

func NewInventory() *Inventory {
	return &Inventory{}
}

func (inv *Inventory) Add(eq Equipment) error {
	for _, existing := range inv.list {
		if existing.GetID() == eq.GetID() {
			return fmt.Errorf("Оборудование с ID %d уже существует", eq.GetID())
		}
	}

	inv.list = append(inv.list, eq)
	return nil
}

func (inv *Inventory) Len() int {
	return len(inv.list)
}

func (inv *Inventory) ForEach(fn func(id int, name string) error) error {
	for _, v := range inv.list {
		if err := fn(v.GetID(), v.GetName()); err != nil {
			return err
		}
	}
	return nil
}
