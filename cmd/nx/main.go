package main

import (
	"xabbo.io/nx/cmd/nx/cmd"

	_ "xabbo.io/nx/cmd/nx/cmd/figure"
	_ "xabbo.io/nx/cmd/nx/cmd/figure/convert"
	_ "xabbo.io/nx/cmd/nx/cmd/figure/info"

	_ "xabbo.io/nx/cmd/nx/cmd/furni"
	_ "xabbo.io/nx/cmd/nx/cmd/furni/info"
	_ "xabbo.io/nx/cmd/nx/cmd/furni/search"

	_ "xabbo.io/nx/cmd/nx/cmd/get"
	_ "xabbo.io/nx/cmd/nx/cmd/get/furni"

	_ "xabbo.io/nx/cmd/nx/cmd/profile"

	_ "xabbo.io/nx/cmd/nx/cmd/visual"

	_ "xabbo.io/nx/cmd/nx/cmd/imager"
	_ "xabbo.io/nx/cmd/nx/cmd/imager/avatar"
	_ "xabbo.io/nx/cmd/nx/cmd/imager/furni"

	_ "xabbo.io/nx/cmd/nx/cmd/texts"

	_ "xabbo.io/nx/cmd/nx/cmd/vars"

	_ "xabbo.io/nx/cmd/nx/cmd/extract"
)

func main() {
	cmd.Execute()
}
