package vmgen

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/end-r/efp"
)

// VM ...
type VM struct {
	Name           string
	Author         string
	Receiver       string
	Instructions   map[string]instruction
	ProgramCounter int
	Current        instruction
	stats          *stats
	Stack          *Stack
	Memory         []interface{}
}

// FuelFunction ...
type FuelFunction func(*VM, []byte) int

// ExecuteFunction ...
type ExecuteFunction func(*VM, []byte)

// Instruction for the current FireVM instance
type instruction struct {
	opcode       string
	description  string
	execute      ExecuteFunction
	fuel         int
	fuelFunction FuelFunction
	count        int
}

const prototype = "vmgen.efp"

// CreateVM creates a new FireVM instance
func CreateVM(path string, executes map[string]ExecuteFunction, fuels map[string]FuelFunction) *VM {
	p, errs := efp.PrototypeFile(prototype)
	if errs != nil {
		fmt.Printf("Invalid prototype file\n")
		fmt.Println(errs)
		return nil
	}
	e, errs := p.ValidateFile(path)
	if errs != nil {
		fmt.Printf("Invalid VM file %s\n", path)
		return nil
	}

	var vm VM
	// no need to check for nil: would have errored
	vm.Author = e.FirstField("author").Value()
	vm.Name = e.FirstField("name").Value()
	vm.Receiver = e.FirstField("receiver").Value()

	vm.Instructions = make(map[string]instruction)
	for _, e := range e.Elements("instruction") {
		var i instruction
		i.description = e.FirstField("description").Value()

		// try to get fuel as an integer
		fuel, err := strconv.ParseInt(e.FirstField("fuel").Value(), 10, 64)

		// if not, it's a fuel function
		if err != nil {
			i.fuel = int(fuel)
		} else {
			i.fuelFunction = fuels[e.FirstField("fuel").Value()]
		}

		i.execute = executes[e.FirstField("execute").Value()]

		opcode := e.Parameter(0).Value()
		i.opcode = opcode

		vm.Instructions[opcode] = i
	}

	vm.stats = new(stats)

	vm.Stack = new(Stack)
	return &vm
}

func (vm *VM) createParameters(params []string) []reflect.Value {
	var vals []reflect.Value
	vals = append(vals, reflect.ValueOf(vm))
	for _, p := range params {
		vals = append(vals, reflect.ValueOf(p))
	}
	return vals
}

// ExecuteFile parses opcodes from a file
func (vm *VM) ExecuteFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	bytes := make([]byte, fi.Size())
	_, err = f.Read(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (vm *VM) executeInstruction(opcode string, params []byte) {
	i := vm.Instructions[opcode]
	i.execute(vm, params)
	vm.stats.operations++
	if i.fuelFunction != nil {
		vm.stats.fuelConsumption += i.fuel
	} else {
		vm.stats.fuelConsumption += i.fuelFunction(vm, params)
	}
}

func version() string {
	return fmt.Sprintf("%d.%d.%d", 0, 0, 1)
}
