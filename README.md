#go console

golang console for embedded code

## init

init as singleton 
```go
c_ptr := console.GetPtr()
```

## use
set output handler
```go
type someHandler struct {}

func (p *someHandler) Output(text []byte) {
    fmt.Println(string(text))
}

c_ptr.Output(h)
```

## variables
assignment of variables to be able to access them
```go
c_ptr.AddBool("bb", &bb)
c_ptr.AddFloat("ff", &ff)
c_ptr.AddInt("ii", &ii)
c_ptr.AddString("ss", &ss)
```

update varibles from code
```go
c_ptr.Exec("set ii 22")
```

## functions
function assignment:
```go
t := &testFunc{}
c_ptr.AddCommand("test", t.ShareMethod)
```

exec function with args
```go
c_ptr.Exec("test")
c_ptr.Exec("test true")
```
## local help
default commands:
```go
c_ptr.Exec("ls")
```

```bash
ls
variables:
bb,ff,ii,ss,
console commands:
ls - prints all variables and commands
get - get value of variables
set - sets value of variable
help - this help
```

created for [http://e154.ru](http://e154.ru/ "e154")
