// Package ghligh, a pdf highlights swiss knife
//
// Copyright (c) 2024 Francesco Orlando
//
// This file is part of a program licensed under the GNU General Public License, version 2.
// You should have received a copy of the GNU General Public License along with this program.
// If not, see the LICENSE.md file in the root directory of this repository or visit
// <https://www.gnu.org/licenses/old-licenses/gpl-2.0.html>.
package ghligh


import (
	"github.com/cheggaaa/go-poppler"

	"os"
	"encoding/json"
	"sync"

	"fmt"
)

// This is different from poppler's annot_mapping
// it is the list of annotations mapped to the page index
type AnnotsMap map[int][]AnnotJSON

type GhlighDoc struct {
	doc		*poppler.Document
	mu		sync.Mutex

	Path		string			`json:"file"`
	Hash		string			`json:"hash"`
	AnnotsBuffer	AnnotsMap		`json:"highlights"`
}

func Open(filename string) (*GhlighDoc, error) {
	var err error

	g := &GhlighDoc{}

	g.doc, err = poppler.Open(filename)
	if err != nil {
		fmt.Errorf("%s: error opening pdf %v", os.Args[0], err)
		return nil, err
	}
	g.Path = filename

	return g, nil
}

func (d *GhlighDoc) Close() {
	d.AnnotsBuffer = nil
	d.Hash = ""
	if d.doc != nil {
		d.doc.Close()
	}
}

func (d *GhlighDoc) Export() (json.RawMessage, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.AnnotsBuffer = d.loadAnnots()
	d.Hash = d.HashDoc()
	defer func() {
		d.AnnotsBuffer = nil
		d.Hash = ""
	}()

	jsonData, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (d *GhlighDoc) Import (annotsMap AnnotsMap) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	annots_count := 0

	var err error
	d.AnnotsBuffer = annotsMap

	for key := range d.AnnotsBuffer {
		page := d.doc.GetPage(key)
		for _, annot := range(d.AnnotsBuffer[key]){
			a := d.jsonToAnnot(annot)
			if !isInPage(a, page){
				annots_count += 1
				page.AddAnnot(*a)
			}

		}
		page.Close()
	}

	d.AnnotsBuffer = nil
	return annots_count, err
}


func integrityCheck(tizio *GhlighDoc, caio *GhlighDoc){

}

func (d *GhlighDoc) Save() (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	tempFile, err := os.CreateTemp("", ".ghligh_*.pdf")
	if err != nil {
		return false, err
	}
	defer os.Remove(tempFile.Name())


	ok, err := d.doc.Save(tempFile.Name())
	if !ok {
		return false, err
	}

	/* integrity check */
	newDoc, err := Open(tempFile.Name())
	if err != nil {
		return false, err
	}

	if newDoc.HashDoc() != d.HashDoc(){
		return false, fmt.Errorf("After saving document %s to %s its hash doesn't correspond the the old one", d.Path, tempFile.Name())
	}

	err = os.Rename(tempFile.Name(), d.Path)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (d *GhlighDoc) Cat(ch chan string) {
	defer close(ch)


	n_pages := d.doc.GetNPages()
	for i := 0; i < n_pages; i++ {
		page := d.doc.GetPage(i)
		annots := page.GetAnnots()
		for _, annot := range(annots){
			if annot.Type() == poppler.AnnotHighlight {
				annotText := page.AnnotText(*annot)
				ch <- annotText
			}
		}

		page.Close()
	}
}

func (d *GhlighDoc) HasHighlights() bool{
	n_pages := d.doc.GetNPages()
	for i := 0; i < n_pages; i++ {
		page := d.doc.GetPage(i)
		annots := page.GetAnnots()
		for _, annot := range(annots){
			if annot.Type() == poppler.AnnotHighlight {
				return true
			}
		}

		page.Close()
	}
	return false
}
