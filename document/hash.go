package document

import (
	"fmt"
	"os"

	"reflect"
	"unsafe"

	"math"

	"runtime"
	"sync"

	"crypto/hmac"
	"crypto/sha256"

	"github.com/orlandofra/ghligh/go-poppler"
)

var ghlighKey = []byte("ghligh-pdf-doc")

var bufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, os.Getpagesize())
	},
}

type pageResult struct {
	index int
	buf   []byte
}

func sqrtInt(n int) int {
	return int(math.Sqrt(float64(n)))
}

func continueAt(i, n int) bool {
	// Very unlikely to edit a pdf and add a page in the center
	return i < sqrtInt(n)/2 || i > n-sqrtInt(n)/2
}

// generate identifier from document based on document text (use layout instead)
func (d *GhlighDoc) HashDoc() string {
	nPages := d.doc.GetNPages()

	hmacHash := hmac.New(sha256.New, ghlighKey)
	resultsCh := make(chan pageResult, nPages)

	var wg sync.WaitGroup

	maxWorkers := runtime.NumCPU() + 1
	sem := make(chan struct{}, maxWorkers)

	for i := 0; continueAt(i, nPages); i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()

			page := d.doc.GetPage(i)
			text := page.Text()
			page.Close()

			buf := bufPool.Get().([]byte)
			buf = buf[:0]
			buf = append(buf, []byte(text)...)

			resultsCh <- pageResult{index: i, buf: buf}
		}(i)
	}

	wg.Wait()
	close(resultsCh)

	results := make([][]byte, nPages)
	for res := range resultsCh {
		results[res.index] = res.buf
	}

	for i := 0; continueAt(i, nPages); i++ {
		hmacHash.Write(results[i])
		hmacHash.Write([]byte{byte(i)})
		bufPool.Put(results[i])
	}
	hmacHash.Write([]byte{byte(nPages)})

	return fmt.Sprintf("%x", hmacHash.Sum(nil))
}

func rectToBytes(r *poppler.Rectangle) []byte {
	size := int(unsafe.Sizeof(*r))

	var sliceHeader reflect.SliceHeader
	sliceHeader.Data = uintptr(unsafe.Pointer(r))
	sliceHeader.Len = size
	sliceHeader.Cap = size

	return *(*[]byte)(unsafe.Pointer(&sliceHeader))
}
