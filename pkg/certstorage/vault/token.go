package vault

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"strconv"
	"time"
)

type TokenData struct {
	ExpireTime time.Time
	Renewable  bool
}

func (token *TokenData) MinutesUntilExpiry() int64 {
	return int64(token.ExpireTime.Sub(time.Now()).Minutes())
}

func (token *TokenData) PrettyExpiryDate() string {
	return fmt.Sprintf("%s (%d minutes)", token.ExpireTime.Format(time.RFC822), token.MinutesUntilExpiry())
}

func FromSecret(secret *api.Secret) *TokenData {
	if secret == nil {
		return &TokenData{}
	}

	ttl, err := strconv.Atoi(fmt.Sprintf("%v", secret.Data["ttl"]))
	var expiry time.Time
	if err == nil {
		expiry = time.Now().Add(time.Duration(ttl) * time.Second)
	}

	renewable, err := strconv.ParseBool(fmt.Sprintf("%t", secret.Data["renewable"]))
	if err != nil {
		renewable = false
	}
	return &TokenData{
		ExpireTime: expiry,
		Renewable:  renewable,
	}
}
