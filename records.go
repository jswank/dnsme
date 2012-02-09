package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var records = &Command{
	Run:         runRecords,
	CustomFlags: flagsRecords,
	UsageLine:   "records [filter flags] <domain>",
	Short:       "return records in a domain",
	Long: `
'records' returns a list of all record objects for the specified domain.

Optional filter flags can be specified:

-gtdLocation <DEFAULT | US_EAST | US_WEST | ASIA> is the location.

-type <A | CNAME | MX | NS | PTR | SRV | AAAA | HTTPRED | TXT> is the
record type.

-name <text> an exact match of the record name

-nameContains <text> a substring match of the record name

-value <text> an exact match of the record value

-valueContains <text> an exact match of the record value

`,
}

func flagsRecords(f *flag.FlagSet) {
	f.String("gtdLocation", "", "")
	f.String("type", "", "")
	f.String("name", "", "")
	f.String("nameContains", "", "")
	f.String("value", "", "")
	f.String("valueContains", "", "")
}

func runRecords(cmd *Command, args []string) (err error) {

	if len(args) == 0 {
		err = errors.New("domain not specified")
		return
	}

	domain := args[0]

	values := &url.Values{}
	for _, param := range []string{"gtdLocation", "type", "name", "nameContains", "value", "valueContains"} {
		if cmd.Flag.Lookup(param).Value.String() != "" {
			values.Set(param, cmd.Flag.Lookup(param).Value.String())
		}
	}

	records, err := getDomainRecords(domain, values)
	if err != nil {
		return
	}
	/*
		for _, record := range records {
			if record.Data == "" {
				record.Data = domain + "."
			}
		}
	*/

	switch outputType {
	default:
		{
			for _, record := range records {
				tmpl(os.Stdout, recordTemplate, record)
			}
		}
	case "json":
		{
			b, _ := json.Marshal(records)
			os.Stdout.Write(b)
		}
	case "csv":
		{
			var rs [][]string
			for _, record := range records {
				r := []string{record.Name, strconv.Itoa(record.TTL), record.Type, record.Data, strconv.Itoa(record.ID), record.GtdLocation}
				rs = append(rs, r)
			}
			w := csv.NewWriter(os.Stdout)
			w.WriteAll(rs)
		}
	}
	return

}

var record = &Command{
	Run:         runRecord,
	CustomFlags: flagsRecord,
	UsageLine:   "record -id <record id> <domain>",
	Short:       "returns a specific record id from a domain",
	Long:        "'record' returns a domain record.",
}

func flagsRecord(f *flag.FlagSet) {
	f.String("id", "", "")
}

func runRecord(cmd *Command, args []string) (err error) {

	if len(args) == 0 {
		err = errors.New("domain not specified")
		return
	}

	domain := args[0]

	id := cmd.Flag.Lookup("id").Value.String()
	if id == "" {
		err = errors.New("record id not specified")
		return
	}

	record, err := getDomainRecord(id, domain)
	if err != nil {
		return
	}

	switch outputType {
	default:
		tmpl(os.Stdout, recordTemplate, record)
	case "json":
		{
			b, _ := json.Marshal(record)
			os.Stdout.Write(b)
		}
	case "csv":
		{
			r := []string{record.Name, strconv.Itoa(record.TTL), record.Type, record.Data, strconv.Itoa(record.ID), record.GtdLocation}
			w := csv.NewWriter(os.Stdout)
			w.Write(r)
			w.Flush()
		}
	}

	return

}

var deleteRecord = &Command{
	Run:         runDeleteRecord,
	CustomFlags: flagsDeleteRecord,
	UsageLine:   "delete-record -id <record id> <domain>",
	Short:       "delete a record from the domain",
	Long:        "'delete-record' deleted a record from the domain.",
}

func flagsDeleteRecord(f *flag.FlagSet) {
	f.String("id", "", "record id")
}

func runDeleteRecord(cmd *Command, args []string) (err error) {

	if len(args) == 0 {
		err = errors.New("domain not specified")
		return
	}

	domain := args[0]
	id := cmd.Flag.Lookup("id").Value.String()
	if id == "" {
		err = errors.New("record id not specified")
		return
	}

	req, err := http.NewRequest("DELETE", api_url+"/domains/"+domain+"/records/"+id, nil)
	if err != nil {
		return
	}
	addDnsmeHeaders(req)

	_, err = makeRequest(req)
	if err != nil {
		return
	}

	return
}

