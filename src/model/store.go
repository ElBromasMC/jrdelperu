package model

import (
	"errors"
	"strings"
)

type Type string

const (
	VidrioType   Type = "VIDRIO"
	AluminioType Type = "ALUMINIO"
	UPVCType     Type = "UPVC"
)

type Category struct {
	Id              int
	Type            Type
	Name            string
	Description     string
	LongDescription string
	Slug            string
	Img             Image
}

type Item struct {
	Id              int
	Category        Category
	Name            string
	Description     string
	LongDescription string
	Slug            string
	Img             Image
}

type CategoryFeature struct {
	Id          int
	Category    Category
	Name        string
	Description string
}

type Product struct {
	Id   int
	Item Item
	Name string
	Slug string
}

func (cat Category) Normalize() (Category, error) {
	// Trim name
	cat.Name = strings.TrimSpace(cat.Name)

	if len(cat.Name) == 0 {
		return Category{}, errors.New("invalid name")
	}

	return cat, nil
}

func (i Item) Normalize() (Item, error) {
	// Trim name
	i.Name = strings.TrimSpace(i.Name)

	if len(i.Name) == 0 {
		return Item{}, errors.New("invalid name")
	}

	return i, nil
}

func (p Product) Normalize() (Product, error) {
	// Trim name
	p.Name = strings.TrimSpace(p.Name)

	if len(p.Name) == 0 {
		return Product{}, errors.New("invalid name")
	}

	return p, nil
}
