package console

import (
	"strings"
	"fmt"
	"strconv"
)

const (
	BOOL int = iota
	INT64
	FLOAT64
	STRING
	history_size = 512
)

type variable struct {
	typ_e int
	b *bool
	i *int64
	f *float64
	s *string
}

type console struct {
	// commands store
	commands map[string]func(string, string, *string)

	// variables store
	variables map[string]*variable

	// command history
	history []string

	// output handler
	output Handler
}

var instantiated *console = nil

func GetPtr() *console {
	if instantiated == nil {
		instantiated = &console{
			commands:		make(map[string]func(string, string, *string)),
			variables:		make(map[string]*variable),
			history:		make([]string, history_size),
			output:			nil,
		};

		instantiated.AddCommand("ls", instantiated.ls)
		instantiated.AddCommand("get", instantiated.get)
		instantiated.AddCommand("set", instantiated.set)
		instantiated.AddCommand("help", instantiated.help)
	}
	return instantiated;
}

//-------------------------------------------------------------------------------------
// if command exist, remove first
// after add
func (c *console) AddCommand(name string, h func(key, value string, result *string)) {
	c.RemoveCommand(name)
	c.commands[name] = h
}
// remove command from command base
func (c *console) RemoveCommand(name string) {
	if c.commands[name] != nil {
		delete(c.commands, name)
	}
}

func (c *console) Printf(text string) {
	if text != "" {
		c.output.Output(([]byte(text)))
	}
}

func (c *console) Exec(command string) {

	c.Printf(command)

	if command == "" {
		return
	}

	pos := strings.Split(command, " ")
	var key, value string

	key = pos[0]
	if len(pos) > 1 {
		value = command[len(key)+1:]
	}

	if c.commands[key] != nil {
		var result string
		c.commands[key](key, value, &result)
		c.Printf(result)
	} else {
		c.Printf(fmt.Sprintf("unknown command %s", command))
	}

	return
}

//-------------------------------------------------------------------------------------
// setters
//-------------------------------------------------------------------------------------
func (c *console) AddBool(name string, value *bool) {

	v := new(variable)
	v.typ_e = BOOL
	v.b = value
	c.variables[name] = v
}

func (c *console) AddInt(name string, value *int64) {

	v := new(variable)
	v.typ_e = INT64
	v.i = value
	c.variables[name] = v
}

func (c *console) AddFloat(name string, value *float64) {

	v := new(variable)
	v.typ_e = FLOAT64
	v.f = value
	c.variables[name] = v
}

func (c *console) AddString(name string, value *string) {

	v := new(variable)
	v.typ_e = STRING
	v.s = value
	c.variables[name] = v
}

//-------------------------------------------------------------------------------------
// register commands
//-------------------------------------------------------------------------------------
func (c *console) set(name, value string, help *string) {

	if value == "help" {
		*help = "sets value of variable"
		return
	}

	cmd := strings.Split(value, " ")
	name = cmd[0]

//	if c.variables[name] == nil {
//		return "usage: set <variable> <value>\n"
//	}

	// if non string? and 1 arg, bad command
	if len(cmd) != 2 && c.variables[name] == nil && c.variables[name].typ_e != STRING {
		c.Printf("usage: set <variable> <value>\n")
		return
	}

	// need for empty field
	if len(cmd) == 1 && c.variables[name].typ_e == STRING {
		value = ""
	} else {
		value = cmd[1]
	}

	// check if exist
	if c.variables[name] == nil {
		// variable not found
		c.Printf(fmt.Sprintf("unknown variable %s\n", name))
		return
	}

	// variable is found
	switch c.variables[name].typ_e {
	case BOOL:
		if value == "1" {
			*c.variables[name].b = true
		} else if value == "0" {
			*c.variables[name].b = false
		} else if value == "true" {
			*c.variables[name].b = true
		} else if value == "false" {
			*c.variables[name].b = false
		}
		c.Printf(fmt.Sprintf("%s set to %t", name, *c.variables[name].b))
	case INT64:
		*c.variables[name].i, _ = strconv.ParseInt(value, 10, 64)
		c.Printf(fmt.Sprintf("%s set to %d", name, *c.variables[name].i))
	case FLOAT64:
		*c.variables[name].f, _ = strconv.ParseFloat(value, 64)
		c.Printf(fmt.Sprintf("%s set to %f", name, *c.variables[name].f))
	case STRING:
		*c.variables[name].s = value
		c.Printf(fmt.Sprintf("%s set to %s", name, *c.variables[name].s))
	default:
		c.Printf(fmt.Sprintf("unknown variable type %s", name))
	}
}

func (c *console) get(name, value string, help *string) {

	if value == "help" {
		*help = "get value of variables"
		return
	}

	cmd := strings.Split(value, " ")
	if value == "" || len(cmd) > 1 {
		c.Printf("usage: get <variable>")
		return
	}

	if c.variables[value] == nil {
		c.Printf(fmt.Sprintf("unknown variable %s", value))
		return
	}

	switch c.variables[value].typ_e {
	case BOOL:
		c.Printf(fmt.Sprintf("%t", *c.variables[value].b))
	case INT64:
		c.Printf(fmt.Sprintf("%d", *c.variables[value].i))
	case FLOAT64:
		c.Printf(fmt.Sprintf("%f", *c.variables[value].f))
	case STRING:
		c.Printf(fmt.Sprintf("%s", *c.variables[value].s))
	default:
		c.Printf(fmt.Sprintf("unknown variable type %s", value))
	}
}

func (c *console) ls(key, value string, help *string) {

	if value == "help" {
		*help = "prints all variables and commands"
		return
	}

	// variables
	v_size := len(c.variables)
	if v_size > 0 {
		c.Printf(fmt.Sprintf("variables:"))
		i := 0
		var v string
		for name, _ := range c.variables {
			v += name
			if i != v_size {
				v += ","
			}
			i++
		}
		c.Printf(v)
	}


	// commands
	var h string
	c.help("", "", &h)
	c.Printf(h)
}

func (c *console) help(key, value string, help *string) {

	if value == "help" {
		*help = "this help"
		return
	}

	// commands
	v_size := len(c.commands)
	if v_size > 0 {
		c.Printf(fmt.Sprintf("console commands:"))
		for name, command := range c.commands {
			var h string
			command("", "help", &h)
			c.Printf(fmt.Sprintf("%s - %s", name, h))
		}
	}
}

func (c *console) Output(h Handler) {
	c.output = h
}
