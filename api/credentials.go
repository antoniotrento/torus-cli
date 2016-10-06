package api

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/arigatomachine/cli/apitypes"
)

// CredentialsClient provides access to unencrypted credentials for viewing,
// and encrypts credentials when setting.
type CredentialsClient struct {
	client *Client
}

// Get returns all credentials at the given path.
func (c *CredentialsClient) Get(ctx context.Context, path string) ([]apitypes.CredentialEnvelope, error) {
	v := &url.Values{}
	v.Set("path", path)

	req, _, err := c.client.NewRequest("GET", "/credentials", v, nil, false)
	if err != nil {
		return nil, err
	}

	resp := []apitypes.CredentialResp{}

	_, err = c.client.Do(ctx, req, &resp, nil, nil)
	if err != nil {
		return nil, err
	}

	creds := make([]apitypes.CredentialEnvelope, len(resp))
	for i, c := range resp {
		v, err := createEnvelopeFromResp(c)
		if err != nil {
			return nil, err
		}
		creds[i] = *v
	}

	return creds, err
}

// Create creates the given credential
func (c *CredentialsClient) Create(ctx context.Context, cred *apitypes.Credential,
	progress *ProgressFunc) (*apitypes.CredentialEnvelope, error) {

	env := apitypes.CredentialEnvelope{Version: 1, Body: cred}
	req, reqID, err := c.client.NewRequest("POST", "/credentials", nil, &env, false)
	if err != nil {
		return nil, err
	}

	resp := apitypes.CredentialResp{}
	_, err = c.client.Do(ctx, req, &resp, &reqID, progress)
	if err != nil {
		return nil, err
	}

	out, err := createEnvelopeFromResp(resp)
	return out, err
}

func createEnvelopeFromResp(c apitypes.CredentialResp) (*apitypes.CredentialEnvelope, error) {
	var envelope apitypes.CredentialEnvelope
	switch c.Version {
	case 1:
		var cBody apitypes.Credential
		cBodyV1 := apitypes.CredentialV1{}

		err := json.Unmarshal(c.Body, &cBodyV1)
		if err != nil {
			return nil, err
		}

		cBody = &cBodyV1
		envelope = apitypes.CredentialEnvelope{
			ID:      c.ID,
			Version: c.Version,
			Body:    &cBody,
		}
		break
	default:
		panic("Omg I don't know this version")
	}

	return &envelope, nil
}
