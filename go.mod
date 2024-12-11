module github.com/orlandofra/ghligh

go 1.22.2

replace github.com/orlandofra/ghligh/ghligh => ./ghligh

replace github.com/orlandofra/ghligh/go-poppler => ./go-poppler

require github.com/orlandofra/ghligh/ghligh v0.0.0-00010101000000-000000000000

require (
	github.com/orlandofra/ghligh/go-poppler v0.0.0-00010101000000-000000000000 // indirect
	github.com/ungerik/go-cairo v0.0.0-20240304075741-47de8851d267 // indirect
)
