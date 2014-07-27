package main

import (
	"encoding/json"
	"fmt"
	"sort"
	//	"os"
)

type exportDomain struct {
	Domain  apiDomain   `json:"domain"`
	Records []apiRecord `json:"records"`
}

var exportData = &Command{
	Run: runExport,
	//	CustomFlags: flagsExport,
	UsageLine: "export [<domain>]",
	Short:     "export domain info & records",
	Long:      "'export' returns all domain information suitable for importing",
}

func runExport(cmd *Command, args []string) (err error) {

	var domains apiDomainList

	if len(args) > 0 {
		domains.List = args[0:]
	} else {
		domains, err = getDomainList()
		if err != nil {
			return
		}
		sort.Strings(domains.List)
	}

	var exportDomains []exportDomain

	for _, domain := range domains.List {
		var d exportDomain
		d.Domain, err = getDomainInfo(domain)
		if err != nil {
			return
		}
		d.Records, err = getDomainRecords(domain, nil)
		if err != nil {
			return
		}
		exportDomains = append(exportDomains, d)
	}

	b, err := json.Marshal(exportDomains)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", b)
	//	outputExportDomains(exportDomains)

	return

}
