package main

import (
	"fmt"

	"github.com/ThiagoDonadel/loan-management/app/registry"
	"github.com/ThiagoDonadel/loan-management/infra"
)

func main() {

	fmt.Println("HELLO")
	if err := infra.ConnectToDatabase(); err != nil {
		panic(err)
	}
	registry.Initialialize(infra.DBConnection)
	infra.StartGinServer()
}
