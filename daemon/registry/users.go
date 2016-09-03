package registry

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/arigatomachine/cli/envelope"
	"github.com/arigatomachine/cli/primitive"

	"github.com/arigatomachine/cli/daemon/crypto"
)

// Users represents the  registry `/users` endpoints.
type Users struct {
	client *Client
}

// GetSelf returns the logged in user.
func (u *Users) GetSelf(ctx context.Context, token string) (*envelope.Unsigned, error) {
	req, err := u.client.NewTokenRequest(token, "GET", "/users/self", nil, nil)
	if err != nil {
		log.Printf("Error making api request: %s", err)
		return nil, err
	}

	self := envelope.Unsigned{}
	_, err = u.client.Do(ctx, req, &self)
	if err != nil {
		log.Printf("Error making api request: %s", err)
		return nil, err
	}

	err = validateSelf(&self)
	if err != nil {
		log.Printf("Invalid user self: %s", err)
		return nil, err
	}

	return &self, nil
}

// Create attempts to register a new user
func (u *Users) Create(ctx context.Context, userObj User) (*envelope.Unsigned, error) {
	req, err := u.client.NewRequest("POST", "/users", nil, userObj)
	if err != nil {
		log.Printf("Error making api request: %s", err)
		return nil, err
	}

	user := envelope.Unsigned{}
	_, err = u.client.Do(ctx, req, &user)
	if err != nil {
		log.Printf("Error making api request: %s", err)
		return nil, err
	}

	err = validateSelf(&user)
	if err != nil {
		log.Printf("Invalid user object: %s", err)
		return nil, err
	}

	return &user, nil
}

func validateSelf(s *envelope.Unsigned) error {
	if s.Version != 1 {
		return errors.New("version must be 1")
	}

	body := s.Body.(*primitive.User)

	if body == nil {
		return errors.New("missing body")
	}

	if body.Master == nil {
		return errors.New("missing master key section")
	}

	if body.Master.Alg != crypto.Triplesec {
		return fmt.Errorf("Unknown alg: %s", body.Master.Alg)
	}

	if len(*body.Master.Value) == 0 {
		return errors.New("Zero length master key found")
	}

	return nil
}

// User contains fields for signup
type User struct {
	ID      string          `json:"id"`
	Version int             `json:"version"`
	Body    *primitive.User `json:"body"`
}
