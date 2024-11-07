module github.com/dirodriguezm/xmatch/catalog_indexer

go 1.23.0

require github.com/stretchr/testify v1.9.0

require github.com/dirodriguezm/xmatch/service v0.0.0

replace github.com/dirodriguezm/xmatch/service => ../service

replace github.com/dirodriguezm/healpix => ../healpix

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dirodriguezm/healpix v0.0.0-20241017225944-6b9a84e4353c
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
