package tokenglobal

import (
	"fmt"
	"nrs_customer_module_backend/internal/global"
	"nrs_customer_module_backend/internal/types"
	"time"

	"golang.org/x/oauth2"
)

func GenerateToken(tokenValiditySeconds int64) (response *types.TokenResponse, err error) {

	// Generate a unique access token
	accessToken := global.GenerateAccessToken()

	// Generate a unique refresh token
	refreshToken := global.GenerateRefreshToken()

	tokenValidityDuration := time.Second * time.Duration(tokenValiditySeconds)
	// Calculate token expiration time
	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return nil, err
	}

	expirationTime := currentTime.Add(tokenValidityDuration)

	// Assuming authentication is successful, generate a token
	token := oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		RefreshToken: refreshToken,
		Expiry:       expirationTime,
	}

	// Serialize token response
	response = &types.TokenResponse{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		ExpiresIn:    int(tokenValiditySeconds), // Using time.Until
		RefreshToken: token.RefreshToken,
		ExpiredAt:    fmt.Sprint(expirationTime.Format(global.DefaultTimeFormat)),
	}

	return response, nil
}

func RefreshAccessToken(tokenValiditySeconds int64, refreshToken string) (response *types.TokenResponse, err error) {

	// Generate a unique access token
	accessToken := global.GenerateAccessToken()

	tokenValidityDuration := time.Second * time.Duration(tokenValiditySeconds)
	// Calculate token expiration time
	currentTime, err := global.TimeInSingapore()
	if err != nil {
		return nil, err
	}
	expirationTime := currentTime.Add(tokenValidityDuration)

	// Assuming authentication is successful, generate a token
	token := oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		RefreshToken: refreshToken,
		Expiry:       expirationTime,
	}

	// Serialize token response
	response = &types.TokenResponse{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		ExpiresIn:    int(tokenValiditySeconds), // Using time.Until
		RefreshToken: token.RefreshToken,
		ExpiredAt:    fmt.Sprint(expirationTime.Format(global.DefaultTimeFormat)),
	}

	return response, nil
}
