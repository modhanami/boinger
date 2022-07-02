package endpoints

type UserTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

func NewUserTokenResponse(accessToken string) UserTokenResponse {
	return UserTokenResponse{AccessToken: accessToken}
}
