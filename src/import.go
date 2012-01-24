package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

var importData = &Command{
	Run:         runImport,
	CustomFlags: flagsImport,
	UsageLine:   "import [-file <export file>]",
	Short:       "import domain info & records",
	Long: `'import' imports JSON-encoded information into DNS Made Easy

If no export file is specified, then standard input is used.
	
`,
}

func flagsImport(f *flag.FlagSet) {
	f.String("file", "-", "Import file")
}

func getReader(input string) (reader io.Reader, err error) {

	if input == "" || input == "-" {
		reader = os.Stdin
	} else {
		reader, err = os.Open(input)
		if err != nil {
			return
		}
	}
	return
}

func runImport(cmd *Command, args []string) (err error) {

	var import_domains []exportDomain

	// open file
	r, err := getReader(cmd.Flag.Lookup("file").Value.String())
	if err != nil {
		return
	}

	// parse it
	d := json.NewDecoder(r)
	err = d.Decode(&import_domains)
	if err != nil {
		return
	}

	for _, d := range import_domains {
		// get domain info
		// if it does not exist, create it
		_, err = getDomainInfo(d.Domain.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			_, e := addDomain(d.Domain)
			if e != nil {
				fmt.Fprintf(os.Stderr, "couldn't create domain %s: , %s", d.Domain.Name, e)
				continue
			}
		}

		for _, record := range d.Records {
			record.ID = 0
			_, err = addDomainRecord(d.Domain.Name, record)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error adding record to domain %s: %+v, %s", d.Domain.Name, record, err)
			}
		}
	}

	return

}
