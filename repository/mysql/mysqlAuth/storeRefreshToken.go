package mysqlAuth

import (
	"context"
	"crypto/sha256"
	"golang.org/x/crypto/bcrypt"
	"mymodule/entity"
	"mymodule/pkg/richerr"
	"mymodule/pkg/types"
	"mymodule/service/authService"
	"time"
)

func (d DB) StoreRefreshToken(ctx context.Context, refreshToken string, userInfo entity.User, cfg authService.Config) error {
	userAgent := ctx.Value(types.Key).(string)
	tokenHash := sha256.Sum256([]byte(refreshToken))
	hashedToken, gErr := bcrypt.GenerateFromPassword(tokenHash[:], bcrypt.DefaultCost)
	if gErr != nil {
		return richerr.New().
			SetMsg(gErr.Error()).
			SetOperation("mysqlAuth.StoreRefreshToken").
			SetKind(richerr.KindUnexpected).
			SetWrappedErr(gErr)
	}
	result, eErr := d.conn.NewConn().Exec("insert into refresh_tokens (user_id,token_hash,expires_at,revoked,device_info) values(?,?,?,?,?)", userInfo.ID, string(hashedToken), time.Now().Add(cfg.RefreshTokenExpireTime), 0, userAgent)
	if eErr != nil {
		return richerr.New().
			SetMsg(eErr.Error()).
			SetKind(richerr.KindUnexpected).
			SetWrappedErr(eErr).
			SetOperation("mysqlAuth.StoreRefreshToken")
	}

	_, lErr := result.LastInsertId()
	if lErr != nil {
		return richerr.New().
			SetKind(richerr.KindUnexpected).
			SetWrappedErr(lErr).
			SetOperation("mysqlAuth.StoreRefreshToken").
			SetMsg(lErr.Error())
	}

	return nil
}
