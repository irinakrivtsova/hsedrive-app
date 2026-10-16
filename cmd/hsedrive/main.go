package main

import (
	"os"

	"hsedrive-app/internal/cli"
	"hsedrive-app/internal/inventory"
	"hsedrive-app/internal/inspection"
	"hsedrive-app/internal/park"
	"hsedrive-app/internal/print"
	"hsedrive-app/internal/reports"
)

func main() {
	// создание парка
	simplicity := 3
	energyConsumption := 8
	// указывать ли мне эту зависимость как TechInspection или она имеет тип CheckPolicy
	inspPolicy := inspection.NewTechInspection(energyConsumption, simplicity)
	serviceCenter := inspection.NewServiceCenter(inspPolicy)
	storage := park.NewPark(serviceCenter)

	// создание инвентаря
	invStorage := inventory.NewInventory()

	// отчеты и печать
	reporter := reports.NewReportData(storage, invStorage)
	printerVehs := print.NewPrinter(storage)
	printerEquips := print.NewPrinter(invStorage)

	// command line interface
	cli := cli.NewCLI(os.Stdin, os.Stdout, os.Stderr, storage, invStorage, printerVehs, printerEquips, reporter)
	cli.Run()
}