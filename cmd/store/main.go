package main

import (
	"fmt"
	"log"

	"github.com/radisvaliullin/challenge/pkg/store"
)

func main() {
	fmt.Println("main")

	conf := store.Config{
		Addr: ":8080",
	}
	if err := store.New(conf).Run(); err != nil {
		log.Printf("main: fail: %v", err)
	}
}
