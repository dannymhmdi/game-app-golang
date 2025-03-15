package slice

import (
	"mymodule/contract/golang/presence"
	"mymodule/entity"
)

func Uint64ToUintMapper(list []uint64) []uint {
	s := make([]uint, len(list))
	for i, id := range list {
		s[i] = uint(id)
	}
	return s
}

func UintToUint64Mapper(list []uint) []uint64 {
	s := make([]uint64, len(list))
	for i, id := range list {
		s[i] = uint64(id)
	}
	return s
}

func OnlinePlayerMapperToProtobuf(list []entity.OnlinePlayer) []*presence.OnlinePlayer {
	s := make([]*presence.OnlinePlayer, len(list))

	for i, player := range list {
		s[i] = &presence.OnlinePlayer{
			UserId:    uint64(player.UserId),
			Timestamp: player.Timestamp,
		}
	}

	return s

}

func OnlinePlayerMapperToParams(list []*presence.OnlinePlayer) []entity.OnlinePlayer {
	s := make([]entity.OnlinePlayer, len(list))
	for i, player := range list {
		s[i] = entity.OnlinePlayer{
			UserId:    uint(player.UserId),
			Timestamp: player.Timestamp,
		}
	}

	return s

}
