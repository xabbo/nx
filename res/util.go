package res

import (
	"bytes"
	"encoding/xml"

	"golang.org/x/net/html/charset"

	"b7c.io/swfx"
)

func getBinaryTag(swf *swfx.Swf, symbol string) *swfx.DefineBinaryData {
	if id, ok := swf.Symbols[symbol]; ok {
		if tag, ok := swf.Characters[id].(*swfx.DefineBinaryData); ok {
			return tag
		}
	}
	return nil
}

func getImageTag(swf *swfx.Swf, symbol string) swfx.ImageTag {
	if id, ok := swf.Symbols[symbol]; ok {
		if tag, ok := swf.Characters[id].(swfx.ImageTag); ok {
			return tag
		}
	}
	return nil
}

func decodeXml(b []byte, v any) error {
	reader := bytes.NewReader(b)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder.Decode(v)
}
