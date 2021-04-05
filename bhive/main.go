package main

import (
  _	"embed"
)

//go:embed scale.py
var pyScale []byte

//go:embed hx711.py
var pyHx711 []byte

func main() {
	print(string(pyScale))
	print(string(pyHx711))
}
