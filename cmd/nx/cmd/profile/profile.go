package profile

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	_root "xabbo.b7c.io/nx/cmd/nx/cmd"
	"xabbo.b7c.io/nx/cmd/nx/spinner"
	"xabbo.b7c.io/nx/cmd/nx/util"

	"xabbo.b7c.io/nx"
	"xabbo.b7c.io/nx/web"
)

var Cmd = &cobra.Command{
	Use:     "profile [name|unique-id]",
	Aliases: []string{"user"},
	Short:   "Gets a user's profile information",
	RunE:    run,
}

var opts struct {
	outputJson bool
}

func init() {
	f := Cmd.Flags()
	f.BoolVar(&opts.outputJson, "json", false, "Output raw JSON data")

	_root.Cmd.AddCommand(Cmd)
}

func run(cmd *cobra.Command, args []string) (err error) {
	// TODO: check CalledAs() user or profile
	if len(args) == 0 {
		return fmt.Errorf("no name or unique id provided")
	}
	cmd.SilenceUsage = true

	api := nx.NewApiClient(_root.Host)

	userName := args[0]

	if opts.outputJson {
		var data []byte
		err = spinner.DoErr("Loading user...", func() (err error) {
			data, err = api.GetRawUser(userName)
			return
		})
		if err != nil {
			return
		}
		os.Stdout.Write(data)
		fmt.Println()
		return
	}

	var user web.User
	err = spinner.DoErr("Loading user...", func() (err error) {
		user, err = api.GetUserByName(userName)
		return err
	})
	if err != nil {
		return
	}

	util.RenderUserInfo(user)

	return nil
}
