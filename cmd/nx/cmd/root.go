package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var hotels = map[string]string{
	"us":  "www.habbo.com",
	"ous": "origins.habbo.com",
	"es":  "www.habbo.es",
	"fi":  "www.habbo.fi",
	"it":  "www.habbo.it",
	"nl":  "www.habbo.nl",
	"de":  "www.habbo.de",
	"fr":  "www.habbo.fr",
	"br":  "www.habbo.com.br",
	"tr":  "www.habbo.com.tr",
	"s2":  "sandbox.habbo.com",
}

var Cmd = &cobra.Command{
	Use:               "nx",
	Short:             "A command-line toolkit for Habbo Hotel.",
	PersistentPreRunE: preRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		if opts.showHotels {
			for id, host := range hotels {
				fmt.Printf("%s: %s\n", id, host)
			}
			return nil
		}

		return cmd.Usage()
	},
}

var opts struct {
	showHotels bool
}

var (
	Hotel string
	Host  string
)

func init() {
	Cmd.CompletionOptions.DisableDefaultCmd = true

	defaultHotel := "us"
	if envHotel, exist := os.LookupEnv("HOTEL"); exist {
		if _, ok := hotels[envHotel]; ok {
			defaultHotel = envHotel
		}
	}

	pf := Cmd.PersistentFlags()
	pf.StringVar(&Hotel, "hotel", defaultHotel, "The hotel to fetch information from")

	f := Cmd.Flags()
	f.BoolVar(&opts.showHotels, "hotels", false, "Show a list of supported hotels")
}

func preRun(cmd *cobra.Command, args []string) error {
	if !cmd.Flags().Lookup("hotel").Changed {
		hotel, ok := os.LookupEnv("HOTEL")
		if ok {
			Hotel = hotel
		}
	}

	var ok bool
	Host, ok = hotels[Hotel]
	if !ok {
		return fmt.Errorf("unknown hotel: %q", Hotel)
	}
	return nil
}

func Execute() {
	Cmd.SetOut(os.Stdout)
	Cmd.SetErr(os.Stderr)
	err := Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