var updateRecord = &Command{
	Run:         runUpdateRecord,
	CustomFlags: flagsUpdateRecord,
	UsageLine: `update-record -id <record id> -name <name> -data <record data>
    [-ttl <ttl>] [-type <record type>] [-gtdLocation <gtdLocation>] 
    [-password <password>] <domain>`,
	Short: "update an existing record",
	Long: `
'update-record' updates an existing record object in the specified
domain.
 
-id is the unique record identifier.

-name <record name> is the record name. An empty value indicates that
the base domain is used.

-data is the record data. Content varies based on record type:
    A: <host IP>
    IPv6 host IP>
    CNAME: <target name>
    MX: <priority> <target name>
    NS: <name server>
    PTR: <target name>
    SRV: <priority> <weight> <port> <target name>
    TXT: <text value>

-ttl is the amount of time in seconds a record will be cached before
being refreshed.  Default value is "3600" (one hour).

-type (optional) is the record type. Values: A, AAAA, CNAME, MX, NS,
PTR, SRV, TXT

-gtdLocation (optional) is the Global Traffic Director location. Values:
DEFAULT, US_EAST, US_WEST, EUROPE

-password is the password required for dynamic DNS updates

`,
}

func flagsUpdateRecord(f *flag.FlagSet) {
	f.String("id", "", "")
	f.String("name", "", "")
	f.String("type", "", "")
	f.String("data", "", "")
	f.String("ttl", "3600", "")
	f.String("gtdLocation", "DEFAULT", "")
	f.String("password", "", "")
}

func runUpdateRecord(cmd *Command, args []string) (err error) {

	if len(args) == 0 {
		err = errors.New("domain not specified")
		return
	}

	domain := args[0]

	rec := &apiRecord{}
	rec.ID, _ = strconv.Atoi(cmd.Flag.Lookup("id").Value.String())
	rec.Name = cmd.Flag.Lookup("name").Value.String()
	rec.Type = cmd.Flag.Lookup("type").Value.String()
	rec.Data = cmd.Flag.Lookup("data").Value.String()
	rec.TTL, _ = strconv.Atoi(cmd.Flag.Lookup("ttl").Value.String())
	rec.GtdLocation = cmd.Flag.Lookup("gtdLocation").Value.String()
	rec.Password = cmd.Flag.Lookup("password").Value.String()

	_, err = addDomainRecord(domain, *rec)
	if err != nil {
		return
	}

	return

}

var addRecord = &Command{
	Run:         runAddRecord,
	CustomFlags: flagsUpdateRecord,
	UsageLine: `add-record -name <name> -type <record type> [-ttl <ttl>]
    -data <record data> [-gtdLocation <gtdLocation>] [-password <password>]
    <domain>`,
	Short: "add a new record",
	Long: `
'add-record' adds a record object to the specified domain.

-name <record name> is the record name. An empty value indicates that
the base domain is used.

-data is the record data. Content varies based on record type:
    A: <host IP>
    IPv6 host IP>
    CNAME: <target name>
    MX: <priority> <target name>
    NS: <name server>
    PTR: <target name>
    SRV: <priority> <weight> <port> <target name>
    TXT: <text value>

-ttl is the amount of time in seconds a record will be cached before
being refreshed.  Default value is 3600 (one hour).

-type (optional) is the record type. Values: A, AAAA, CNAME, MX, NS,
PTR, SRV, TXT

-gtdLocation is the Global Traffic Director location. Values:
DEFAULT, US_EAST, US_WEST, EUROPE

-password is the password required for dynamic DNS updates

`,
}

func runAddRecord(cmd *Command, args []string) (err error) {

	if len(args) == 0 {
		err = errors.New("domain not specified")
		return
	}

	domain := args[0]

	rec := &apiRecord{}
	rec.Name = cmd.Flag.Lookup("name").Value.String()
	rec.Type = cmd.Flag.Lookup("type").Value.String()
	rec.Data = cmd.Flag.Lookup("data").Value.String()
	rec.TTL, _ = strconv.Atoi(cmd.Flag.Lookup("ttl").Value.String())
	rec.GtdLocation = cmd.Flag.Lookup("gtdLocation").Value.String()

	if rec.Type == "" || rec.Data == "" || rec.TTL == 0 {
		err = errors.New("missing required parameters")
		return
	}

	record, err := addDomainRecord(domain, *rec)
	if err != nil {
		return
	}

	if record.Data == "" {
		record.Data = domain + "."
	}

	switch outputType {
	default:
		tmpl(os.Stdout, recordTemplate, record)
	case "json":
		{
			b, _ := json.Marshal(record)
			os.Stdout.Write(b)
		}
	case "csv":
		{
			r := []string{record.Name, strconv.Itoa(record.TTL), record.Type, record.Data, strconv.Itoa(record.ID), record.GtdLocation}
			w := csv.NewWriter(os.Stdout)
			w.Write(r)
			w.Flush()
		}
	}

	return

}
