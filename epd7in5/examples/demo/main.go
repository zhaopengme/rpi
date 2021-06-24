package main

import (
	"fmt"
	"github.com/gandaldf/rpi/epd7in5/epd"
	"time"
)

func main() {
	e := epd.CreateEpd()
	defer e.Close()
	defer e.Clear()
	e.Init()
	e.Clear()

	fmt.Printf("Display\n")
	e.DisplayBlack(MyImg)
	fmt.Printf("sleeping\n")
	time.Sleep(5 * time.Second)
}
