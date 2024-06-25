package main

import (
	"fmt"
	"github.com/lukasFindura/gocliselect"
)

func main() {
	menu := gocliselect.NewMenu("Chose a colour", 0)

	menu.AddItem("Red", "red")
	menu.AddItem("Blue", "blue")
	menu.AddItem("Green", "green")
	menu.AddItem("Yellow", "yellow")
	menu.AddItem("Cyan", "cyan")

	if choice := menu.Display(menu); choice != nil {
		fmt.Printf("Choice: %s\n", choice.Text)
	}
}