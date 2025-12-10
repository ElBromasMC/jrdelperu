package model

type Type int

const (
	VidrioType Type = iota
	AluminioType
	UPVCType
)

type Category struct {
	Id              int
	Type            Type
	Slug            string
	Name            string
	Description     string
	LongDescription string
	Img             Image
}

type Item struct {
	Id              int
	Category        Category
	Slug            string
	Name            string
	Description     string
	LongDescription string
	Img             Image
	SecondaryImg    Image
}

type CategoryFeature struct {
	Id          int
	Category    Category
	Name        string
	Description string
}

// SiteDocumentInfo represents a site-wide document with its URL
type SiteDocumentInfo struct {
	Key         string
	DisplayName string
	URL         string
	HasFile     bool
}

// PDFInfo represents a PDF from a category or item
type PDFInfo struct {
	Name        string // Category or item name
	DisplayName string // PDF display name
	URL         string
}
