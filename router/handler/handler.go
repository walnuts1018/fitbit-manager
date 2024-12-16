package handler

import (
	"github.com/walnuts1018/fitbit-manager/config"
	"github.com/walnuts1018/fitbit-manager/usecase"
)

type Handler struct {
	usecase       *usecase.Usecase
	defaultUserID string
}

func NewHandler(defaultUserID config.UserID, usecase *usecase.Usecase) (Handler, error) {
	return Handler{
		usecase,
		string(defaultUserID),
	}, nil
}
