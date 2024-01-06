package json

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ProductDataContainer struct {
	ProductData ProductData `json:"productdata"`
}

type ProductData struct {
	Products []ProductInfo `json:"product"`
}

type ProductInfo struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (info *ProductInfo) UnmarshalJSON(d []byte) (err error) {
	type productInfo ProductInfo
	x := struct {
		productInfo
		Code any `json:"code"`
	}{productInfo: productInfo(*info)}

	err = json.Unmarshal(d, &x)
	if err != nil {
		return
	}

	*info = ProductInfo(x.productInfo)
	switch v := x.Code.(type) {
	case string:
		(*info).Code = v
	case float64:
		// ofc they have one random product code with a number
		// when every other one is a string.
		(*info).Code = strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Errorf("unknown type %T for ProductInfo.Code", x.Code)
	}

	return
}
