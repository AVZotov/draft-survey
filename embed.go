package draft_survey

import (
	"embed"
)

//go:embed data/dictionaries/packing.json
//go:embed data/dictionaries/cargo_types.json
//go:embed data/dictionaries/countries.json
//go:embed data/dictionaries/ports.json
var Dictionaries embed.FS

//go:embed docs/glossary_en.md
//go:embed docs/glossary_ru.md
//go:embed docs/CALCULATION.md
var Docs embed.FS

//go:embed web/static
var StaticFiles embed.FS
