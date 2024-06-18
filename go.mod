module github.com/cosys-io/cosys

go 1.22.3

replace github.com/cosys-io/cosys/experiment/invoice/invoicer => ./experiment/invoice/invoicer

require (
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/spf13/cobra v1.8.0
	golang.org/x/text v0.14.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)
