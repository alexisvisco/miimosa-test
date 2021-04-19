package server

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"miimosa-test/internal/config"
	"miimosa-test/pkg/sessions"
	"testing"
	"time"
)

var secret = []byte("hello world")

var sessionServer = NewSessionServer(config.App{
	JWTSecret:     string(secret),
	JWTExpiration: time.Second * 1,
})

func TestSessionServer_Create(t *testing.T) {
	token := mustCreateToken(t)

	claims := new(jwt.StandardClaims)

	tok, err := jwt.ParseWithClaims(token.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if !tok.Valid {
		t.Fatal("issued token must be valid")
	}
}

func TestSessionServer_Validate(t *testing.T) {

	var (
		validToken   = mustCreateToken(t).Token
		expiredToken = mustCreateTokenThatExpire(t)
	)

	tests := []struct {
		name    string
		token   string
		wantErr bool
		test    func(reply *sessions.TokenReply)
	}{
		{
			name:    "invalid character",
			token:   "lalala.lalala.lalala",
			wantErr: true,
			test: func(r *sessions.TokenReply) {
				if r.Valid {
					t.Fatal("token must be invalid")
				}
			},
		},
		{
			name:    "invalid signature token",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			wantErr: true,
			test: func(r *sessions.TokenReply) {
				if r.Valid {
					t.Fatal("token must be invalid")
				}
			},
		},
		{
			name:    "token expired",
			token:   expiredToken,
			wantErr: true,
			test: func(r *sessions.TokenReply) {
				if r.Valid {
					t.Fatal("token must be invalid")
				}
			},
		},
		{
			name:    "valid token",
			token:   validToken,
			wantErr: false,
			test: func(r *sessions.TokenReply) {
				if !r.Valid {
					t.Fatal("token not valid")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sessionServer.Validate(context.Background(), &sessions.ValidateTokenRequest{
				Token: tt.token,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.test != nil {
				tt.test(got)
			}
		})
	}
}

func mustCreateToken(t *testing.T) *sessions.TokenReply {
	token, err := sessionServer.Create(context.Background(), &sessions.CreateRequest{UserId: "123"})
	if err != nil {
		t.Fatal(err)
	}

	return token
}

func mustCreateTokenThatExpire(t *testing.T) string {
	claims := jwt.StandardClaims{
		Audience:  "user",
		ExpiresAt: time.Now().Add(-(time.Second * 5)).Unix(),
		Id:        uuid.New().String(),
		IssuedAt:  time.Now().Add(-(time.Second * 10)).Unix(),
		Issuer:    "session-server",
		NotBefore: time.Now().Add(-(time.Second * 10)).Unix(),
		Subject:   "123",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	str, err := token.SignedString(secret)
	if err != nil {
		t.Fatal(err)
		return ""
	}

	return str
}
