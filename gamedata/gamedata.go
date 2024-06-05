package gamedata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"slices"
	"time"

	"golang.org/x/sync/errgroup"

	"b7c.io/swfx"

	j "xabbo.b7c.io/nx/json"
	"xabbo.b7c.io/nx/res"
)

type GamedataType string

const (
	GamedataHashes    GamedataType = "hashes"
	GamedataFurni     GamedataType = "furnidata"
	GamedataProduct   GamedataType = "productdata"
	GamedataVariables GamedataType = "external_variables"
	GamedataTexts     GamedataType = "external_texts"
	GamedataFigure    GamedataType = "figurepartlist"
	GamedataFigureMap GamedataType = "figuremap"
	GamedataAvatar    GamedataType = "HabboAvatarActions"

	keyFlashClientUrl          = "flash.client.url"
	habboAvatarActionsFilename = "HabboAvatarActions.xml"
)

// Map of hashed game data types.
var hashTypeMap = map[GamedataType]reflect.Type{
	GamedataFurni:     reflect.TypeOf((*FurniData)(nil)).Elem(),
	GamedataFigure:    reflect.TypeOf((*FigureData)(nil)).Elem(),
	GamedataProduct:   reflect.TypeOf((*ProductData)(nil)).Elem(),
	GamedataTexts:     reflect.TypeOf((*ExternalTexts)(nil)).Elem(),
	GamedataVariables: reflect.TypeOf((*ExternalVariables)(nil)).Elem(),
}

type Gamedata interface {
	UnmarshalBytes(data []byte) error
}

type GamedataManager struct {
	Client        *http.Client
	Host          string
	CacheDir      string
	Hashes        map[GamedataType]string
	Furni         FurniData
	Figure        FigureData
	Product       ProductData
	Texts         ExternalTexts
	Variables     ExternalVariables
	FigureMap     FigureMap
	AvatarActions AvatarActions
	Assets        res.LibraryManager

	currentHashes *j.GamedataHashes
	lastFetched   map[GamedataType]time.Time
}

func (mgr *GamedataManager) FurniLoaded() bool {
	return len(mgr.Furni) > 0
}

func (mgr *GamedataManager) FigureLoaded() bool {
	return len(mgr.Figure.Sets) > 0
}

func (mgr *GamedataManager) ProductsLoaded() bool {
	return len(mgr.Product) > 0
}

func (mgr *GamedataManager) TextsLoaded() bool {
	return len(mgr.Texts) > 0
}

func (mgr *GamedataManager) VariablesLoaded() bool {
	return len(mgr.Variables) > 0
}

func (mgr *GamedataManager) FigureMapLoaded() bool {
	return len(mgr.FigureMap.Libs) > 0
}

func (mgr *GamedataManager) Loaded(types ...GamedataType) bool {
	for _, t := range types {
		switch t {
		case GamedataFurni:
			if !mgr.FurniLoaded() {
				return false
			}
		case GamedataFigure:
			if !mgr.FigureLoaded() {
				return false
			}
		case GamedataProduct:
			if !mgr.ProductsLoaded() {
				return false
			}
		case GamedataTexts:
			if !mgr.TextsLoaded() {
				return false
			}
		case GamedataVariables:
			if !mgr.VariablesLoaded() {
				return false
			}
		case GamedataFigureMap:
			if !mgr.FigureMapLoaded() {
				return false
			}
		default:
			panic(fmt.Errorf("unknown game data type %q", t))
		}
	}
	return true
}

func NewGamedataManager(host string) *GamedataManager {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = ""
	}
	cacheDir = filepath.Join(cacheDir, "b7c", "nx")
	return &GamedataManager{
		Client:   &http.Client{},
		Host:     host,
		Hashes:   make(map[GamedataType]string),
		CacheDir: cacheDir,
		Assets:   res.NewManager(),
	}
}

