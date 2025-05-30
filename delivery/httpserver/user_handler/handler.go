package user_handler

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"mymodule/config"
	"mymodule/params"
	"mymodule/pkg/richerr"
	"mymodule/pkg/types"
	"mymodule/service/authService"
	"mymodule/service/presenceService"
	"mymodule/service/userService"
	"mymodule/validator/uservalidator"
	"net/http"
	"time"
)

type contextKey string

const Key contextKey = "userAgent"

type Handler struct {
	authSvc       authService.Service
	userSvc       userService.Service
	presenceSvc   presenceService.Service
	userValidator uservalidator.Validator
	authSignKey   []byte
}

func New(authSvc authService.Service, userSvc userService.Service, presenceSvc presenceService.Service, validator uservalidator.Validator, signKey []byte) *Handler {
	return &Handler{
		authSvc:       authSvc,
		userSvc:       userSvc,
		presenceSvc:   presenceSvc,
		userValidator: validator,
		authSignKey:   signKey,
	}
}

func (h Handler) userRegisterHandler(c echo.Context) error {
	bd := params.RegisterRequest{}
	if bErr := c.Bind(&bd); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, bErr.Error())
	}

	vErr := h.userValidator.ValidateRegisterCredentials(bd)
	if vErr != nil {
		code, msg, op := richerr.CheckTypeErr(vErr)
		return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"message": msg, "operation": op})
	}

	createdUSer, rErr := h.userSvc.Register(bd)
	if rErr != nil {
		code, msg, op := richerr.CheckTypeErr(rErr)
		return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"message": msg, "operation": op})
	}

	return c.JSON(http.StatusOK, createdUSer)

}

func (h Handler) userLoginHandler(c echo.Context) error {
	bd := params.LoginRequest{}
	parentCtx := context.WithValue(context.Background(), types.Key, c.Request().UserAgent())
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*5)
	defer cancel()
	if bErr := c.Bind(&bd); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, bErr.Error())
	}

	vErr := h.userValidator.ValidateLoginCredentials(bd)
	if vErr != nil {
		code, msg, op := richerr.CheckTypeErr(vErr)
		return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"message": msg, "operation": op})
	}

	cookie, cErr := c.Request().Cookie("refresh-token")
	if cErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"message": cErr.Error(), "op": "user_hanler.userLoginHandler"})
	}

	fmt.Println("cookie", cookie)

	loginResp, lErr := h.userSvc.Login(ctx, bd)
	if lErr != nil {
		code, msg, op := richerr.CheckTypeErr(lErr)
		return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"message": msg, "operation": op})
	}

	appConfig := config.Load()
	c.SetCookie(&http.Cookie{
		Name:     "refresh-token",
		Value:    loginResp.Token.RefreshToken,
		Expires:  time.Now().Add(appConfig.AuthConfig.RefreshTokenExpireTime),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, loginResp)
}

func (h Handler) userProfileHandler(c echo.Context) error {
	//ToDO write a funtion to get claims: getClaims
	claim := c.Get("claim")
	customClaims, ok := claim.(*authService.CustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"message": "token is not valid"})
	}
	userInfo, gErr := h.userSvc.GetUserProfile(customClaims.UserId)
	if gErr != nil {
		code, msg, op := richerr.CheckTypeErr(gErr)
		return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"message": msg, "operation": op})
	}

	generatedAccessToken, isOk := c.Get("generatedNewAccessToken").(string)

	if isOk && generatedAccessToken != "" {
		userInfo.RegeneratedToken = generatedAccessToken
		return c.JSON(http.StatusOK, userInfo)
	}

	return c.JSON(http.StatusOK, userInfo)
}
