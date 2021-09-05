package gen

import (
	"fmt"
	"os"

	"github.com/rwlist/gjrpc/pkg/gen/argparse"
)

func FromCmdline(args *argparse.CliArgs) {
	switch args.Command {
	case "gen:server:router":
		err := cmdServerRouter(args.Args)
		if err != nil {
			fmt.Println("Error happened,", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command \"%s\"\n", args.Command)
		os.Exit(1)
	}
}

func cmdServerRouter(args map[string]string) error {
	protoPkg := extractByKey(args, "protoPkg")
	handlersStruct := extractByKey(args, "handlersStruct")
	out := extractByKey(args, "out")

	src, err := generateServerRouter(&genServerRouterArgs{
		protoPkg:       protoPkg,
		handlersStruct: handlersStruct,
	})
	if err != nil {
		return err
	}

	return renderToFile(src, out)
}
