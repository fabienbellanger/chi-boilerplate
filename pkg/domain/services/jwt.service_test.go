package services

import (
	"chi_boilerplate/pkg/domain/entities"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	type args struct {
		user     entities.User
		lifetime time.Duration
		algo     string
		secret   string
	}

	type result struct {
		jwt JWT
		err error
	}

	lifetime := time.Duration(2)

	tests := []struct {
		name   string
		args   args
		wanted result
	}{
		{
			name: "Invalid algo",
			args: args{
				user:     entities.User{},
				lifetime: lifetime,
				algo:     "",
				secret:   "my-secret",
			},
			wanted: result{
				jwt: JWT{},
				err: errors.New("unsupported JWT algo: must be HS512 or ES384"),
			},
		},
		{
			name: "Invalid algo",
			args: args{
				user:     entities.User{},
				lifetime: lifetime,
				algo:     "HS512",
				secret:   "secret",
			},
			wanted: result{
				jwt: JWT{},
				err: errors.New("secret must have at least 8 characters"),
			},
		},
		{
			name: "Valid",
			args: args{
				user:     entities.User{},
				lifetime: lifetime,
				algo:     "HS512",
				secret:   "my-secret",
			},
			wanted: result{
				jwt: JWT{},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwt, err := NewJWT(
				tt.args.user.ID,
				tt.args.lifetime,
				tt.args.algo,
				tt.args.secret,
			)
			got := result{jwt, err}

			if got.err != nil {
				assert.Equal(t, got.jwt.Value, tt.wanted.jwt.Value)
			} else {
				assert.Greater(t, len(got.jwt.Value), 0)
				assert.Greater(t, got.jwt.ExpiredAt, time.Now().Add(lifetime*time.Hour-time.Minute))
				assert.Less(t, got.jwt.ExpiredAt, time.Now().Add(lifetime*time.Hour+time.Minute))
			}
			assert.Equal(t, got.err, tt.wanted.err)
		})
	}
}
