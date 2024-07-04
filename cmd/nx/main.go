package main

import (
	"xabbo.b7c.io/nx/cmd/nx/cmd"

	_ "xabbo.b7c.io/nx/cmd/nx/cmd/figure"
	_ "xabbo.b7c.io/nx/cmd/nx/cmd/figure/info"

	_ "xabbo.b7c.io/nx/cmd/nx/cmd/furni"
	_ "xabbo.b7c.io/nx/cmd/nx/cmd/furni/info"
	_ "xabbo.b7c.io/nx/cmd/nx/cmd/furni/search"

	_ "xabbo.b7c.io/nx/cmd/nx/cmd/get"
	_ "xabbo.b7c.io/nx/cmd/nx/cmd/get/furni"

	_ "xabbo.b7c.io/nx/cmd/nx/cmd/profile"

	_ "xabbo.b7c.io/nx/cmd/nx/cmd/visual"

	_ "xabbo.b7c.io/nx/cmd/nx/cmd/render"
	_ "xabbo.b7c.io/nx/cmd/nx/cmd/render/avatar"
	_ "xabbo.b7c.io/nx/cmd/nx/cmd/render/furni"

	_ "xabbo.b7c.io/nx/cmd/nx/cmd/texts"

	_ "xabbo.b7c.io/nx/cmd/nx/cmd/vars"
	
	_ "xabbo.b7c.io/nx/cmd/nx/cmd/extract"
)

func main() {
	cmd.Execute()
}
