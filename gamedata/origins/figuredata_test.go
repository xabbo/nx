package origins

import (
	"testing"
)

func TestUnmarshalFiguredata(t *testing.T) {
	figureDataStr := `[
		"M": [
			"hr": [
				[
					"s": 1,
					"p": ["sh": "2"],
					"c": ["AAAAAA","BBBBBB","CCCCCC","DDDDDD"]
				]
			]
		]
	]`
	_, err := ParseFigureData([]byte(figureDataStr))
	if err != nil {
		t.Fatal(err)
	}
}
