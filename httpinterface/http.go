package httpinterface

type HTTPInterface interface {
	Hello() string
	Test() bool
}

type SomeInterface interface {
	Method1() string
	Method2() bool
}
