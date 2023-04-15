package main

import "github.com/spf13/pflag"

var (
	cfg = pflag.StringP("csv", "", "./data.csv", "import data from vector csv file.")
)

func main() {
	pflag.Parse()
}
