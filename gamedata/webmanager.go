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
	"strconv"
	"strings"
	"time"

	"b7c.io/swfx"
	"golang.org/x/sync/errgroup"
	j "xabbo.b7c.io/nx/raw/json"
	"xabbo.b7c.io/nx/res"
)

type webGameDataManager struct {
	client   *http.Client
	host     string
	cacheDir string
	hashes   map[Type]string

	figure        *FigureData
	figureMap     *FigureMap
	avatarActions AvatarActions
	furni         FurniData
	products      ProductData
	texts         ExternalTexts
	variables     ExternalVariables
	assets        res.LibraryManager

	currentHashes *j.GameDataHashes
	lastFetched   map[Type]time.Time
}

func (mgr *webGameDataManager) Library(name string) res.AssetLibrary {
	return mgr.assets.Library(name)
}

func (mgr *webGameDataManager) Libraries() []string {
	return mgr.assets.Libraries()
}

func (mgr *webGameDataManager) LibraryExists(name string) bool {
	return mgr.assets.LibraryExists(name)
}

func (mgr *webGameDataManager) AddLibrary(lib res.AssetLibrary) bool {
	return mgr.assets.AddLibrary(lib)
}

type bytesUnmarshaler interface {
	UnmarshalBytes(data []byte) error
}

// Creates a new web-based game data manager.
// The provided manager fetches assets from the web and caches assets to disk.
// The cache directory is located under `xabbo/nx` within the user's cache directory.
func NewManager(host string) Manager {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = ".cache"
	}
	cacheDir = filepath.Join(cacheDir, "xabbo", "nx")

	return &webGameDataManager{
		client:   &http.Client{},
		host:     host,
		hashes:   make(map[Type]string),
		cacheDir: cacheDir,
		assets:   res.NewManager(),
	}
}

func (mgr *webGameDataManager) Figure() *FigureData {
	return mgr.figure
}

func (mgr *webGameDataManager) FigureMap() *FigureMap {
	return mgr.figureMap
}

func (mgr *webGameDataManager) AvatarActions() AvatarActions {
	return mgr.avatarActions
}

func (mgr *webGameDataManager) Furni() FurniData {
	return mgr.furni
}

func (mgr *webGameDataManager) Products() ProductData {
	return mgr.products
}

func (mgr *webGameDataManager) Texts() ExternalTexts {
	return mgr.texts
}

func (mgr *webGameDataManager) Variables() ExternalVariables {
	return mgr.variables
}

func (mgr *webGameDataManager) Loaded(types ...Type) bool {
	for _, t := range types {
		switch t {
		case GameDataFurni:
			if mgr.furni == nil {
				return false
			}
		case GameDataFigure:
			if mgr.figure == nil {
				return false
			}
		case GameDataProduct:
			if mgr.products == nil {
				return false
			}
		case GameDataTexts:
			if mgr.texts == nil {
				return false
			}
		case GameDataVariables:
			if mgr.variables == nil {
				return false
			}
		case GameDataFigureMap:
			if mgr.figureMap == nil {
				return false
			}
		default:
			panic(fmt.Errorf("unknown game data type %q", t))
		}
	}
	return true
}

