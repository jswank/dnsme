# dnsme

This is a command line tool that can be used to manage DNS information
hosted at DNS Made Easy (http://www.dnsmadeeasy.com) using their REST
API.

## Installation

Statically-linked binaries are available for download at
https://github.com/jswank/dnsme/releases

The tool is written in Go (http://golang.org). To
compile, the Go tool chain must be installed- see
http://golang.org/doc/install for installation instructions, and
http://golang.org/doc/articles/go_command.html for info on the `go`
command.

Once the go tool chain is installed and environment configured,
compilation steps are:

	$ go get github.com/jswank/dnsme
	$ go install github.com/jswank/dnsme

This will install the `dnsme` command as $GOPATH/bin/dnsme

## Usage

	$ export DNSME_API_URL=http://api.dnsmadeeasy.com/V1.2
	$ export DNSME_API_KEY=8j7dn64b-83jc-48jd-0913-98wrhjd601df
	$ export DNSME_SECRET_KEY=93jsmq86-11hs-00ls-tnd8-8djdnb98a74c
	$ dnsme --help

	usage: dnsme command [arguments]

	'dnsme' is a command-line interface to the DNS Made Easy REST API.
	In order to use it, an API key pair must be configured for your account:
	this can be trivially created by visiting your account information page
	at https://cp.dnsmadeeasy.com/account/info.

	The following environment variables should be set:

		DNSME_API_URL = http://api.dnsmadeeasy.com/V1.2
		DNSME_API_KEY = API key
		DNSME_SECRET_KEY = Secret key

	Available commands are:

		domains          lists all domains
		domain           returns information about a domain
		add-domain       adds a domains
		delete-domain    deletes a domain
		secondaries      lists all secondary domains
		secondary        returns information about a secondary domain
		add-secondary    adds a secondary domains
		delete-secondary deletes a secondary domain
		records          return records in a domain
		record           returns a specific record id from a domain
		add-record       add a new record
		update-record    update an existing record
		delete-record    delete a record from the domain
		import           import domain info & records from JSON
		export           export domain info & records into JSON 

	Use "dnsme help [command]" for more information about a command.

	Global flags:

	The -d flag can be used to print raw HTTP requests and responses to
	stderr.

	The flag "-o" specifies the output type.  Available output types are
	"csv", "json", or the default text-based "std".

## Examples

### List primary domains

	$ ./dnsme domains
	example.com
	example.org

### List secondary domains

	$ ./dnsme secondaries
	example.net

### Show all records in a zone

	$ ./dnsme records example.com

	@     1800   A     92.250.168.100                  ; id=7693172, gtd=DEFAULT
	dev   1800   A     92.250.168.91                   ; id=7700700, gtd=DEFAULT
	dev2  1800   A     98.74.181.1                     ; id=7700747, gtd=DEFAULT
	www   1800   CNAME example.com.                    ; id=7693173, gtd=DEFAULT
	m     1800   CNAME example.com.                    ; id=7786353, gtd=DEFAULT
	admin 1800   CNAME example.com.                    ; id=7786354, gtd=DEFAULT
	@     1800   MX    10 mailstore1.secureserver.net. ; id=7693175, gtd=DEFAULT
	@     1800   MX    0 smtp.secureserver.net.        ; id=7693176, gtd=DEFAULT

### Update a record

	$ ./dnsme record -id 7693175 -o json example.com
	{"name":"","id":7693175,"type":"MX","data":"10 mailstore1.secureserver.net.",
	"gtdLocation":"DEFAULT","ttl":1800,"password":""}

## Todo

* Support for HTTP-RED records
* Smarter domain import
