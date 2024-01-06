package nx

import (
	"bufio"
	"bytes"
	"fmt"
	"path"
	"strings"
)

type ExternalVariables map[string]string
type ExternalTexts map[string]string

func (texts *ExternalTexts) UnmarshalBytes(data []byte) (err error) {
	*texts = ExternalTexts(readKeyValueMap(data))
	return
}

func (vars *ExternalVariables) UnmarshalBytes(data []byte) (err error) {
	*vars = ExternalVariables(readKeyValueMap(data))
	return
}

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
