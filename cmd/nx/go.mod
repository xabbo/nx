module github.com/xabbo/nx/cli/nx

go 1.21.1

require (
	cli v0.0.0-00010101000000-000000000000
	github.com/disintegration/imaging v1.6.2
	github.com/phrozen/blend v0.0.0-20210220204729-f26b6cf7a28e
	github.com/spf13/cobra v1.8.0
	github.com/theckman/yacspin v0.13.12
	github.com/b7c/swfx v0.0.0-20231108034426-db9ec2dafc55
	golang.org/x/net v0.18.0
)

require (
	github.com/dustin/go-humanize v1.0.1
	github.com/fatih/color v1.16.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jedib0t/go-pretty/v6 v6.5.0
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/image v0.14.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/term v0.16.0
	golang.org/x/text v0.14.0 // indirect
)

replace github.com/b7c/nx => ../../

replace cli => ./
