/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/orlandofra/ghligh/document"
	"github.com/spf13/cobra"
)

type Accumulatore struct {
	value int
}

func newAccumulatore() *Accumulatore {
	return &Accumulatore{
		value: 0,
	}
}

func (acc *Accumulatore) Goto(currentPage int) *Accumulatore {
	if acc.value == 0 {
		// if no argument specified go to current page (refresh screen)
		acc.value = currentPage
	} else {
		// accumualtor seems off by one
		acc.value -= 1
	}

	return acc
}

func (acc *Accumulatore) Prev() *Accumulatore {
	if acc.value == 0 {
		acc.value = 1
	}
	acc.value *= -1
	return acc
}

func (acc *Accumulatore) Next() *Accumulatore {
	if acc.value == 0 {
		acc.value = 1
	}
	return acc
}

func (acc *Accumulatore) Pop() int {
	oldValue := acc.value
	acc.value = 0
	return oldValue
}

func (acc *Accumulatore) Value() int {
	return acc.value
}

func (acc *Accumulatore) EatRune(r rune) {
	acc.value = acc.value*10 + int(r-'0')
}

type browsableDoc struct {
	doc         *document.GhlighDoc
	currentPage int

	title  string
	author string
	nPages int
}

func (d *browsableDoc) getCurrentPage() int {
	return d.currentPage
}
func (d *browsableDoc) setCurrentPage(i int) {
	switch {
	case i >= d.nPages:
		d.currentPage = d.nPages - 1
	case i < 0:
		d.currentPage = 0
	default:
		d.currentPage = i
	}
}

func openBrowsableDoc(path string) (*browsableDoc, error) {
	// open and init document
	doc, err := document.Open(path)
	if err != nil {
		return nil, err
	}

	b := &browsableDoc{
		doc:         doc,
		currentPage: 0,
		title:       doc.Info().Title,
		author:      doc.Info().Author,
		nPages:      doc.GetNPages(),
	}
	return b, nil
}

func (d *browsableDoc) header() string {
	return fmt.Sprintf("[yellow]%s - %s",
		d.title,
		d.author,
	)
}

func (d *browsableDoc) status(accumulator int) string {

	accumulated := ""
	if accumulator != 0 {
		accumulated = fmt.Sprintf("[%d] --", accumulator)
	}

	return fmt.Sprintf("%s page %d of %d",
		accumulated,
		d.currentPage+1,
		d.nPages,
	)

}

func (d *browsableDoc) getCurrentPageText() string {
	text, _ := d.doc.GetPageText(d.currentPage)
	return text
}

type Browser struct {
	app         *tview.Application
	pageContent *tview.TextArea
	header      *tview.TextView
	status      *tview.TextView
	layout      *tview.Flex

	docs    []*browsableDoc
	currDoc int

	accumulator *Accumulatore
}

func newBrowser(paths []string) *Browser {
	app := tview.NewApplication()
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	// a read-only textArea is more convenient than a textView
	pageContent := tview.NewTextArea().
		SetWordWrap(true)

	status := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 1, 0, false).
		AddItem(pageContent, 0, 1, true).
		AddItem(status, 1, 0, false)

	var docs []*browsableDoc
	for _, path := range paths {
		doc, err := openBrowsableDoc(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
		}
		docs = append(docs, doc)
	}

	// make the textArea readonly, we will handle the navigation
	// completely
	pageContent.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		default:
			return nil
		}
	})

	app.SetFocus(status)

	return &Browser{
		app:         app,
		header:      header,
		pageContent: pageContent,
		status:      status,
		layout:      layout,
		docs:        docs,
		accumulator: newAccumulatore(),
	}

}

func (b *Browser) currentDoc() *browsableDoc {
	return b.docs[b.currDoc]
}

func (b *Browser) updateStatus() {
	b.status.SetText(b.currentDoc().status(b.accumulator.Value()))
}

func (b *Browser) writeStatus(s string) {
	b.status.SetText(s)
}

func (b *Browser) getLineText() string {
	_, y, _, _ := b.pageContent.GetCursor()
	text := b.pageContent.GetText()
	return strings.Split(text, "\n")[y]
}

func (b *Browser) updateHeader() {
	b.header.SetText(b.currentDoc().header())
}

func (b *Browser) updateContent() {
	text := b.currentDoc().getCurrentPageText()

	b.pageContent.
		SetText(text, false).
		SetBorder(true)
}

func (b *Browser) updatePage(newPage int) {
	// lock mutex here
	b.currentDoc().setCurrentPage(newPage)
	b.updateHeader()
	b.updateContent()
	// unlock mutex here

	b.updateStatus()
}

func (b *Browser) setCurrentDoc(i int) {
	if len(b.docs) == 0 {
		return
	}

	b.currDoc = (b.currDoc + i + len(b.docs)) % len(b.docs)
}

func (b *Browser) Run() {
	b.updatePage(0)
	b.app.SetInputCapture(b.handle)
	if err := b.app.SetRoot(b.layout, true).Run(); err != nil {
		panic(err)
	}
}

func (b *Browser) currentPage() int {
	return b.currentDoc().getCurrentPage()
}
func (b *Browser) handle(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		b.accumulator.Pop()
		b.updateStatus()
	case tcell.KeyRune:
		switch event.Rune() {
		case 'Q':
			b.app.Stop()
			return event
			//TODO		case 'q': QUIT CURRENT DOCUMENT
			//			b.app.Stop()
			//			return event
		//case 'v':
		//if b.pageContent.HasFocus() {
		//b.app.SetFocus(b.status)
		//} else {
		//b.app.SetFocus(b.pageContent)
		//}
		case 'g':
			b.updatePage(b.accumulator.Goto(b.currentPage()).Pop())
		case 'G':
			b.updatePage(b.accumulator.Goto(b.currDoc).Pop())
		case 'N':
			b.setCurrentDoc(b.accumulator.Next().Pop())
			b.updatePage(b.currentPage())
		case 'P':
			b.setCurrentDoc(b.accumulator.Prev().Pop())
			b.updatePage(b.currentPage())
		case 'n':
			b.updatePage(b.currentPage() + b.accumulator.Next().Pop())
		case 'p':
			b.updatePage(b.currentPage() + b.accumulator.Prev().Pop())
		case 'j':
			//	b.pageContent.ScrollTo(b.accumulator.Next().Pop(), 0)
			b.updateStatus()
		case 'k':
			//	b.pageContent.ScrollTo(b.accumulator.Prev().Pop(), 0)
			b.updateStatus()
		default:
			if event.Rune() >= '0' && event.Rune() <= '9' {
				b.accumulator.EatRune(event.Rune())
				b.updateStatus()
			} else {
				// reset accumulator
				b.accumulator.Pop()
				b.updateStatus()
			}
		}
	}
	return event
}

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse pdf file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		var docs []*document.GhlighDoc
		for _, arg := range args {
			doc, err := document.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not open %s: %v\n", arg, err)
				continue
			}
			docs = append(docs, doc)
		}

		if len(docs) == 0 {
			fmt.Fprintf(os.Stderr, "please specify at least one valid document\n")
			return
		}

		// if more then one file show menu
		browser := newBrowser(args)
		browser.Run()

	},
}

func init() {
	rootCmd.AddCommand(browseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// browseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// browseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
