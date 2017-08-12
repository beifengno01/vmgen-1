package vmgen

import (
	"os"
	"strconv"
)

// GenerateReadMe generates a read me
func (vm *VM) GenerateReadMe(name string) {
	f, err := os.Create(name)
	if err != nil {
		return
	}
	defer f.Close()
	f.Write([]byte("| " + "Opcode"))
	f.Write([]byte("| " + "Description"))
	f.Write([]byte("| " + "Fuel"))
	f.Write([]byte("|\n"))
	for _, v := range vm.Instructions {
		f.Write([]byte("| " + v.opcode))
		f.Write([]byte("| " + v.description))
		f.Write([]byte("| " + strconv.Itoa(v.fuel)))
		f.Write([]byte("|\n"))
	}
	f.Sync()
}