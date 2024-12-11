package ghligh

import (
	"github.com/orlandofra/ghligh/go-poppler"

	"encoding/json"
	"log"

	// debug
	//"fmt"
)

type GhlighDoc struct {
	doc		*poppler.Document
	Path		string		`json:"file"`
	AnnotsBuffer	json.RawMessage	`json:"annots"`
}

func Open(filename string) (*GhlighDoc, error) {
	var err error

	g := &GhlighDoc{}

	g.doc, err = poppler.Open(filename)
	if err != nil {
		log.Printf("error opening pdf: %v", err)
		return nil, err
	}
	g.Path = filename

	return g, nil
}

func (d *GhlighDoc) Close() {
	if d.doc != nil {
		d.doc.Close()
	}
}

func (d *GhlighDoc) getAnnots() json.RawMessage{
	/* TODO Refactor */
	n := d.doc.GetNPages()
	var ajs []AnnotJSON
	ajsMap := make(map[int][]AnnotJSON)


	for i := 0; i < n; i++ {

		page := d.doc.GetPage(i)

		annots := page.GetAnnots()
		for _, annot := range annots {
			if annot.Type() == poppler.AnnotHighlight {
				aj := NewAnnotJSON(*annot)
				ajs = append(ajs, aj)
			}
		}

		page.Close()

		if len(ajs) > 0 {
			ajsMap[i] = ajs
		}


	}

	jsonData, err := json.Marshal(ajsMap)
	if err != nil {
		log.Fatalf("Error serializing annotations to JSON: %v", err)
		return nil
	}

	return jsonData
}

func (d *GhlighDoc) Export() string {
	//fmt.Printf("%v", d.getAnnots())

	d.AnnotsBuffer = d.getAnnots()
	defer func() {
		/* mutex maybe??? */
		d.AnnotsBuffer = nil
	}()

	jsonData, err := json.Marshal(d)
	if err != nil {
		return ""
		log.Fatalf("Error serializing annotations to JSON: %v", err)
	}
	return string(jsonData)
}

