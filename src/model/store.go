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
}

type CategoryFeature struct {
	Id          int
	Category    Category
	Name        string
	Description string
}
