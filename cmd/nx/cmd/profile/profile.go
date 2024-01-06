package profile

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"

	root "cli/cmd"
	"cli/spinner"
	"cli/util"

	"github.com/b7c/nx"
	"github.com/b7c/nx/web"
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

	body, err := fetchRawUser(args[0], !outputJson)
	if err != nil {
		return err
	}

	if body == nil {
		return fmt.Errorf("user not found")
	}

	if outputJson {
		fmt.Println(string(body))
		return nil
	}

	var user web.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return err
	}

	util.RenderUserInfo(user)

	return nil
}

func fetchRawUser(name string, spin bool) (data []byte, err error) {
	if spin {
		spinner.Message("Fetching profile...")
		spinner.Start()
		defer spinner.Stop()
	}

	q := url.Values{}
	q.Add("name", name)
	u := url.URL{
		Scheme:   "https",
		Host:     root.Host,
		Path:     "/api/public/users",
		RawQuery: q.Encode(),
	}

	res, err := http.Get(u.String())
	if err != nil {
		return
	}

	switch res.StatusCode {
	case http.StatusOK:
		data, err = io.ReadAll(res.Body)
	case http.StatusNotFound:
		return
	case http.StatusServiceUnavailable:
		err = fmt.Errorf("maintenance")
	default:
		err = fmt.Errorf("server responded %s", res.Status)
	}

	return
}
