package res

type LibraryManager interface {
	Library(name string) AssetLibrary
	Libraries() []string
	LibraryExists(name string) bool
	Load(LibraryLoader) error
}

type LibraryLoader interface {
	Load() (AssetLibrary, error)
}

type AssetLibrary interface {
	Name() string
	Asset(name string) (Asset, error)
	Assets() []string
	AssetExists(name string) bool
}

type FurniLibrary interface {
	Index() *Index
	Manifest() *Manifest
}
