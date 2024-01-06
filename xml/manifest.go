package xml

type Manifest struct {
	Libraries []Library `xml:"library"`
}

type Library struct {
	Name    string  `xml:"name,attr"`
	Version string  `xml:"version,attr"`
	Assets  []Asset `xml:"assets>asset"`
}
