package res

type LibraryManager interface {
	Library(name string) AssetLibrary
	Libraries() []string
	LibraryExists(name string) bool
	LoadLibrary(LibraryLoader) error
}

type LibraryLoader interface {
	Load() (AssetLibrary, error)
}

type AssetLibrary interface {
	Name() string
	Asset(name string) (*Asset, error)
	Assets() []string
	AssetExists(name string) bool
}

type FurniLibraryLoader interface {
	Load() (FurniLibrary, error)
}

type FurniLibrary interface {
	AssetLibrary
	Index() *Index
	Manifest() *Manifest
	Logic() *Logic
	Visualizations() map[int]*Visualization
}
