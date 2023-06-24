package endpointstructs

import "image"

// /data/<brand>/<date_scraped>/<product id>/<color>/ then indexed from 0 to n for each picture there is of that color

type UniqueProductIdentifier struct {
	Brand string
	Id    string
}

type UniqueProductExpanded struct {
	Brand          string
	Id             string
	UrlOnBrandSite string
	Price          string
	Description    string
	ClothingType   string // tee, shorts, sweats, etc.
	ColorAttrs     []ColorAttr
}

// dateScraped is an attribute of the Color not of the product because the same product
// can have multiple colors added to the site at different dates
// ColorAttr can also uniquely identify the filesystem location if we also have the brand and product Id
type ColorAttr struct {
	colorName   string
	dateScraped string
}

type ProductContainerForFrontend struct {
	Upi       UniqueProductIdentifier
	ColorAttr ColorAttr
	imageImgs []image.Image
}

type ProductContainerForWritingToDB struct {
	Upi       UniqueProductIdentifier
	ColorAttr ColorAttr
	imageURLs []string
}
