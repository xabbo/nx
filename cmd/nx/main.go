package main

import (
	"xabbo.b7c.io/nx/cmd/nx/cmd"
	_ "xabbo.b7c.io/nx/cmd/nx/cmd/figure"
	_ "xabbo.b7c.io/nx/cmd/nx/cmd/furni"
	_ "xabbo.b7c.io/nx/cmd/nx/cmd/get"
	_ "xabbo.b7c.io/nx/cmd/nx/cmd/profile"
	_ "xabbo.b7c.io/nx/cmd/nx/cmd/render"
	_ "xabbo.b7c.io/nx/cmd/nx/cmd/texts"
	_ "xabbo.b7c.io/nx/cmd/nx/cmd/vars"
)

func main() {
	cmd.Execute()
}