// Loads the specified game data types.
// If no types are specified, all types are loaded.
func (mgr *GamedataManager) Load(types ...GamedataType) (err error) {
	err = os.MkdirAll(mgr.CacheDir, 0755)
	if err != nil {
		return
	}

	hashes, err := mgr.GetHashes()
	if err != nil {
		return
	}
	mgr.currentHashes = hashes

	g := &errgroup.Group{}

	for i := range hashes.Hashes {
		hash := hashes.Hashes[i]
		gameDataType := GamedataType(hash.Name)

		if currentHash, ok := mgr.Hashes[gameDataType]; ok && currentHash == hash.Hash {
			// Game data already loaded with the same hash, continue
			continue
		}

		if len(types) > 0 && !slices.Contains(types, gameDataType) {
			continue
		}

		if t, ok := hashTypeMap[gameDataType]; ok {
			g.Go(func() error {
				ptr := reflect.New(t)
				gd := ptr.Interface().(Gamedata)

				data, err := mgr.DownloadHash(hash)
				if err != nil {
					return err
				}

				err = gd.UnmarshalBytes(data)
				if err != nil {
					return err
				}

				switch v := gd.(type) {
				case *FurniData:
					mgr.Furni = *v
				case *FigureData:
					mgr.Figure = *v
				case *ProductData:
					mgr.Product = *v
				case *ExternalTexts:
					mgr.Texts = *v
				case *ExternalVariables:
					mgr.Variables = *v
				default:
					return fmt.Errorf("unknown game data type: %T", v)
				}

				mgr.Hashes[gameDataType] = hash.Hash
				return nil
			})
		} else {
			return fmt.Errorf("unknown runtime type for gamedata type: %s", gameDataType)
		}
	}

	err = g.Wait()
	if err != nil {
		return
	}

	// Load figure map & avatar actions, depends on external variables
	if len(types) == 0 || slices.Contains(types, GamedataFigureMap) {
		clientUrl, exist := mgr.Variables[keyFlashClientUrl]
		if !exist {
			err = fmt.Errorf("unable to load figure map - failed to retrieve %s from external variables",
				keyFlashClientUrl)
			return
		}
		version := path.Base(clientUrl)
		filePath := filepath.Join(mgr.CacheDir, mgr.Host, string(GamedataFigureMap), version)

		var data []byte
		data, err = mgr.fetchOrGetCached(filePath, clientUrl+"figuremap.xml", 0)
		if err != nil {
			return
		}

		err = mgr.FigureMap.UnmarshalBytes(data)
		if err != nil {
			return
		}
	}

	if len(types) == 0 || slices.Contains(types, GamedataAvatar) {
		clientUrl, exist := mgr.Variables[keyFlashClientUrl]
		if !exist {
			err = fmt.Errorf("unable to load avatar actions - failed to retrieve %s from external variables", keyFlashClientUrl)
			return
		}
		version := path.Base(clientUrl)
		filePath := filepath.Join(mgr.CacheDir, mgr.Host, string(GamedataAvatar), version)

		var data []byte
		data, err = mgr.fetchOrGetCached(filePath, clientUrl+habboAvatarActionsFilename, 0)
		if err != nil {
			return
		}

		err = mgr.AvatarActions.UnmarshalBytes(data)
		if err != nil {
			return
		}
	}

	return
}

func (mgr *GamedataManager) LoadFurni(libraries ...string) (err error) {
	err = fmt.Errorf("not implemented")
	return
}

func (mgr *GamedataManager) LoadFigureParts(libraries ...string) (err error) {
	if !mgr.VariablesLoaded() {
		err = fmt.Errorf("variables not loaded")
		return
	}

	clientUrl, ok := mgr.Variables[keyFlashClientUrl]
	if !ok {
		err = fmt.Errorf("failed to find client url in external variables")
		return
	}

	for _, libraryName := range libraries {
		if mgr.Assets.LibraryExists(libraryName) {
			continue
		}

		filePath := filepath.Join(mgr.CacheDir, "swf", "figure", libraryName+".swf")
		libraryUrl := clientUrl + libraryName + ".swf"

		var data []byte
		data, err = mgr.fetchOrGetCached(filePath, libraryUrl, 0)
		if err != nil {
			return
		}

		var swf *swfx.Swf
		swf, err = swfx.ReadSwf(bytes.NewReader(data))
		if err != nil {
			return
		}

		err = mgr.Assets.Load(res.NewSwfFigureLibraryLoader(swf))
		if err != nil {
			return
		}
	}

	return
}

func (mgr *GamedataManager) GetHashes() (hashes *j.GamedataHashes, err error) {
	if mgr.currentHashes != nil {
		if lastFetched, ok := mgr.lastFetched[GamedataHashes]; ok {
			if time.Since(lastFetched).Hours() < 4 {
				hashes = mgr.currentHashes
				return
			}
		}
	}
	data, err := mgr.fetchOrGetCached(
		filepath.Join(mgr.CacheDir, mgr.Host, "hashes.json"),
		"https://"+mgr.Host+"/gamedata/hashes2",
		time.Hour*4,
	)
	if err == nil {
		err = json.Unmarshal(data, &hashes)
		if err == nil {
			mgr.currentHashes = hashes
		}
	}
	return
}

func (mgr *GamedataManager) DownloadHash(hash j.GamedataHash) (data []byte, err error) {
	data, err = mgr.fetchOrGetCached(
		filepath.Join(mgr.CacheDir, mgr.Host, hash.Name, hash.Hash),
		hash.Url+"/"+hash.Hash,
		time.Hour*24*365,
	)
	return
}

func (mgr *GamedataManager) fetchOrGetCached(filePath string, url string, refetchThreshold time.Duration) (data []byte, err error) {
	dir := filepath.Dir(filePath)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	stats, err := f.Stat()
	if err != nil {
		return
	}

	size := stats.Size()
	if size > 0 && (refetchThreshold == 0 || time.Since(stats.ModTime()) <= refetchThreshold) {
		data = make([]byte, size)
		_, err = io.ReadFull(f, data)
		if err != nil {
			return
		}
	}

	if len(data) == 0 {
		err = f.Truncate(0)
		if err != nil {
			return
		}

		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			return
		}

		var res *http.Response

		res, err = mgr.Client.Get(url)
		if err != nil {
			return
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			err = fmt.Errorf("server responded %s", res.Status)
			return
		}

		buf := bytes.Buffer{}
		_, err = io.Copy(io.MultiWriter(f, &buf), res.Body)
		if err != nil {
			return
		}

		data = buf.Bytes()
	}

	return
}
