ghligh
======

pdf hi`ghligh`ts swiss knife

ghligh uses JSON as the format to store highlights from PDF files.
It uniquely identifies each PDF by calculating its hash,
ensuring it can accurately track and import highlights associated with the correct file.

### dependencies
install poppler-glib and cairo (for Debian/Ubuntu):
```sh
	apt-get install libpoppler-glib-dev libcairo2-dev
```

### usage
`ghligh` consist of subcommands:
- `import`  import document highlights stored in JSON files into pdf documents
- `export`  export documents' highlights inside a JSON
- `cat`  print documents' highlights
- `ls`  given a set of files/directories show the pdf that contains highlights

```sh
	ghligh export --from file1.pdf --from file2.pdf --to highlighs.json
	ghligh import --from highlights.json --to file1.pdf --to file2.pdf
```

```sh
	ghligh cat file1.pdf file2.pdf
	ghligh ls
	ghligh ls ./economy-books
```
