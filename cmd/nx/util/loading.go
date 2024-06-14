package util

import (
	gd "xabbo.b7c.io/nx/gamedata"

	"xabbo.b7c.io/nx/cmd/nx/spinner"
)

func LoadGameData(mgr gd.Manager, message string, types ...gd.Type) error {
	return spinner.DoErr(message, func() error {
		return mgr.Load(types...)
	})
}

func LoadTexts(mgr gd.Manager) error {
	return LoadGameData(mgr, "Loading external texts...", gd.GameDataTexts)
}

func LoadFurni(mgr gd.Manager) error {
	return LoadGameData(mgr, "Loading furni data...", gd.GameDataFurni)
}

func LoadFigure(mgr gd.Manager) error {
	return LoadGameData(mgr, "Loading figure data...", gd.GameDataFigure)
}
