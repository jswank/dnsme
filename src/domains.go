package main

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"strings"
)

var listDomains = &Command{
	Run:       runListDomains,
	UsageLine: "list",
	Short:     "lists all domains",
	Long:      "'list' lists all domains currently available.",
}

func runListDomains(cmd *Command, args []string) (err error) {

	domains, err := getDomainList()
	if err != nil {
		return
	}

	switch outputType {
	default:
		{
			tmpl(os.Stdout, domainListTemplate, domains)
		}
	case "json":
		{
			b, _ := json.Marshal(domains.List)
			os.Stdout.Write(b)
		}
	}

	return

}

var infoDomain = &Command{
	Run: runInfoDomain,
	//	CustomFlags: false,
	UsageLine: "info <domain>",
	Short:     "returns information about a domain",
	Long:      "'info' returns information about a domain.",
}

func runInfoDomain(cmd *Command, args []string) (err error) {

	if len(args) == 0 {
		err = errors.New("domain not specified")
		return
	}

	domain := args[0]

	info, err := getDomainInfo(domain)
	if err != nil {
		return
	}

	switch outputType {
	default:
		{
			tmpl(os.Stdout, domainInfoTemplate, info)
		}
	case "csv":
		{
			tmpl(os.Stdout, domainInfoTemplateCSV, info)
		}
	case "json":
		{
			b, _ := json.Marshal(info)
			os.Stdout.Write(b)
		}
	}

	return

}

var delDomain = &Command{
	Run:       runDeleteDomain,
	UsageLine: "delete-domain <domain>",
	Short:     "deletes a domain",
	Long:      "'delete-domain' removes a domain.",
}

func runDeleteDomain(cmd *Command, args []string) (err error) {

	if len(args) == 0 {
		err = errors.New("domain not specified")
		return
	}

	domain := args[0]

	err = deleteDomain(domain)
	if err != nil {
		return
	}

	return

}

var addNewDomain = &Command{
	Run:         runAddDomain,
	CustomFlags: flagsAddDomain,
	UsageLine:   "add-domain [-ns <list of nameservers>] [-gtdEnabled <true|false>] <domain>",
	Short:       "adds a domains",
	Long: `
'add-domain' creates a domain entry with the specified name.

The -ns flag is a list of comma separated strings defining the name servers associated
with this domain.

The -gtdEnabled flag defines whether or not this domain uses the Global Traffic Director.
Default value is false.

	`,
}

func flagsAddDomain(f *flag.FlagSet) {
	f.String("ns", "", "")
	f.Bool("gtd", false, "")
}

func runAddDomain(cmd *Command, args []string) (err error) {

	if len(args) == 0 {
		err = errors.New("domain not specified")
	}

	domain := &apiDomain{}
	domain.Name = args[0]
	if cmd.Flag.Lookup("ns").Value.String() != "" {
		for _, ns := range strings.Split(cmd.Flag.Lookup("ns").Value.String(), ",") {
			domain.NameServers = append(domain.NameServers, strings.TrimSpace(ns))
		}
	}
	if cmd.Flag.Lookup("gtd").Value.String() == "true" {
		domain.GtdEnabled = true
	}

	info, err := addDomain(*domain)
	if err != nil {
		return
	}

	switch outputType {
	default:
		{
			tmpl(os.Stdout, domainInfoTemplate, info)
		}
	case "csv":
		{
			tmpl(os.Stdout, domainInfoTemplateCSV, info)
		}
	case "json":
		{
			b, _ := json.Marshal(info)
			os.Stdout.Write(b)
		}
	}

	return

}
