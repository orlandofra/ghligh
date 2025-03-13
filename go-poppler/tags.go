package poppler

// #cgo pkg-config: poppler-glib
// #include <poppler.h>
// #include <stdlib.h>
// #include <glib.h>
// #include <unistd.h>
import "C"
import "strings"



var zeroRect = Rectangle{X1: 0, X2: 0, Y1: 0, Y2: 0}

func (d *Document) Tag(text string){
	am := C.poppler_annot_mapping_new();

	pRect := rectangleToPopplerRectangle(zeroRect)

	annot := Annot {
		am: am,
	}
	defer annot.Close()

	am.annot = C.poppler_annot_text_new(d.doc, &pRect)
	annot.SetContents(text)
	annot.SetFlags(AnnotFlagHidden | AnnotFlagInvisible)

	zeroPage := d.GetPage(0)
	zeroPage.AddAnnot(annot)
	defer zeroPage.Close()
}


func (d *Document) GetTags(filter string) []*Annot {
	page := d.GetPage(0)
	//defer page.Close()

	annots := page.GetAnnots()
	var tags []*Annot
	for _, a := range(annots) {
		if a.Type() == AnnotText &&
		  rectEq(a.Rect(), zeroRect) &&
		  a.Flags() & AnnotFlagHidden != 0 &&
		  a.Flags() & AnnotFlagInvisible != 0 &&
		  strings.HasPrefix(a.Contents(), filter){
			tags = append(tags, a)
		}
	}

	return tags
}

func (d *Document) RemoveTags(filter string){
//	d.GetPage(0).AddAnnot(annot)
	// TODO 
}
