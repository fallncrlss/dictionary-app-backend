package controller

import (
	"net/http"
	"strings"

	"github.com/fallncrlss/dictionary-app-backend/src/internal/lib/customerrors"
	"github.com/fallncrlss/dictionary-app-backend/src/internal/lib/enums"
	"github.com/fallncrlss/dictionary-app-backend/src/pkg/service"
	echoLog "github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type WordController interface {
	Get(ctx echo.Context) error
	Search(ctx echo.Context) error
}

type wordController struct {
	services *service.Manager
	logger   *echo.Logger
}

func NewWordControllers(services *service.Manager, logger *echo.Logger) WordController {
	return &wordController{
		services: services,
		logger:   logger,
	}
}

func (wc *wordController) Get(ctx echo.Context) error {
	languageCode := ctx.Param("language")
	name := strings.ToLower(ctx.Param("name"))
	language, isCorrect := enums.GetLanguageCodes()[languageCode]

	if !isCorrect {
		return errors.Wrapf(customerrors.ErrIncorrectArgument, "arguments (language='%s')", language)
	}

	echoLog.Debugf("Getting word '%s' (language='%s') from database...", name, language)

	wordData, err := wc.services.Word.GetWordWithDB(name, language)
	if err != nil {
		switch {
		case errors.Is(err, customerrors.ErrUnableFetchInstance), errors.Is(err, customerrors.ErrFetchedInstanceIsNil):
			echoLog.Debugf("Unable to get word from DB, getting word '%s' from web... Error: %s", name, err)

			wordData, err = wc.services.Word.GetWordWithWeb(name, languageCode)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			echoLog.Debugf("Saving word '%s' into DB...", name)

			err := wc.services.Word.SaveWordToDB(wordData)
			if err != nil {
				echoLog.Warn(err)
			}

		case errors.Is(err, customerrors.ErrIncorrectArgument):
			return echo.NewHTTPError(http.StatusBadRequest, err)

		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "Could not get word"))
		}
	}

	echoLog.Debugf("Return successful response: %s", wordData)

	err = ctx.JSON(http.StatusOK, wordData)
	if err != nil {
		return errors.Wrap(err, "sending a JSON response with status code failed")
	}

	return nil
}

func (wc *wordController) Search(ctx echo.Context) error {
	searchInput := ctx.QueryParam("q")

	result, err := wc.services.Word.SearchByInput(searchInput)
	if err != nil {
		return err
	}

	err = ctx.JSON(http.StatusOK, result)
	if err != nil {
		return errors.Wrap(err, "sending a JSON response with status code failed")
	}

	echoLog.Debugf("Return successful response: %s", result)
	return nil
}
