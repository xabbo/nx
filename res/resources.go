package res

import (
	"fmt"
	"image"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

type assetManager struct {
	libs map[string]AssetLibrary
}

func NewManager() LibraryManager {
	return &assetManager{
		libs: map[string]AssetLibrary{},
	}
}

func (mgr *assetManager) Library(name string) AssetLibrary {
	return mgr.libs[name]
}

func (mgr *assetManager) Libraries() []string {
	return maps.Keys(mgr.libs)
}

func (mgr *assetManager) LibraryExists(name string) bool {
	_, exists := mgr.libs[name]
	return exists
}

func (mgr *assetManager) AddLibrary(lib AssetLibrary) bool {
	if _, exists := mgr.libs[lib.Name()]; exists {
		// err = fmt.Errorf("library already loaded: %q", library.Name())
		return false
	} else {
		mgr.libs[lib.Name()] = lib
		return true
	}
}

func parsePoint(s string) (pt image.Point, err error) {
	split := strings.Split(s, ",")
	if len(split) != 2 {
		err = fmt.Errorf("invalid point")
		return
	}
	var x, y int
	x, err = strconv.Atoi(split[0])
	if err != nil {
		err = fmt.Errorf("invalid point")
		return
	}
	y, err = strconv.Atoi(split[1])
	if err != nil {
		err = fmt.Errorf("invalid point")
		return
	}
	pt = image.Pt(x, y)
	return
}
