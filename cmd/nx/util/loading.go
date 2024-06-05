package util

import (
	gd "xabbo.b7c.io/nx/gamedata"

	"xabbo.b7c.io/nx/cmd/nx/spinner"
)

func LoadGamedata(mgr *gd.GamedataManager, message string, types ...gd.GamedataType) error {
	return spinner.DoErr(message, func() error {
		return mgr.Load(types...)
	})
}

func LoadTexts(mgr *gd.GamedataManager) error {
	return LoadGamedata(mgr, "Loading external texts...", gd.GamedataTexts)
}

func LoadFurni(mgr *gd.GamedataManager) error {
	return LoadGamedata(mgr, "Loading furni data...", gd.GamedataFurni)
}

func LoadFigure(mgr *gd.GamedataManager) error {
	return LoadGamedata(mgr, "Loading figure data...", gd.GamedataFigure)
}
