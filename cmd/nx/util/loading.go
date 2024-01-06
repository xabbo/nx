package util

import (
	"github.com/b7c/nx"

	"cli/spinner"
)

func LoadGamedata(mgr *nx.GamedataManager, message string, types ...nx.GamedataType) error {
	return spinner.DoErr(message, func() error {
		return mgr.Load(types...)
	})
}

func LoadTexts(mgr *nx.GamedataManager) error {
	return LoadGamedata(mgr, "Loading external texts...", nx.GamedataTexts)
}

func LoadFurni(mgr *nx.GamedataManager) error {
	return LoadGamedata(mgr, "Loading furni data...", nx.GamedataFurni)
}

func LoadFigure(mgr *nx.GamedataManager) error {
	return LoadGamedata(mgr, "Loading figure data...", nx.GamedataFigure)
}
