package ghligh

import (
	"github.com/orlandofra/ghligh/go-poppler"

	"encoding/json"
	"log"
)

type AnnotJSON struct {
	Type		poppler.AnnotType  `json:"type"`
	Index		int        `json:"index"`
	Date		string     `json:"date"`
	Rect		poppler.Rectangle  `json:"rect"`
	Color		poppler.Color      `json:"color"`
	Name		string     `json:"name"`
	Contents	string     `json:"contents"`
	Flags		poppler.AnnotFlag  `json:"flags"`
	Quads		[]poppler.Quad     `json:"quads"`
}

func NewAnnotJSON (a poppler.Annot) (AnnotJSON) {
	var aj AnnotJSON
	aj.Type = a.Type()
	aj.Index = a.Index()
	aj.Date = a.Date()
	aj.Rect = a.Rect()
	aj.Color = a.Color()
	aj.Name = a.Name()
	aj.Contents = a.Contents()
	aj.Flags = a.Flags()
	aj.Quads = a.Quads()

	return aj
}

func (aj *AnnotJSON) Export() string {
	jsonData, err := json.Marshal(aj)
	if err != nil {
		log.Printf("Error serializing AnnotJSON: %v", err)
		return ""
	}

	return string(jsonData)
}

