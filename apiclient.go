package nx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime/debug"

	"github.com/b7c/nx/web"
)

var (
	ErrUserNotFound = errors.New("the user was not found")
	ErrUserBanned   = errors.New("the user is banned")
	ErrMaintenance  = errors.New("the server is under maintenance")
)

type ApiClient struct {
	Http *http.Client
	Host string
	// Issue an extra request to determine if a user who is not found previously existed,
	// indicating that they have been permanently banned.
	// If so, the error returned will be ErrUserBanned.
	CheckBan bool
	Agent    string // The user agent to use.
}

type apiClientRt struct {
	client    *ApiClient
	transport http.RoundTripper
}

func (rt *apiClientRt) RoundTrip(r *http.Request) (*http.Response, error) {
	if rt.client.Agent != "" {
		r.Header.Set("User-Agent", rt.client.Agent)
	}
	return rt.transport.RoundTrip(r)
}

type httpError struct {
	Response *http.Response
}

func (e httpError) Error() string {
	return "server responded " + e.Response.Status
}

func NewApiClient(host string) *ApiClient {
	var agent string
	bi, ok := debug.ReadBuildInfo()
	if ok {
		version := bi.Main.Version
		if version == "(devel)" {
			version = "dev"
		}
		agent = "nx/" + version
	}
	client := &ApiClient{
		Http:     &http.Client{},
		Host:     host,
		Agent:    agent,
		CheckBan: true,
	}
	client.Http.Transport = &apiClientRt{
		client:    client,
		transport: http.DefaultTransport,
	}
	return client
}

func (c *ApiClient) url(path string, query url.Values) *url.URL {
	return &url.URL{
		Scheme:   "https",
		Host:     c.Host,
		Path:     path,
		RawQuery: query.Encode(),
	}
}

func (c *ApiClient) urlUserById(id HabboId) *url.URL {
	return c.url("/api/public/users/"+id.String(), nil)
}

func (c *ApiClient) urlUserByName(name string) *url.URL {
	return c.url("/api/public/users", url.Values{"name": {name}})
}

func (c *ApiClient) urlProfileById(id HabboId) *url.URL {
	return c.url("/api/public/users/"+id.String()+"/profile", nil)
}

func (c *ApiClient) urlAvatarImage(values url.Values) *url.URL {
	return c.url("/habbo-imaging/avatarimage", values)
}

func (c *ApiClient) urlAvatarImageName(name string) *url.URL {
	return c.urlAvatarImage(url.Values{"user": {name}})
}

func (c *ApiClient) doRequest(req *http.Request) (res *http.Response, data []byte, err error) {
	res, err = c.Http.Do(req)
	if err == nil {
		defer res.Body.Close()
		data, err = io.ReadAll(res.Body)
		if err == nil {
			switch res.StatusCode {
			case http.StatusServiceUnavailable:
				err = ErrMaintenance
			}
		}
	}
	return
}

func (c *ApiClient) getUrl(u *url.URL) (res *http.Response, data []byte, err error) {
	return c.doRequest(&http.Request{URL: u})
}

func (c *ApiClient) getRawUserUrl(u *url.URL) (data []byte, err error) {
	res, data, err := c.getUrl(u)
	if err != nil {
		return
	}
	switch res.StatusCode {
	case http.StatusOK:
		// No error
	case http.StatusNotFound:
		err = ErrUserNotFound
	default:
		err = httpError{Response: res}
	}
	return
}

// Gets the raw response of the specified user's info.
func (c *ApiClient) GetRawUser(name string) (data []byte, err error) {
	res, data, err := c.getUrl(c.urlUserByName(name))
	if err != nil {
		return
	}
	switch res.StatusCode {
	case http.StatusOK:
		// No error
	case http.StatusNotFound:
		err = ErrUserNotFound
	default:
		err = httpError{Response: res}
	}
	return
}

func (c *ApiClient) getUserUrl(u *url.URL) (user web.User, err error) {
	res, data, err := c.getUrl(u)
	if err != nil {
		return
	}

	switch res.StatusCode {
	case http.StatusOK:
		err = json.Unmarshal(data, &user)
	case http.StatusNotFound:
		err = ErrUserNotFound
	case http.StatusServiceUnavailable:
		err = ErrMaintenance
	default:
		err = errors.New("server responded " + res.Status)
	}

	return
}

func (c *ApiClient) getProfileUrl(u *url.URL) (profile web.Profile, err error) {
	res, data, err := c.getUrl(u)
	if err != nil {
		return
	}

	switch res.StatusCode {
	case http.StatusOK:
		err = json.Unmarshal(data, &profile)
	case http.StatusNotFound:
		err = ErrUserNotFound
	case http.StatusServiceUnavailable:
		err = ErrMaintenance
	default:
		err = fmt.Errorf("server responded %s", res.Status)
	}

	return
}

// Checks if a user exists by sending a HEAD request to the avatar imaging API.
// This can determine whether a user exists even if their profile is not found due to being permanently banned.
func (c *ApiClient) GetUserExists(name string) (exists bool, err error) {
	res, _, err := c.doRequest(&http.Request{
		Method: http.MethodHead,
		URL:    c.urlAvatarImageName(name),
	})
	if err != nil {
		return
	}

	switch res.StatusCode {
	case http.StatusOK:
		exists = true
	case http.StatusNotFound:
		exists = false
	default:
		err = httpError{res}
	}

	return
}

// Gets a user's information by their name.
// If the user was not found, and CheckBan is true,
// an extra request will be issued to determine
// whether the user exists and was not found due to being permanently banned.
func (c *ApiClient) GetUserByName(name string) (user web.User, err error) {
	user, err = c.getUserUrl(c.urlUserByName(name))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) && c.CheckBan {
			var exists bool
			exists, err = c.GetUserExists(name)
			if err != nil {
				return
			}
			if exists {
				err = ErrUserBanned
			}
		}
	}
	return
}

// Gets a user's information by their unique HabboId.
func (c *ApiClient) GetUser(uid HabboId) (user web.User, err error) {
	if uid.Kind != HabboIdKindUser {
		err = fmt.Errorf("non-user HabboId specified")
		return
	}

	return c.getUserUrl(c.urlUserById(uid))
}

func (c *ApiClient) GetProfile(uid HabboId) (profile web.Profile, err error) {
	if uid.Kind != HabboIdKindUser {
		err = fmt.Errorf("non-user HabboId specified")
		return
	}

	return c.getProfileUrl(c.urlProfileById(uid))
}
