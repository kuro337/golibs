package httpinterface

type HTTPInterface interface {
	Hello() string
	Test() bool
}

type SomeInterface struct {
	Method1() string
	Method2() bool 
}
