package okta

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	verifier "github.com/okta/okta-jwt-verifier-golang"
)

// GenerateNonce returns a value that is returned in the ID token. It is used to mitigate replay attacks.
func GenerateNonce() (string, error) {
	nonceBytes := make([]byte, 32)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("could not generate nonce")
	}

	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}

// VerifyToken function, returns whether token is valid or not
func VerifyToken(t string, nonce string, oktaClientID string, issuer string) (*verifier.Jwt, error) {
	tv := map[string]string{}
	tv["nonce"] = nonce
	tv["aud"] = oktaClientID
	jv := verifier.JwtVerifier{
		Issuer:           issuer,
		ClaimsToValidate: tv,
	}

	result, err := jv.New().VerifyIdToken(t)

	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if result != nil {
		return result, nil
	}

	return nil, fmt.Errorf("token could not be verified: %s", "")
}
