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
	"crypto/sha256"
	"crypto/hmac"

	"fmt"
)

/* hash.go
*  standalone file to access and edit the hash function easily
*/

var ghlighKey = []byte("ghligh")

func (d *GhlighDoc) HashDoc()(string){
	/* This is the magic and cpu expensive algorithm
	*  to hash a document. We use the text of every page
	*  to get a document hash.
	*/
	n_pages := d.doc.GetNPages()

	hmacHash := hmac.New(sha256.New, ghlighKey)

	for i := 0; i < n_pages; i++ {
		page := d.doc.GetPage(i)
		pageText := page.Text()
		page.Close()

		hmacHash.Write([]byte(pageText))
	}

	result := hmacHash.Sum(nil)

	return fmt.Sprintf("%x", result)
}
