package gen

import (
	"github.com/rwlist/gjrpc/internal/gen/protog"
	"github.com/rwlist/gjrpc/internal/gen/tsproto"
)

type genTypescriptProtoArgs struct {
	protoPkg string
}

func generateTypescriptProto(args *genTypescriptProtoArgs) (string, error) {
	proto, err := protog.Parse(args.protoPkg)
	if err != nil {
		return "", err
	}

	return tsproto.GenerateSource(proto)
}
