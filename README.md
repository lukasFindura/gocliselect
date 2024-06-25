# Golang CLI Select
## Fork of https://github.com/Nexidian/gocliselect
In this fork I implemented submenu options so the items can be in a nested structure.

## Examples

[examples/simple/main.go]()

![](examples/simple/example.gif)

[examples/advanced/main.go]()

![](examples/advanced/example.gif)

## Known issues

#### Text flickering
It seems `\033[J` (clear from cursor down) is causing random text flickering.