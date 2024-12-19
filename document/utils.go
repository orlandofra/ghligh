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

	"encoding/json"

)

func (d *GhlighDoc) loadAnnots() AnnotsMap {
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


func unmarshallHighlights(jsonData string) (AnnotsMap, error) {
	var annotsMap AnnotsMap

	err := json.Unmarshal([]byte(jsonData), &struct {
		Highlights *AnnotsMap `json:"highlights"`
	}{
		Highlights: &annotsMap,
	})

	return annotsMap, err
}


func isInPage(a *poppler.Annot, p *poppler.Page) bool {
	annots := p.GetAnnots()
	for _, annot := range(annots){
		if popplerAnnotsMatch(a, annot){
			return true
		}
	}

	return false
}

