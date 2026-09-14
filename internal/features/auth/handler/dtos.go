package handler

import "go-marketplace/internal/core/domain"

type Response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func toDTO(model domain.TokenPair) Response {
	return Response{
		AccessToken:  model.AccessToken,
		RefreshToken: model.RefreshToken,
	}
}