func (mgr *webGameDataManager) Load(types ...Type) (err error) {
	err = os.MkdirAll(mgr.cacheDir, 0755)
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
		gameDataType := Type(hash.Name)

		if currentHash, ok := mgr.hashes[gameDataType]; ok && currentHash == hash.Hash {
			// Game data already loaded with the same hash, continue
			continue
		}

		if len(types) > 0 && !slices.Contains(types, gameDataType) {
			continue
		}

		if t, ok := hashTypeMap[gameDataType]; ok {
			g.Go(func() error {
				ptr := reflect.New(t)
				gd := ptr.Interface().(bytesUnmarshaler)

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
					mgr.furni = *v
				case *FigureData:
					mgr.figure = v
				case *ProductData:
					mgr.products = *v
				case *ExternalTexts:
					mgr.texts = *v
				case *ExternalVariables:
					mgr.variables = *v
				default:
					return fmt.Errorf("unknown game data type: %T", v)
				}

				mgr.hashes[gameDataType] = hash.Hash
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
	if len(types) == 0 || slices.Contains(types, GameDataFigureMap) {
		clientUrl, exist := mgr.variables[keyFlashClientUrl]
		if !exist {
			err = fmt.Errorf("unable to load figure map - failed to retrieve %s from external variables",
				keyFlashClientUrl)
			return
		}
		version := path.Base(clientUrl)
		filePath := filepath.Join(mgr.cacheDir, mgr.host, string(GameDataFigureMap), version)

		var data []byte
		data, err = mgr.fetchOrGetCached(filePath, clientUrl+"figuremap.xml", 0)
		if err != nil {
			return
		}

		var figureMap FigureMap
		err = figureMap.UnmarshalBytes(data)
		if err != nil {
			return
		}
		mgr.figureMap = &figureMap
	}

	if len(types) == 0 || slices.Contains(types, GameDataAvatar) {
		clientUrl, exist := mgr.variables[keyFlashClientUrl]
		if !exist {
			err = fmt.Errorf("unable to load avatar actions - failed to retrieve %s from external variables", keyFlashClientUrl)
			return
		}
		version := path.Base(clientUrl)
		filePath := filepath.Join(mgr.cacheDir, mgr.host, string(GameDataAvatar), version)

		var data []byte
		data, err = mgr.fetchOrGetCached(filePath, clientUrl+habboAvatarActionsFilename, 0)
		if err != nil {
			return
		}

		var avatarActions AvatarActions
		err = avatarActions.UnmarshalBytes(data)
		if err != nil {
			return
		}
		mgr.avatarActions = avatarActions
	}

	return
}

func (mgr *webGameDataManager) LoadFurni(libraries ...string) (err error) {
	if mgr.variables == nil {
		err = fmt.Errorf("variables not loaded")
		return
	}
	if mgr.furni == nil {
		err = fmt.Errorf("furni data not loaded")
		return
	}

	downloadUrl, ok := mgr.variables["dynamic.download.url"]
	if !ok {
		err = fmt.Errorf("failed to find dynamic download url")
		return
	}

	furniCacheDir := filepath.Join(mgr.cacheDir, "swf", "furni")
	for _, identifier := range libraries {
		fi, ok := mgr.furni[identifier]
		if !ok {
			err = fmt.Errorf("failed to find furni info for %q", identifier)
			return
		}

		libraryName := strings.Split(identifier, "*")[0]

		if mgr.assets.LibraryExists(libraryName) {
			continue
		}

		fpath := filepath.Join(furniCacheDir, libraryName+".swf")
		url := downloadUrl + strconv.Itoa(fi.Revision) + "/" + libraryName + ".swf"
		var data []byte
		data, err = mgr.fetchOrGetCached(fpath, url, 0)
		if err != nil {
			return
		}

		var swf *swfx.Swf
		swf, err = swfx.ReadSwf(bytes.NewReader(data))
		if err != nil {
			return
		}

		var lib res.AssetLibrary
		lib, err = res.LoadFurniLibrarySwf(swf)
		if err != nil {
			return
		}

		mgr.assets.AddLibrary(lib)
	}

	return
}

func (mgr *webGameDataManager) LoadFigureParts(libraries ...string) (err error) {
	if mgr.variables == nil {
		err = fmt.Errorf("variables not loaded")
		return
	}

	clientUrl, ok := mgr.variables[keyFlashClientUrl]
	if !ok {
		err = fmt.Errorf("failed to find client url in external variables")
		return
	}

	for _, libraryName := range libraries {
		if mgr.assets.LibraryExists(libraryName) {
			continue
		}

		filePath := filepath.Join(mgr.cacheDir, "swf", "figure", libraryName+".swf")
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

		var lib res.AssetLibrary
		lib, err = res.LoadFigureLibrarySwf(swf)
		if err != nil {
			return
		}
		mgr.assets.AddLibrary(lib)
	}

	return
}

func (mgr *webGameDataManager) GetHashes() (hashes *j.GameDataHashes, err error) {
	if mgr.currentHashes != nil {
		if lastFetched, ok := mgr.lastFetched[GameDataHashes]; ok {
			if time.Since(lastFetched).Hours() < 4 {
				hashes = mgr.currentHashes
				return
			}
		}
	}
	data, err := mgr.fetchOrGetCached(
		filepath.Join(mgr.cacheDir, mgr.host, "hashes.json"),
		"https://"+mgr.host+"/gamedata/hashes2",
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

func (mgr *webGameDataManager) DownloadHash(hash j.GameDataHash) (data []byte, err error) {
	data, err = mgr.fetchOrGetCached(
		filepath.Join(mgr.cacheDir, mgr.host, hash.Name, hash.Hash),
		hash.Url+"/"+hash.Hash,
		time.Hour*24*365,
	)
	return
}

func (mgr *webGameDataManager) fetchOrGetCached(filePath string, url string, refetchThreshold time.Duration) (data []byte, err error) {
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

		res, err = mgr.client.Get(url)
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
