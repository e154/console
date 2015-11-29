package console

import (
	"fmt"
	"testing"
	"strings"
	"strconv"
)

var (
	ii int64 = 12
	ff float64
	ss string
	bb bool
)

type handler struct {}

func (p *handler) Output(text []byte) {
	fmt.Println(string(text))
}

type testFunc struct {}

func (t *testFunc) ShareMethod(key, value string, help *string) {

	c_ptr := GetPtr()

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

	c_ptr.Printf("run test command with %t argument", arg)
}

func TestCore(t *testing.T) {

	h := &handler{}

	c_ptr := GetPtr()
	c_ptr.Output(h)

	// vars
	c_ptr.AddBool("bb", &bb)
	c_ptr.AddFloat("ff", &ff)
	c_ptr.AddInt("ii", &ii)
	c_ptr.AddString("ss", &ss)

	c_ptr.Exec("set bb true")
	c_ptr.Exec("set ii 22")
	c_ptr.Exec("set ff 1.234")
	c_ptr.Exec("set ss some_string")

	c_ptr.Exec("get ii")
	c_ptr.Exec("ls")

	if !bb { t.Errorf("bb is false") }
	if ff != 1.234 { t.Errorf("ff != 1.234") }
	if ii != 22 { t.Errorf("ii != 22") }
	if ss != "some_string" { t.Errorf("ss != 'some_string'") }

	_t := &testFunc{}
	c_ptr.AddCommand("test", _t.ShareMethod)

	c_ptr.Exec("ls")
	c_ptr.Exec("test")
	c_ptr.Exec("test true")

	c_ptr.Exec("set")
}