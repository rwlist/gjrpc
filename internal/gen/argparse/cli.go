package argparse

import (
	"github.com/pkg/errors"
	"strings"
)

type CliArgs struct {
	Command string
	Args    map[string]string
}

func ParseCliArgs(args []string) (*CliArgs, error) {
	if len(args) == 0 {
		args = []string{"help"}
	}

	res := &CliArgs{
		Command: args[0],
		Args:    map[string]string{},
	}

	args = args[1:]
	for _, arg := range args {
		rawArg := strings.TrimPrefix(arg, "--")
		if rawArg == arg {
			return nil, errors.Errorf("arguments must start with -- but found \"%s\"", arg)
		}

		kv := strings.SplitN(rawArg, "=", 2)
		if len(kv) == 1 {
			kv = append(kv, "")
		}
		if len(kv) != 2 { //nolint:gomnd
			return nil, errors.Errorf("invalid argument found \"%s\"", arg)
		}
		key, value := kv[0], kv[1]

		if _, ok := res.Args[key]; ok {
			return nil, errors.Errorf("duplicate argument %s", key)
		}
		res.Args[key] = value
	}

	return res, nil
}
