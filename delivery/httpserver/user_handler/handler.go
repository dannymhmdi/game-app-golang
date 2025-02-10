package user_handler

import (
	"github.com/labstack/echo/v4"
	"mymodule/dto"
	"mymodule/pkg/richerr"
	"mymodule/service/authservice"
	"mymodule/service/userservice"
	"mymodule/validator/uservalidator"
	"net/http"
)

//	type AuthGenerator interface {
//		ParseToken(tokenString string) (*authservice.CustomClaims, error)
//	}
type Handler struct {
	authSvc       authservice.Service
	userSvc       userservice.Service
	userValidator uservalidator.Validator
	authSignKey   []byte
}

func New(authSvc authservice.Service, userSvc userservice.Service, validator uservalidator.Validator, signKey []byte) *Handler {
	return &Handler{
		authSvc:       authSvc,
		userSvc:       userSvc,
		userValidator: validator,
		authSignKey:   signKey,
	}
}

func (h Handler) userRegisterHandler(c echo.Context) error {
	bd := dto.RegisterRequest{}
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

	bd := dto.LoginRequest{}
	if bErr := c.Bind(&bd); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, bErr.Error())
	}

	vErr := h.userValidator.ValidateLoginCredentials(bd)
	if vErr != nil {
		code, msg, op := richerr.CheckTypeErr(vErr)
		return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"message": msg, "operation": op})
	}

	loginResp, lErr := h.userSvc.Login(bd)
	if lErr != nil {
		code, msg, op := richerr.CheckTypeErr(lErr)
		return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"message": msg, "operation": op})
	}
	return c.JSON(http.StatusOK, loginResp)
}

func (h Handler) userProfileHandler(c echo.Context) error {
	//token := c.Request().Header.Get("Authorization")
	//if token == "" {
	//	return echo.NewHTTPError(http.StatusUnauthorized, "token is empty")
	//}
	//
	//claim, pErr := h.authSvc.ParseToken(token)
	//if pErr != nil {
	//	code, msg, op := richerr.CheckTypeErr(pErr)
	//	return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"message": msg, "operation": op})
	//}
	claim := c.Get("claim")
	customClaims, ok := claim.(*authservice.CustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"message": "token is not valid"})
	}
	userInfo, gErr := h.userSvc.GetUserProfile(customClaims.UserId)
	if gErr != nil {
		code, msg, op := richerr.CheckTypeErr(gErr)
		return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"message": msg, "operation": op})
	}

	return c.JSON(http.StatusOK, userInfo)
}
