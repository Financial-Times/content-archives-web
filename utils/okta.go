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

// VerifyTokens function, returns whether token is valid or not
func VerifyTokens(idToken string, accessToken string, nonce string, oktaClientID string, issuer string) (bool, error) {
	tv := map[string]string{}
	tv["nonce"] = nonce
	tv["aud"] = oktaClientID
	jv := verifier.JwtVerifier{
		Issuer:           issuer,
		ClaimsToValidate: tv,
	}

	_, err := jv.New().VerifyIdToken(idToken)

	if err != nil {
		return false, fmt.Errorf("%s", err)
	}

	tv["aud"] = "api://sso"
	tv["cid"] = oktaClientID
	jv = verifier.JwtVerifier{
		Issuer:           issuer,
		ClaimsToValidate: tv,
	}

	_, err = jv.New().VerifyAccessToken(accessToken)

	if err != nil {
		return false, fmt.Errorf("%s", err)
	}

	return true, nil
}
