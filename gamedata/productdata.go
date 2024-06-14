package gamedata

import (
	"encoding/json"
	"fmt"

	j "xabbo.b7c.io/nx/json"
)

// ProductData maps product info by product code.
type ProductData map[string]ProductInfo

// ProductInfo defines a product code, name and description.
type ProductInfo struct {
	Code        string
	Name        string
	Description string
}

// Unmarshals a JSON document as raw bytes into a ProductData.
func (pd *ProductData) UnmarshalBytes(data []byte) (err error) {
	var jpd j.ProductDataContainer
	err = json.Unmarshal(data, &jpd)
	if err != nil {
		return
	}

	*pd = ProductData{}
	for _, p := range jpd.ProductData.Products {
		if _, exist := (*pd)[p.Code]; exist {
			return fmt.Errorf("duplicate product code: %q", p.Code)
		}
		(*pd)[p.Code] = ProductInfo{
			Code:        p.Code,
			Name:        p.Name,
			Description: p.Description,
		}
	}

	return
}
