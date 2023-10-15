package vault

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
)

const percentage = 75

type AutoRenew struct {
	auth   Auth
	ticker *time.Ticker

	expiration chan time.Duration
	done       chan bool

	once   sync.Once
	client *api.Client
}

func NewAutoRenew(auth Auth, done chan bool) (*AutoRenew, error) {
	if auth == nil {
		return nil, errors.New("no auth impl given")
	}

	a := &AutoRenew{
		auth:       auth,
		expiration: make(chan time.Duration, 1),
		done:       done,
	}

	go a.fuck()

	return a, nil
}

func (t *AutoRenew) fuck() {
	t.ticker = time.NewTicker(24 * time.Hour)
	defer t.ticker.Stop()

	cont := true
	for cont {
		select {
		case <-t.done:
			cont = false
		case expiration := <-t.expiration:
			log.Info().Msgf("Scheduling auto token renewal in %v", expiration)
			if expiration <= 0 {
				log.Error().Msgf("got faulty secret expiration: %v", expiration)
				continue
			}
			t.ticker.Reset(expiration)
		case <-t.ticker.C:
			func() {
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()
				log.Info().Msg("Automatically renewing vault token")
				_, err := t.Login(ctx, t.client)
				if err != nil {
					log.Error().Err(err).Msg("encountered while trying to automatically renew token")
				}
			}()
		}
	}
}

func calculateSafeTtl(seconds int) int {
	return percentage * seconds / 100
}

func (t *AutoRenew) Logout(ctx context.Context, client *api.Client) error {
	return t.auth.Logout(ctx, client)
}

func (t *AutoRenew) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	if client == nil {
		return nil, errors.New("empty client")
	}

	t.once.Do(func() {
		t.client = client
	})

	// try login
	secret, err := t.auth.Login(ctx, client)
	if err != nil {
		return nil, err
	}

	ttl := calculateSafeTtl(secret.Auth.LeaseDuration)
	t.expiration <- time.Second * time.Duration(ttl)
	return secret, err
}
