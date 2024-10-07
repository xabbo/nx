package gamedata

import (
	"encoding/xml"

	x "xabbo.io/nx/raw/xml"
)

type AvatarActions map[string]*AvatarActionInfo

type AvatarActionInfo struct {
	Id string
}

func (actions *AvatarActions) UnmarshalBytes(data []byte) (err error) {
	var xActions struct {
		Actions []x.Action
	}
	err = xml.Unmarshal(data, &xActions)
	if err != nil {
		return
	}

	*actions = AvatarActions{}
	for i := range xActions.Actions {
		xAction := &xActions.Actions[i]
		(*actions)[xAction.Id] = &AvatarActionInfo{
			Id: xAction.Id,
		}
	}

	return
}
