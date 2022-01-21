package main

import (
	"fmt"
	"github.com/rwlist/gjrpc/internal/gen"
	"github.com/rwlist/gjrpc/internal/gen/argparse"
	"os"
)

func main() {
	args, err := argparse.ParseCliArgs(os.Args[1:])
	if err != nil {
		fmt.Println("Failed to parse args,", err)
		return
	}

	gen.FromCmdline(args)
}
