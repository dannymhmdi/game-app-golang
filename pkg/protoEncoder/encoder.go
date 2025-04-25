package protoEncoder

import (
	"encoding/base64"
	"google.golang.org/protobuf/proto"
	"log"
	"mymodule/contract/golang/matchingPlayer"
	"mymodule/entity"
	"mymodule/pkg/slice"
)

func Encoder(msg any) string {
	switch msg.(type) {
	case entity.MatchedPlayers:
		mu := msg.(entity.MatchedPlayers)
		protoMu := matchingPlayer.MatchedPlayers{
			UserIds:   slice.UintToUint64Mapper(mu.UserIDs),
			Category:  string(mu.Category),
			Timestamp: mu.Timestamp,
		}

		return encodeProtoToString(&protoMu)

	default:
		return ""
	}
}

func encodeProtoToString(m proto.Message) string {
	protoMsg, mErr := proto.Marshal(m)
	if mErr != nil {
		log.Fatalf("failed to encode to proto message: %v", mErr)
	}
	return base64.StdEncoding.EncodeToString(protoMsg)
}
