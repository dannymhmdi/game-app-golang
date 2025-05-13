package mysqlMatchStore

import (
	"context"
	"encoding/json"
	"fmt"
	"mymodule/entity"
	"mymodule/pkg/richerr"
)

func (d DB) StoreMatch(ctx context.Context, game entity.Game) (uint, error) {

	jsonIds, mErr := json.Marshal(game.PlayersID)
	if mErr != nil {
		return 0, mErr
	}

	result, eErr := d.conn.NewConn().Exec(`INSERT INTO matches(player_ids,category)VALUES(?,?)`, string(jsonIds), game.Category)
	if eErr != nil {
		fmt.Println("debug cc")
		return 0, richerr.New().
			SetMsg(eErr.Error()).
			SetKind(richerr.KindUnexpected).
			SetWrappedErr(eErr).
			SetOperation("mysqlMatchStore.StoreMatch")
	}

	id, lErr := result.LastInsertId()
	if lErr != nil {
		return 0, richerr.New().
			SetMsg(lErr.Error()).
			SetKind(richerr.KindUnexpected).
			SetWrappedErr(lErr).
			SetOperation("mysqlMatchStore.StoreMatch")
	}

	return uint(id), nil
}
