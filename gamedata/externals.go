package gamedata

import (
	"bufio"
	"bytes"
	"fmt"
	"path"
	"strings"
)

// ExternalVariables defines dynamic variables loaded by the client.
type ExternalVariables map[string]string

// ExternalTexts defines dynamic strings loaded by the client.
type ExternalTexts map[string]string

// Unmarshals a text file as raw bytes into an ExternalTexts.
func (texts *ExternalTexts) UnmarshalBytes(data []byte) (err error) {
	*texts = ExternalTexts(readKeyValueMap(data))
	return
}

// Unmarshals a text file as raw bytes into an ExternalVariables.
func (vars *ExternalVariables) UnmarshalBytes(data []byte) (err error) {
	*vars = ExternalVariables(readKeyValueMap(data))
	return
}

// Gets the client version from the external variables,
// or returns an error if the key is not found.
func (vars *ExternalVariables) ClientVersion() (version string, err error) {
	if clientUrl, ok := (*vars)["flash.client.url"]; ok {
		version = path.Base(clientUrl)
	} else {
		err = fmt.Errorf("key not found")
	}
	return
}

func readKeyValueMap(data []byte) map[string]string {
	m := map[string]string{}

	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		key, value := keyValueSplit(sc.Text())
		m[key] = value
	}

	return m
}

func keyValueSplit(s string) (key, value string) {
	split := strings.SplitN(s, "=", 2)
	switch len(split) {
	case 2:
		value = strings.TrimSpace(split[1])
		fallthrough
	case 1:
		key = strings.TrimSpace(split[0])
	}
	return
}
