package json

type GamedataHashes struct {
	Hashes []GamedataHash `json:"hashes"`
}

type GamedataHash struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Hash string `json:"hash"`
}
