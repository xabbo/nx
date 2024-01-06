package nx

import (
	"encoding/json"
	"fmt"

	j "github.com/b7c/nx/json"
)

type ProductData map[string]ProductInfo

type ProductInfo struct {
	Code        string
	Name        string
	Description string
}

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
