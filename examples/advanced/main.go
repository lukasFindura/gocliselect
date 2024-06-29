package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/buger/goterm"
	"github.com/lukasFindura/gocliselect"
)

type Item struct {
    Color string  `json:"color"`
    Items *[]Item `json:"items,omitempty"` // Using a pointer to differentiate between null and an empty slice
}

type ColorData struct {
    Color string `json:"color"`
    Items []Item `json:"items"`
}

// Parse ColorData struct recursively
func createMenu(color string, items []Item, level int) *gocliselect.Menu {
    var menu *gocliselect.Menu = gocliselect.NewMenu(color, level)
    for _, item := range items {
        if item.Items != nil {
            subMenu := createMenu(item.Color, *item.Items, level+1)
            menu.AddItemMenu(item.Color, strings.ToLower(item.Color), subMenu)
        } else {
            menu.AddItem(item.Color, strings.ToLower(item.Color))
        }
    }

    return menu
}

func main() {
    jsonData := `{
      "color": "colors",
      "items": [
        {
          "color": "Red",
          "items": [
            {
              "color": "Light Red",
              "items": [
                { "color": "Pink", "items": null },
                { "color": "Salmon", "items": null },
                { "color": "Rose", "items": null }
              ]
            },
            {
              "color": "Medium Red",
              "items": [
                { "color": "Cherry", "items": null },
                { "color": "Scarlet", "items": null },
                { "color": "Ruby", "items": null }
              ]
            },
            {
              "color": "Dark Red",
              "items": [
                { "color": "Maroon", "items": null },
                { "color": "Burgundy", "items": null },
                { "color": "Crimson", "items": null }
              ]
            }
          ]
        },
        {
          "color": "Blue",
          "items": [
            {
              "color": "Light Blue",
              "items": [
                { "color": "Sky Blue", "items": null },
                { "color": "Powder Blue", "items": null },
                { "color": "Baby Blue", "items": null }
              ]
            },
            {
              "color": "Medium Blue",
              "items": [
                { "color": "Azure", "items": null },
                { "color": "Cobalt", "items": null },
                { "color": "Sapphire", "items": null }
              ]
            },
            {
              "color": "Dark Blue",
              "items": [
                { "color": "Navy", "items": null },
                { "color": "Midnight Blue", "items": null },
                { "color": "Royal Blue", "items": null }
              ]
            }
          ]
        },
        {
          "color": "Green",
          "items": [
            {
              "color": "Light Green",
              "items": [
                { "color": "Mint", "items": null },
                { "color": "Lime", "items": null },
                { "color": "Seafoam", "items": null }
              ]
            },
            {
              "color": "Medium Green",
              "items": [
                { "color": "Jade", "items": null },
                { "color": "Moss", "items": null },
                { "color": "Sage", "items": null }
              ]
            },
            {
              "color": "Dark Green",
              "items": [
                { "color": "Forest Green", "items": null },
                { "color": "Emerald", "items": null },
                { "color": "Olive", "items": null }
              ]
            }
          ]
        }
      ]
    }`

    gocliselect.Cursor.ItemPrompt = "❯"
    gocliselect.Cursor.SubMenuPrompt = "❯"
    gocliselect.Cursor.ItemColor = goterm.YELLOW
    gocliselect.Cursor.SubMenuColor = goterm.CYAN
    gocliselect.Cursor.Suffix = " "

    var colorData ColorData
    err := json.Unmarshal([]byte(jsonData), &colorData)
    if err != nil {
        fmt.Println("Error unmarshaling JSON:", err)
        return
    }

    // fmt.Printf("Color: %s\n", colorData.Color)
    menu := createMenu(colorData.Color, colorData.Items, 0)
    if _, choice := menu.Display(menu); choice != nil {
      fmt.Printf("Choice: %s\n", choice.Text)
    }
}
