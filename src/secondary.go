package main

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"strings"
)

var listSecondaries = &Command{
	Run:       runListSecondaries,
	UsageLine: "secondaries",
	Short:     "lists all secondary domains",
	Long:      "'secondaries' lists all secondary domains currently available.",
}

func runListSecondaries(cmd *Command, args []string) (err error) {

	domains, err := getSecondaryList()
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

var infoSecondary = &Command{
	Run: runInfoSecondary,
	//	CustomFlags: false,
	UsageLine: "secondary <domain>",
	Short:     "returns information about a secondary domain",
	Long:      "'secondary' returns information about a secondary domain.",
}

func runInfoSecondary(cmd *Command, args []string) (err error) {

	if len(args) == 0 {
		err = errors.New("domain not specified")
		return
	}

	domain := args[0]

	info, err := getSecondary(domain)
	if err != nil {
		return
	}

	switch outputType {
	default:
		{
			tmpl(os.Stdout, secondaryTemplate, info)
		}
	case "csv":
		{
			tmpl(os.Stdout, secondaryTemplateCSV, info)
		}
	case "json":
		{
			b, _ := json.Marshal(info)
			os.Stdout.Write(b)
		}
	}

	return

}

var delSecondary = &Command{
	Run:       runDeleteSecondary,
	UsageLine: "delete-secondary <domain>",
	Short:     "deletes a secondary domain",
	Long:      "'delete-secondary' removes a secondary domain.",
}

func runDeleteSecondary(cmd *Command, args []string) (err error) {

	if len(args) == 0 {
		err = errors.New("domain not specified")
		return
	}

	domain := args[0]

	err = deleteSecondary(domain)
	if err != nil {
		return
	}

	return

}

var addNewSecondary = &Command{
	Run:         runAddSecondary,
	CustomFlags: flagsAddSecondary,
	UsageLine:   "add-secondary -ip <list of ip addresses> <domain>",
	Short:       "adds a secondary domains",
	Long: `
'add-secondary' creates a secondary domain entry with the specified
name.  If the secondary domain already exists, then the existing IP
addresses are replaced with those specified.

The -ip flag is a list of comma separated IP addresses defining the
name servers which are masters for the domain, i.e. where updates will 
be accepted from:

Example:
  $ ./dnsme add-secondary -o json -ip "127.0.0.1, 127.0.0.2" example.com
  {"name":"example.com","ip":["127.0.0.2","127.0.0.1"]}

	`,
}

func flagsAddSecondary(f *flag.FlagSet) {
	f.String("ip", "", "")
}

func runAddSecondary(cmd *Command, args []string) (err error) {

	if len(args) == 0 {
		err = errors.New("domain not specified")
	}

	secondary := &apiSecondary{}
	secondary.Name = args[0]
	if cmd.Flag.Lookup("ip").Value.String() != "" {
		for _, ns := range strings.Split(cmd.Flag.Lookup("ip").Value.String(), ",") {
			secondary.IP = append(secondary.IP, strings.TrimSpace(ns))
		}
	}

	info, err := addSecondary(*secondary)
	if err != nil {
		return
	}

	switch outputType {
	default:
		{
			tmpl(os.Stdout, secondaryTemplate, info)
		}
	case "csv":
		{
			tmpl(os.Stdout, secondaryTemplateCSV, info)
		}
	case "json":
		{
			b, _ := json.Marshal(info)
			os.Stdout.Write(b)
		}
	}

	return

}
