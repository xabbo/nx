package main

import (
	"cli/cmd"
	_ "cli/cmd/figure"
	_ "cli/cmd/furni"
	_ "cli/cmd/get"
	_ "cli/cmd/profile"
	_ "cli/cmd/render"
	_ "cli/cmd/texts"
	_ "cli/cmd/vars"
)

func main() {
	cmd.Execute()
}
