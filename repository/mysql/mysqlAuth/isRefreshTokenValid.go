package mysqlAuth

import (
	"crypto/sha256"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"mymodule/entity"
	"mymodule/logger"
	"mymodule/pkg/richerr"
	"strings"
)

func (d DB) IsRefreshTokenValid(cookieRefreshToken string, userId uint) error {
	cookieRefreshToken = strings.Replace(cookieRefreshToken, "refresh-token=", "", 1)
	tokenHash := sha256.Sum256([]byte(cookieRefreshToken))
	rows, qErr := d.conn.NewConn().Query("select * from refresh_tokens where user_id = ?", userId)
	if qErr != nil {
		return richerr.New().
			SetMsg(qErr.Error()).
			SetKind(richerr.KindUnexpected).
			SetOperation("mysqlAuth.IsRefreshTokenValid").
			SetWrappedErr(qErr)
	}

	var cErr error
	var sErr error
	for rows.Next() {
		refToken := entity.RefreshToken{}
		//var createdAt []uint8
		if sErr = rows.Scan(&refToken.Id, &refToken.UserId, &refToken.TokenHash, &refToken.ExpireAt, &refToken.Revoked, &refToken.DeviceInfo, &refToken.CreatedAt); sErr != nil {
			continue
		}
		//refToken.CreatedAt = createdAt
		if cErr = bcrypt.CompareHashAndPassword([]byte(refToken.TokenHash), tokenHash[:]); cErr != nil {
			continue
		}
		logger.Info("IsRefreshTokenValid", zap.Any("isRefreshTokenValid", refToken))
		return nil
	}

	if sErr != nil {
		fmt.Println("kiiiiir", sErr)
		return richerr.New().
			SetMsg(sErr.Error()).
			SetKind(richerr.KindUnexpected).
			SetOperation("mysqlAuth.IsRefreshTokenValid").
			SetWrappedErr(sErr)
	}

	if cErr != nil {
		fmt.Println("kiiiiir", cErr)
		return richerr.New().
			SetMsg(cErr.Error()).
			SetKind(richerr.KindUnauthorized).
			SetOperation("mysqlAuth.IsRefreshTokenValid").
			SetWrappedErr(cErr)
	}

	return nil

}
