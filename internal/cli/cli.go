package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"hsedrive-app/internal/inspection"
	"hsedrive-app/internal/inventory"
	"hsedrive-app/internal/park"
	"hsedrive-app/internal/print"
	"hsedrive-app/internal/reports"
	"hsedrive-app/internal/vehicles"
)

type CLI struct {
	park        *park.Park
	inventory   *inventory.Inventory
	parkPrinter *print.Printer
	invPrinter  *print.Printer
	reporter    *reports.ReportData

	reader    *bufio.Scanner
	writer    io.Writer
	errWriter io.Writer
}

func NewCLI(r io.Reader, w io.Writer, ew io.Writer, p *park.Park, i *inventory.Inventory, ppr *print.Printer, ipr *print.Printer, rep *reports.ReportData) *CLI {
	return &CLI{reader: bufio.NewScanner(r), writer: w, errWriter: ew, park: p, inventory: i,
		parkPrinter: ppr, invPrinter: ipr,
		reporter: rep}
}

func (c *CLI) Run() {
	for c.reader.Scan() {
		request := strings.Fields(c.reader.Text())

		if len(request) == 0 {
			continue
		}

		switch request[0] {
		case "add":
			if len(request) < 2 {
				fmt.Fprintln(c.errWriter, "ожидается запрос в формате "+
					"add vehicle <тип> <id> <название> <энергия> <сложность> "+
					"или add equip <id> <название>")
				continue
			}
			switch request[1] {
			case "vehicle":
				c.handleAddVehicle(request[2:])
			case "equip":
				c.handleAddEquip(request[2:])
			}
		case "energy":
			fmt.Fprintf(c.writer, "Суммарное суточное потребление энергии"+
				"парка составляет %d у. е.\n", c.park.TotalEnergy())
		case "beginners":
			c.handleBeginners(request)
		case "park":
			if err := c.parkPrinter.Print(c.writer); err != nil {
				fmt.Fprintf(c.errWriter, "ошибка вывода парка: %v\n", err)
			}
		case "inventory":
			if c.inventory.Len() == 0 {
				fmt.Fprintln(c.writer, "Склад инвентаря пуст.")
			} else if err := c.invPrinter.Print(c.writer); err != nil {
				fmt.Fprintf(c.errWriter, "ошибка вывода инвентаря: %v\n", err)
			}
		case "report":
			c.handleReport(request)
		case "types":
			fmt.Fprintln(c.writer, "Доступные типы транспорта: scooter, bicycle, el-scooter, el-bicycle, unicycle")
		case "help":
			c.printHelp()
		case "exit":
			fmt.Fprintln(c.writer, "Завершение работы.")
			return
		}
	}
}

func (c *CLI) handleAddEquip(request []string) {
	if len(request) < 2 {
		fmt.Fprintln(c.writer, "Формат: add item <id> <name>")
		return
	}
	id, err := strconv.Atoi(request[0])
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка: id должен быть числом")
		return
	}
	name := strings.Join(request[1:], " ")

	newEquip := inventory.Equipment{ID: id, Name: name}

	if err := c.inventory.Add(newEquip); err != nil {
		fmt.Fprintf(c.errWriter, "ошибка добавления: %v\n", err)
		return
	}

	fmt.Fprintln(c.writer, "Оборудование успешно добавлено.")
}

func (c *CLI) handleAddVehicle(request []string) {
	if len(request) < 5 {
		fmt.Fprintln(c.writer, "Формат: add vehicle <id> <type: scooter|bicycle|el-scooter|el-bicycle|unicycle> <name> <energy> <complexity>")
		return
	}

	id, errId := strconv.Atoi(request[0])
	vType := strings.ToLower(request[1])
	name := request[2]
	energy, errEnergy := strconv.Atoi(request[3])
	compl, errCompl := strconv.Atoi(request[4])

	if errId != nil || errEnergy != nil || errCompl != nil {
		fmt.Fprintln(c.errWriter, "ошибка: id, energy и complexity должны быть числами")
		return
	}

	var v vehicles.Vehicle
	switch vType {
	case "scooter":
		v = vehicles.NewScooter(id, name, energy, compl)
	case "bicycle":
		v = vehicles.NewBicycle(id, name, energy, compl)
	case "el-scooter":
		v = vehicles.NewElectricScooter(id, name, energy, compl)
	case "el-bicycle":
		v = vehicles.NewElectricBicycle(id, name, energy, compl)
	case "unicycle":
		v = vehicles.NewUnicycle(id, name, energy, compl)
	default:
		fmt.Fprintln(c.writer, "Неизвестный тип транспорта.")
		return
	}

	if err := c.park.AddVehicle(v); err != nil {
		fmt.Fprintf(c.writer, "Ошибка добавления: %v\n", err)
		return
	}
	fmt.Fprintln(c.writer, "Транспорт успешно прошел техосмотр и добавлен в парк.")
}

func (c *CLI) handleBeginners(request []string) {
	maxCompl := 6
	if len(request) > 1 {
		if val, err := strconv.Atoi(request[1]); err == nil {
			maxCompl = val
		}
	}

	policy := inspection.NewBeginnersInspection(maxCompl)
	serv := inspection.NewServiceCenter(policy)
	list := c.park.Filter(serv)

	if len(list) == 0 {
		fmt.Fprintln(c.writer, "Подходящего транспорта для новичков не найдено.")
		return
	}

	fmt.Fprintln(c.writer, "Список устройств для новичков:")

	for _, val := range list {
		fmt.Fprintf(c.writer, "%s с номером %d\n", val.GetName(), val.GetID())
	}
}

func (c *CLI) handleReport(request []string) {
	dest := c.writer

	var file *os.File

	if len(request) > 1 {
		var err error
		file, err = os.Create(request[1])
		if err != nil {
			fmt.Fprintf(c.errWriter, "ошибка создания файла: %v\n", err)
			return
		}
		dest = file
	}

	err := c.reporter.GetReport(dest)

	if file != nil {
		_ = file.Close()
	}

	if err != nil {
		fmt.Fprintf(c.errWriter, "ошибка формирования отчета: %v\n", err)
		return
	}

	if file != nil {
		fmt.Fprintf(c.writer, "отчет успешно сохранен в файл: %s\n", request[1])
	}
}

func (c *CLI) printHelp() {
	fmt.Fprintln(c.writer, `Команды:
  add vehicle <id> <тип> <название> <энергия> <сложность> — добавить транспорт
  add equip <id> <название>                           — добавить инвентарь на склад
  report [filepath]                                   — сводный отчет
  park                                                — список транспорта
  inventory                                           — список инвентаря
  beginners [max complexity]                          — транспорт для новичков
  types                                               — вывести типы транспортных средств
  help                                                — справка
  exit                                                — выход`)
}
