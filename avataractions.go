package nx

import (
	"encoding/xml"

	x "github.com/b7c/nx/xml"
)

type AvatarActions map[string]AvatarActionInfo

type AvatarActionInfo struct {
	Id string
}

func (actions *AvatarActions) UnmarshalBytes(data []byte) (err error) {
	var xactions x.ActionContainer
	err = xml.Unmarshal(data, &xactions)
	if err != nil {
		return
	}

	*actions = AvatarActions{}
	for _, xaction := range xactions.Actions {
		(*actions)[xaction.Id] = AvatarActionInfo{
			Id: xaction.Id,
		}
	}

	return
}
