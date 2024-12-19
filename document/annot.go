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

func annotToJson(a poppler.Annot) (AnnotJSON) {
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

func (d *GhlighDoc) jsonToAnnot(aJson AnnotJSON) *poppler.Annot {

	annot, _ := d.doc.NewAnnot(poppler.AnnotHighlight, aJson.Rect, aJson.Quads)

	annot.SetColor(aJson.Color)
	annot.SetContents(aJson.Contents)
	annot.SetFlags(aJson.Flags)

	return &annot
}

func popplerAnnotsMatch(a *poppler.Annot, b *poppler.Annot) bool {
	aRect := a.Rect()
	bRect := b.Rect()

	aQuads := a.Quads()
	bQuads := b.Quads()

	/* refactor into rect == rect funct */
	if aRect.X1 != bRect.X1 ||
		aRect.Y1 != bRect.Y1 ||
		aRect.X2 != bRect.X2 ||
		aRect.Y2 != bRect.Y2 {
		return false
	}

	if len(aQuads) != len(bQuads) {
		return false
	}

	/* refactor into quads == quads funct */
	for i := range aQuads {
		q1 := aQuads[i]
		q2 := bQuads[i]

		if q1.P1.X != q2.P1.X || q1.P1.Y != q2.P1.Y ||
			q1.P2.X != q2.P2.X || q1.P2.Y != q2.P2.Y ||
			q1.P3.X != q2.P3.X || q1.P3.Y != q2.P3.Y ||
			q1.P4.X != q2.P4.X || q1.P4.Y != q2.P4.Y {
			return false
		}
	}

	return true
}
