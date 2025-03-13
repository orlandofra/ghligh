package document

import (
	"github.com/orlandofra/ghligh/go-poppler"

	"os"
	"sync"

	"strings"

	"fmt"
)

const ghlighFilter = "ghligh-Y2lhbm5v:"

// This is different from poppler's annot_mapping
// it is the list of annotations mapped to the page index
type AnnotsMap map[int][]AnnotJSON

type GhlighDoc struct {
	doc *poppler.Document
	mu  sync.Mutex

	Path         string    `json:"file"`
	HashBuffer   string    `json:"hash"`
	AnnotsBuffer AnnotsMap `json:"highlights,omitempty"`
}

type HighlightedText struct {
	Page     int    `json:"page"`
	Text     string `json:"text"`
	Contents string `json:"contents,omitempty"`
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
	// HashDoc??

	return g, nil
}

func (d *GhlighDoc) Close() {
	d.AnnotsBuffer = nil
	d.HashBuffer = ""
	if d.doc != nil {
		d.doc.Close()
	}
}

func (d *GhlighDoc) Info() poppler.DocumentInfo {
	return d.doc.Info()
}

func (d *GhlighDoc) tagExists(text string) bool {
	for _, tag := range d.GetTags() {
		if tag == text {
			return true
		}
	}
	return false
}

func (d *GhlighDoc) Tag(text string) {
	if !d.tagExists(text) {
		d.doc.Tag(ghlighFilter + text)
	} else {
		fmt.Fprintf(os.Stderr, "warning: tag %s already exist inside %s, i don't do anything\n", text, d.Path)
	}
}

func (d *GhlighDoc) GetTags() []string {
	var tags []string
	annots := d.doc.GetTags(ghlighFilter)
	for _, annot := range annots {
		contents := strings.TrimPrefix(annot.Contents(), ghlighFilter)
		tags = append(tags, contents)
	}
	return tags
}

func (d *GhlighDoc) RemoveTags(tags []string) int {
	zeroPage := d.doc.GetPage(0)
	var removedTags int

	annots := d.doc.GetTags(ghlighFilter)
	for _, annot := range annots {
		contents := strings.TrimPrefix(annot.Contents(), ghlighFilter)
		for _, tag := range tags {
			if tag == contents {
				zeroPage.RemoveAnnot(*annot)
				removedTags += 1
				break
			}
		}
	}
	return removedTags
}

func (d *GhlighDoc) Import(annotsMap AnnotsMap) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	annots_count := 0

	var err error
	d.AnnotsBuffer = annotsMap

	for key := range d.AnnotsBuffer {
		page := d.doc.GetPage(key)
		for _, annot := range d.AnnotsBuffer[key] {
			a := d.jsonToAnnot(annot)
			if !isInPage(a, page) {
				annots_count += 1
				page.AddAnnot(*a)
			}

		}
		page.Close()
	}

	d.AnnotsBuffer = nil
	return annots_count, err
}

func integrityCheck(tizio *GhlighDoc, caio *GhlighDoc) {

}

func (d *GhlighDoc) GetNPages() int {
	return d.doc.GetNPages()
}

func (d *GhlighDoc) GetPageText(i int) (string, error) {
	nPages := d.doc.GetNPages()

	if i < 0 || i > nPages {
		return "", fmt.Errorf("error page %d out of range %d", i, nPages)
	}

	p := d.doc.GetPage(i)
	defer p.Close()

	text := p.Text()
	return text, nil
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

	if newDoc.HashDoc() != d.HashDoc() {
		return false, fmt.Errorf("After saving document %s to %s its hash doesn't correspond the the old one", d.Path, tempFile.Name())
	}

	err = os.Rename(tempFile.Name(), d.Path)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (d *GhlighDoc) Cat() []HighlightedText {
	var highlights []HighlightedText

	n_pages := d.doc.GetNPages()
	for i := 0; i < n_pages; i++ {
		page := d.doc.GetPage(i)
		annots := page.GetAnnots()
		for _, annot := range annots {
			if annot.Type() == poppler.AnnotHighlight {
				annotText := page.AnnotText(*annot)

				highlights = append(highlights, HighlightedText{Page: i, Text: annotText, Contents: annot.Contents()})
			}
		}

		page.Close()
	}
	return highlights
}

func (d *GhlighDoc) HasHighlights() bool {
	// check if is tagged with ls
	if d.tagExists("ls") {
		return true
	}

	// check if it has highlights
	n_pages := d.doc.GetNPages()
	for i := 0; i < n_pages; i++ {
		page := d.doc.GetPage(i)
		annots := page.GetAnnots()
		for _, annot := range annots {
			if annot.Type() == poppler.AnnotHighlight {
				return true
			}
		}

		page.Close()
	}
	return false
}

func (d *GhlighDoc) GetAnnotsBuffer() AnnotsMap {
	annots_json_of_page := make(AnnotsMap)

	n := d.doc.GetNPages()
	var annots_json []AnnotJSON
	for i := 0; i < n; i++ {
		annots_json = nil
		page := d.doc.GetPage(i)

		annots := page.GetAnnots()
		for _, annot := range annots {
			if annot.Type() == poppler.AnnotHighlight {
				annot_json := annotToJson(*annot)
				annots_json = append(annots_json, annot_json)
			}
		}

		page.Close()

		if len(annots_json) > 0 {
			annots_json_of_page[i] = annots_json
		}
	}

	return annots_json_of_page
}
