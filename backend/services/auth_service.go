package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/MR-DHRUV/snake_and_ladders/config"
	"github.com/MR-DHRUV/snake_and_ladders/model"
	"github.com/MR-DHRUV/snake_and_ladders/repository"

	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func GoogleAuthService(code, codeVerifier string) (*model.User, error) {
	userInfo, err := verifyGoogleParams(code, codeVerifier)
	if err != nil {
		return nil, err
	}

	// if user already exists, return the user
	existingUser, err := repository.GetUserByEmail(userInfo.Email)
	if err != nil && existingUser != nil {
		return existingUser, nil
	}

	user := &model.User{
		Email:   userInfo.Email,
		Name:    userInfo.Name,
		Picture: fetchImageAsBase64(userInfo.Picture),
	}

	user, err = repository.SaveUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func verifyGoogleParams(code, codeVerifier string) (*GoogleUserInfo, error) {
	ctx := context.Background()

	conf := &oauth2.Config{
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		RedirectURL:  config.GoogleRedirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	token, err := conf.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		return nil, errors.New("failed to exchange code for token: " + err.Error())
	}

	client := conf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, errors.New("failed to get user info: " + err.Error())
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, errors.New("failed to decode user info: " + err.Error())
	}

	if userInfo.Email == "" || !userInfo.EmailVerified {
		return nil, errors.New("unauthorized: email not verified or missing")
	}

	return &userInfo, nil
}

func GetUserById(userId string) (*model.User, error) {
	user, err := repository.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// FetchImageAsBase64 downloads an image from a URL and returns it as a base64-encoded string.
func fetchImageAsBase64(imageURL string) (string) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return imageURL
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return imageURL
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return imageURL
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg" // default fallback
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)
}
