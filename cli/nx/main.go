package main

import (
	"github.com/xabbo/nx/cli/nx/cmd"
	_ "github.com/xabbo/nx/cli/nx/cmd/figure"
	_ "github.com/xabbo/nx/cli/nx/cmd/furni"
	_ "github.com/xabbo/nx/cli/nx/cmd/get"
	_ "github.com/xabbo/nx/cli/nx/cmd/profile"
	_ "github.com/xabbo/nx/cli/nx/cmd/render"
	_ "github.com/xabbo/nx/cli/nx/cmd/texts"
	_ "github.com/xabbo/nx/cli/nx/cmd/vars"
)

func main() {
	cmd.Execute()
}
