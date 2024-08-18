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

func TestFixFiguredata(t *testing.T) {
	var (
		beforeFix        = `["M":["hr":[["s":1,"p":["sh":"2"],"c":["AAAAAA","BBBBBB","CCCCCC","DDDDDD"]]]]]`
		expectedAfterFix = `{"M":{"hr":[{"s":1,"p":{"sh":"2"},"c":["AAAAAA","BBBBBB","CCCCCC","DDDDDD"]}]}}`
	)

	bytes := []byte(beforeFix)
	err := fixFigureData(bytes)
	if err != nil {
		t.Fatal(err)
	}
	actualAfterFix := string(bytes)

	if expectedAfterFix != actualAfterFix {
		t.Fatalf("Fixed figure data not as expected\n  Expected: %q\n    Actual: %q", expectedAfterFix, actualAfterFix)
	}
}
