package nx

import (
	"io"
	"os"
	"testing"
)

func TestFurniData(t *testing.T) {
	f, err := os.Open("testdata/gamedata/furnidata.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	fd := FurniData{}
	err = fd.UnmarshalBytes(data)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Loaded %d furni", len(fd))
}

func TestFigureData(t *testing.T) {
	f, err := os.Open("testdata/gamedata/figuredata.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	fd := FigureData{}
	err = fd.UnmarshalBytes(data)
	if err != nil {
		t.Fatal(err)
	}

	totalColors := 0
	for _, palette := range fd.Palettes {
		totalColors += len(palette)
	}
	totalSets := 0
	for _, setMap := range fd.Sets {
		totalSets += len(setMap)
	}
	t.Logf("Loaded %d palettes / %d colors, %d part set types / %d part sets",
		len(fd.Palettes), totalColors, len(fd.Sets), totalSets)
}

func TestProductData(t *testing.T) {
	f, err := os.Open("testdata/gamedata/productdata.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	pd := ProductData{}
	err = pd.UnmarshalBytes(data)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Loaded %d products", len(pd))
}

func TestFigureMap(t *testing.T) {
	f, err := os.Open("testdata/gamedata/figuremap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	fm := FigureMap{}
	err = fm.UnmarshalBytes(data)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Read %d figure part libraries, %d figure parts", len(fm.Libs), len(fm.Parts))
}
