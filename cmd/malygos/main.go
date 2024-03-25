package main

import (
	"github.com/nrz-incubator/malygos/pkg/malygos"
	_ "go.uber.org/automaxprocs"
)

func main() {
	app := malygos.New()
	if err := app.Run(); err != nil {
		panic(err)
	}
}
