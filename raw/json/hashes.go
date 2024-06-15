package json

type GameDataHashes struct {
	Hashes []GameDataHash `json:"hashes"`
}

type GameDataHash struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Hash string `json:"hash"`
}
