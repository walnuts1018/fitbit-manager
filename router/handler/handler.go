package handler

import "github.com/walnuts1018/fitbit-manager/usecase"

type Handler struct {
	usecase *usecase.Usecase
}

func NewHandler(usecase *usecase.Usecase) (Handler, error) {
	return Handler{
		usecase,
	}, nil
}
