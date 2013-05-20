package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	// Production Values
	API_URL    = "http://api.dnsmadeeasy.com/V1.2"

	// Development Values
	//API_URL = "http://api.sandbox.dnsmadeeasy.com/V1.2"
)

var (
	api_url    string
	api_key    string
	secret_key string

	outputType        string
	debug             bool
	requestsRemaining int
)

var commands = []*Command{
	listDomains,
	infoDomain,
	addNewDomain,
	delDomain,
	listSecondaries,
	infoSecondary,
	addNewSecondary,
	delSecondary,
	records,
	record,
	addRecord,
	updateRecord,
	deleteRecord,
	importData,
	exportData,
	/*
		addRecord,
		search, */
}

// A Command is an implementation of a go command
// like go build or go fix.
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string) error

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'dnsme help' output.
	Short string

	// Long is the long message shown in the 'dnsme help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command will do its own
	// flag parsing.
	CustomFlags func(cmd *flag.FlagSet)
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}

func main() {

	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
		return
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	api_url = os.Getenv("DNSME_API_URL")
	if api_url == "" {
		api_url = API_URL
	}

	api_key = os.Getenv("DNSME_API_KEY")
	if api_key == "" {
		fmt.Fprint(os.Stderr, "DNSME_API_KEY environment variable is not set\n")
		os.Exit(1)
	}

	secret_key = os.Getenv("DNSME_SECRET_KEY")
	if secret_key == "" {
		fmt.Fprint(os.Stderr, "DNSME_SECRET_KEY environment variable is not set\n")
		os.Exit(1)
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			addGlobalFlags(&cmd.Flag)
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags != nil {
				cmd.CustomFlags(&cmd.Flag)
			}
			//				args = args[1:]
			//			} else {
			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()
			//			}
			err := cmd.Run(cmd, args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command %#q\n\n", args[0])
	usage()

}

func addGlobalFlags(fs *flag.FlagSet) {
	fs.StringVar(&outputType, "o", "std", "Output type (std, json, csv)")
	fs.BoolVar(&debug, "d", false, "Debug output")
}

func printUsage(w io.Writer) {
	tmpl(w, usageTemplate, commands)
}

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

// help implements the 'help' command.
func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		// not exit 2: succeeded at 'dnsme help'.
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: dnsme help command\n\nToo many arguments given.\n")
		os.Exit(2) // failed at 'dnsme help'
	}

	arg := args[0]
	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			// not exit 2: succeeded at 'dnsme help cmd'.
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q.  Run 'dnsme help'.\n", arg)
	os.Exit(2) // failed at 'dnsme help cmd'
}

/*
func printOutput(r *http.Response) {
	io.Copy(os.Stdout, r.Body)
	fmt.Println()
}
*/
