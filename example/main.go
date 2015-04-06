package main

import (
	"console"
	"fmt"
	"strings"
	"strconv"
)

var (
	ii int64 = 12
	ff float64
	bb bool
	ss string
)

type handler struct {}

func (p *handler) Output(text []byte) {
	fmt.Println(string(text))
}

type testFunc struct {}

func (t *testFunc) ShareMethod(key, value string, help *string) {

	c_ptr := console.GetPtr()

	if value == "help" {
		*help = "test command"
	}

	cmd := strings.Split(value, " ")
	if len(cmd) != 1 {
		c_ptr.Printf("usage: test <bool>")
		return
	}

	arg, err := strconv.ParseBool(cmd[0])
	if err != nil {
		c_ptr.Printf("usage: test <bool>")
		return
	}

	c_ptr.Printf(fmt.Sprintf("run test command with %t argument", arg))
}

func main() {

	h := &handler{}

	c_ptr := console.GetPtr()
	c_ptr.Output(h)

	// vars
	c_ptr.AddBool("bb", &bb)
	c_ptr.AddFloat("ff", &ff)
	c_ptr.AddInt("ii", &ii)
	c_ptr.AddString("ss", &ss)

	c_ptr.Exec("set ii 22")
	c_ptr.Exec("get ii")
	c_ptr.Exec("ls")


	t := &testFunc{}
	c_ptr.AddCommand("test", t.ShareMethod)

	c_ptr.Exec("ls")
	c_ptr.Exec("test")
	c_ptr.Exec("test true")
}
