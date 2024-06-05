package profile

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	root "xabbo.b7c.io/nx/cmd/nx/cmd"
	"xabbo.b7c.io/nx/cmd/nx/spinner"
	"xabbo.b7c.io/nx/cmd/nx/util"

	"xabbo.b7c.io/nx"
	"xabbo.b7c.io/nx/web"
)

var (
	outputJson bool
)

var profileCmd = &cobra.Command{
	Use:     "profile [name|unique-id]",
	Aliases: []string{"user"},
	Short:   "Gets a user's profile information",
	RunE:    run,
}

func init() {
	root.Cmd.AddCommand(profileCmd)

	profileCmd.Flags().BoolVar(&outputJson, "json", false, "Output raw JSON data")
}

func run(cmd *cobra.Command, args []string) (err error) {
	// TODO: check CalledAs() user or profile
	if len(args) == 0 {
		return fmt.Errorf("no name or unique id provided")
	}
	cmd.SilenceUsage = true

	api := nx.NewApiClient(root.Host)

	userName := args[0]

	if outputJson {
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
