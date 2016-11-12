package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/jaden-young/murder-hobos/util/xmlspellparse"
)

var fpath string

func init() {
	flag.StringVar(&fpath, "f", "", "file to open")
}

func main() {
	flag.Parse()

	file, err := os.Open(fpath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	b, _ := ioutil.ReadAll(file)

	var c xmlspellparse.Compendium
	xml.Unmarshal(b, &c)

	spew.Dump(c)

}
