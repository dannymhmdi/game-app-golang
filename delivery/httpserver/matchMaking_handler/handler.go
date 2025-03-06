package matchMaking_handler

import (
	"github.com/labstack/echo/v4"
	"mymodule/dto"
	"mymodule/pkg/richerr"
	"mymodule/service/authservice"
	"mymodule/service/matchmakingService"
	"mymodule/validator/matchMakingValidator"
	"net/http"
)

type Handler struct {
	matchMakingSvc       matchmakingService.Service
	authSvc              authservice.Service
	signingKey           []byte
	matchMakingValidator matchMakingValidator.Validator
}

func New(matchMakingSvc matchmakingService.Service, authSvc authservice.Service, signingKey []byte, validator matchMakingValidator.Validator) *Handler {
	return &Handler{
		matchMakingSvc:       matchMakingSvc,
		authSvc:              authSvc,
		signingKey:           signingKey,
		matchMakingValidator: validator,
	}
}

func (h Handler) MatchMakingHandler(c echo.Context) error {
	claim := c.Get("claim").(*authservice.CustomClaims)
	var bd dto.AddToWaitingListRequest
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

	return c.JSON(http.StatusOK, res)
}
