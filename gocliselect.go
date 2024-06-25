package gocliselect

import (
	"fmt"
	"github.com/buger/goterm"
	"github.com/pkg/term"
	"log"
	"strings"
)

type CursorConfig struct {
    Color int
    ItemPrompt string
    SubMenuPrompt string
    Suffix string
	IndentMultiplier int
}

var Cursor = CursorConfig{
    Color:  goterm.YELLOW,
    ItemPrompt: ">",
	SubMenuPrompt: ">",
    Suffix: "  ",
	IndentMultiplier: 4,
}

// Raw input keycodes
var up byte = 65
var down byte = 66
var right byte = 67
var left byte = 68
var escape byte = 27
var ctrl_c byte = 3
var enter byte = 13
var keys = map[byte]bool {
	up: true,
	down: true,
	right: true,
	left: true,
}

var linesOnInput int = 0

type Menu struct {
	Prompt  	string
	CursorPos 	int
	Level		int  // TODO ideally this is computed on the fly
	MenuItems 	[]*MenuItem
	ParentMenu  *Menu
}

type MenuItem struct {
	Text     string
	ID       string
	SubMenu  *Menu
}

func NewMenu(prompt string, level int) *Menu {
	return &Menu{
		Prompt: prompt,
		MenuItems: make([]*MenuItem, 0),
		Level: level,
	}
}

// AddItemMenu will add a new item menu 
func (m *Menu) AddItemMenu(option string, id string, subMenu *Menu) *Menu {
	subMenu.ParentMenu = m
	menuItem := &MenuItem{
		Text: option,
		ID: id,
		SubMenu: subMenu,
	}

	m.MenuItems = append(m.MenuItems, menuItem)
	return m
}

// AddItem will add a new item
func (m *Menu) AddItem(option string, id string) *Menu {
	menuItem := &MenuItem{
		Text: option,
		ID: id,
	}

	m.MenuItems = append(m.MenuItems, menuItem)
	return m
}

func (m *Menu) renderRecursivelyUp(root *Menu) {
	if m.ParentMenu != root {
		m.ParentMenu.renderRecursivelyUp(root)
	}
	for index, menuItem := range m.ParentMenu.MenuItems[:m.ParentMenu.CursorPos + 1] {
		menuItemText := menuItem.Text
		cursor := " " + Cursor.Suffix
		if index == m.ParentMenu.CursorPos {
			menuItemText = goterm.Color(menuItemText, Cursor.Color)
		}
		fmt.Printf("\r%s%s%s\n", strings.Repeat(" ", m.ParentMenu.Level * Cursor.IndentMultiplier), cursor, menuItemText)
		linesOnInput++
	}
}

func (m *Menu) renderRecursivelyDown(root *Menu) {
	for _, menuItem := range m.ParentMenu.MenuItems[m.ParentMenu.CursorPos + 1:] {
		menuItemText := menuItem.Text
		cursor := " " + Cursor.Suffix
		fmt.Printf("\r%s%s%s\n", strings.Repeat(" ", m.ParentMenu.Level * Cursor.IndentMultiplier), cursor, menuItemText)
		linesOnInput++
	}
	if m.ParentMenu != root {
		m.ParentMenu.renderRecursivelyDown(root)
	}
}

func (m *Menu) RenderMenu(root *Menu) {
	// move cursor up N lines
	if linesOnInput > 0 {
		fmt.Printf("\r\033[%dA", linesOnInput)
		linesOnInput = 0
	}
	// clear screen from cursor down
	fmt.Print("\033[J")
	if m != root {
		m.renderRecursivelyUp(root)
	}
	for index, menuItem := range m.MenuItems {
		menuItemText := menuItem.Text
		cursor := " " + Cursor.Suffix
		if index == m.CursorPos {
			if menuItem.SubMenu != nil {
				cursor = goterm.Color(Cursor.SubMenuPrompt + Cursor.Suffix, Cursor.Color)
			} else {
				cursor = goterm.Color(Cursor.ItemPrompt + Cursor.Suffix, Cursor.Color)
			}
			menuItemText = goterm.Color(menuItemText, Cursor.Color)
		}
		fmt.Printf("\r%s%s%s\n", strings.Repeat(" ", m.Level * Cursor.IndentMultiplier), cursor, menuItemText)
		linesOnInput++
	}
	if m != root {
		m.renderRecursivelyDown(root)
	}
}

// Display will display the given menu and awaits user selection
// It returns the users selected choice
func (m *Menu) Display(root *Menu) *MenuItem {
	defer func() {
		// Show cursor again.
		fmt.Printf("\033[?25h")
	}()

	m.RenderMenu(root)

	// Turn the terminal cursor off
	fmt.Printf("\033[?25l")

	for {
		switch keyCode := getInput(); keyCode {

		case escape, ctrl_c:
			return nil

		case left:
			if m.ParentMenu != nil {
				return m.ParentMenu.Display(root)
			}
			return nil

		case right:
			menuItem := m.MenuItems[m.CursorPos]
			if menuItem.SubMenu != nil {
				return menuItem.SubMenu.Display(root)
			}
			linesOnInput = 0
			return menuItem

		case up:
			m.CursorPos = (m.CursorPos + len(m.MenuItems) - 1) % len(m.MenuItems)
			m.RenderMenu(root)

		case down:
			m.CursorPos = (m.CursorPos + 1) % len(m.MenuItems)
			m.RenderMenu(root)
		}
	}
}

// getInput will read raw input from the terminal
// It returns the raw ASCII value inputted
func getInput() byte {
	t, _ := term.Open("/dev/tty")

	err := term.RawMode(t)
	if err != nil {
		log.Fatal(err)
	}

	var read int
	readBytes := make([]byte, 3)
	read, err = t.Read(readBytes)

	t.Restore()
	t.Close()

	// Arrow keys are prefixed with the ANSI escape code which take up the first two bytes.
	// The third byte is the key specific value we are looking for.
	// For example the left arrow key is '<esc>[A' while the right is '<esc>[C'
	// See: https://en.wikipedia.org/wiki/ANSI_escape_code
	if read == 3 {
		if _, ok := keys[readBytes[2]]; ok {
			return readBytes[2]
		}
	} else {
		return readBytes[0]
	}

	return 0
}