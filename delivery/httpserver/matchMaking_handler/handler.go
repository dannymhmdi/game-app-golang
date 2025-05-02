package matchMaking_handler

import (
	"github.com/labstack/echo/v4"
	"mymodule/params"
	"mymodule/pkg/richerr"
	"mymodule/service/authService"
	"mymodule/service/matchmakingService"
	"mymodule/service/presenceService"
	"mymodule/validator/matchMakingValidator"
	"net/http"
)

type Handler struct {
	matchMakingSvc       matchmakingService.Service
	presenceSvc          presenceService.Service
	authSvc              authService.Service
	signingKey           []byte
	matchMakingValidator matchMakingValidator.Validator
}

func New(matchMakingSvc matchmakingService.Service, authSvc authService.Service, signingKey []byte, validator matchMakingValidator.Validator) *Handler {
	return &Handler{
		matchMakingSvc: matchMakingSvc,
		//presenceSvc:          presenceSvc,
		authSvc:              authSvc,
		signingKey:           signingKey,
		matchMakingValidator: validator,
	}
}

func (h Handler) MatchMakingHandler(c echo.Context) error {
	claim := c.Get("claim").(*authService.CustomClaims)
	var bd params.AddToWaitingListRequest
	if bErr := c.Bind(&bd); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": bErr.Error()})
	}

	bd.UserId = claim.UserId
	vErr := h.matchMakingValidator.ValidateMatchMakingCredentials(bd)
	if vErr != nil {
		code, msg, op := richerr.CheckTypeErr(vErr)
		return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"error": msg, "operation": op})
	}

	res, aErr := h.matchMakingSvc.AddToWaitingList(bd)
	if aErr != nil {
		code, msg, op := richerr.CheckTypeErr(aErr)
		return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"error": msg, "operation": op})
	}

	//presenceParams:=params.PresenseRequest{
	//	UserId: claim.UserId,
	//}
	//
	//_,pErr:=h.presenceSvc.Presence(context.Background(),presenceParams)
	//if pErr != nil {
	//	code , msg,op:= richerr.CheckTypeErr(pErr)
	//	return echo.NewHTTPError(richerr.MapKindToHttpErr(code), echo.Map{"error": msg,"operation": op})
	//}
	//
	return c.JSON(http.StatusOK, res)
}
