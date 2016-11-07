package api

import (
	"context"

	"github.com/manifoldco/torus-cli/apitypes"
	"github.com/manifoldco/torus-cli/identity"
	"github.com/manifoldco/torus-cli/primitive"
)

// UsersClient makes proxied requests to the registry's users endpoints
type UsersClient struct {
	client *Client
}

// UserResult is the payload returned for a user object
type UserResult struct {
	ID      *identity.ID    `json:"id"`
	Version uint8           `json:"version"`
	Body    *primitive.User `json:"body"`
}

// Signup will have the daemon create a new user request
func (u *UsersClient) Signup(ctx context.Context, signup *apitypes.Signup, output *ProgressFunc) (*UserResult, error) {
	req, _, err := u.client.NewRequest("POST", "/signup", nil, &signup, false)
	if err != nil {
		return nil, err
	}

	user := UserResult{}
	_, err = u.client.Do(ctx, req, &user, nil, output)
	return &user, err
}

// VerifyEmail will confirm the user's email with the registry
func (u *UsersClient) VerifyEmail(ctx context.Context, verifyCode string) error {
	verify := apitypes.VerifyEmail{
		Code: verifyCode,
	}
	req, _, err := u.client.NewRequest("POST", "/users/verify", nil, &verify, true)
	if err != nil {
		return err
	}

	_, err = u.client.Do(ctx, req, nil, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

type userUpdateEmail struct {
	Email string `json:"email"`
}

// UpdateEmail updates the user profile's email field
func (u *UsersClient) UpdateEmail(ctx context.Context, email string) (*UserResult, error) {
	updateEmail := userUpdateEmail{Email: email}
	req, _, err := u.client.NewRequest("PATCH", "/users/self", nil, &updateEmail, true)
	if err != nil {
		return nil, err
	}

	user := UserResult{}
	_, err = u.client.Do(ctx, req, &user, nil, nil)
	return &user, err
}

type userUpdateName struct {
	Name string `json:"name"`
}

// UpdateName updates the user profile's name field
func (u *UsersClient) UpdateName(ctx context.Context, name string) (*UserResult, error) {
	updateName := userUpdateName{Name: name}
	req, _, err := u.client.NewRequest("PATCH", "/users/self", nil, &updateName, true)
	if err != nil {
		return nil, err
	}

	user := UserResult{}
	_, err = u.client.Do(ctx, req, &user, nil, nil)
	return &user, err
}
