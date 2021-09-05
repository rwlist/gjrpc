package main

import (
	"fmt"
	"os"

	"github.com/rwlist/gjrpc/pkg/gen"
	"github.com/rwlist/gjrpc/pkg/gen/argparse"
)

func main() {
	args, err := argparse.ParseCliArgs(os.Args[1:])
	if err != nil {
		fmt.Println("Failed to parse args,", err)
		return
	}

	gen.FromCmdline(args)
}
