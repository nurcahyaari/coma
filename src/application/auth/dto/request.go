package dto

type RequestValidateToken struct {
	AccessToken string
}

type RequestGenerateToken struct {
	// key can be username
	Key string `json:"key"`
	// secret can be user password
	Secret string `json:"secret"`
}

type ResponseGenerateToken struct {
	AccessToken     string `json:"accessToken"`
	AccessTokenExp  string `json:"accessTokenExp"`
	RefreshToken    string `json:"refreshToken"`
	RefreshTokenExp string `json:"refreshTokenExp"`
}
