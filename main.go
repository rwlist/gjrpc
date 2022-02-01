package main

import (
	"fmt"
	"os"

	"github.com/rwlist/gjrpc/internal/gen"
	"github.com/rwlist/gjrpc/internal/gen/argparse"
)

func main() {
	args, err := argparse.ParseCliArgs(os.Args[1:])
	if err != nil {
		fmt.Printf("Failed to parse args, %+v", err)
		return
	}

	gen.FromCmdline(args)
}
