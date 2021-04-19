package server

import (
	"context"
	"github.com/cockroachdb/errors/grpc/status"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"miimosa-test/internal/config"
	"miimosa-test/pkg/sessions"
	"time"
)

type SessionServer struct {
	cfg    config.App
	logger *logrus.Entry

	sessions.UnimplementedSessionsServer
}

func NewSessionServer(cfg config.App) *SessionServer {
	return &SessionServer{cfg: cfg, logger: logrus.WithField("server", "sessions")}
}

func (s SessionServer) Create(_ context.Context, request *sessions.CreateRequest) (*sessions.TokenReply, error) {
	var (
		log = s.logger.WithField("operation", "create")
		now = time.Now()
	)

	_, err := uuid.Parse(request.UserId)
	if err != nil {
		err = status.WrapErr(codes.Internal, "unable to parse user id", err)

		log.
			WithField("duration-ns", time.Since(now).Nanoseconds()).
			WithError(err).
			Error()

		return nil, err
	}

	claims := jwt.StandardClaims{
		Audience:  "user",
		ExpiresAt: time.Now().Add(s.cfg.JWTExpiration).Unix(),
		Id:        uuid.New().String(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "session-server",
		NotBefore: time.Now().Unix(),
		Subject:   request.UserId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenAsString, err := token.SignedString([]byte(s.cfg.JWTSecret))

	if err != nil {
		err = status.WrapErr(codes.Internal, "unable to sign jwt", err)

		log.
			WithField("jwt-id", claims.Id).
			WithField("duration-ns", time.Since(now).Nanoseconds()).
			WithError(err).
			Error()

		return nil, err
	}

	log.
		WithField("jwt-id", claims.Id).
		WithField("duration-ns", time.Since(now).Nanoseconds()).
		Info("created json web token")

	return &sessions.TokenReply{
		Valid:     true,
		Token:     tokenAsString,
		ExpiredAt: claims.ExpiresAt,
		IssuedAt:  claims.IssuedAt,
	}, nil
}

func (s SessionServer) Validate(_ context.Context, request *sessions.ValidateTokenRequest) (*sessions.TokenReply, error) {
	var (
		log = s.logger.WithField("operation", "validate")
		now = time.Now()
	)

	claims := new(jwt.StandardClaims)
	tok, err := jwt.ParseWithClaims(request.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		err = status.WrapErr(codes.InvalidArgument, "unable to parse jwt", err)

		if tok.Claims != nil {
			log = log.WithField("jwt-id", claims.Id).WithField("jwt-valid", tok.Valid)
		}

		log.
			WithField("duration-ns", time.Since(now).Nanoseconds()).
			WithError(err).
			Error()
	} else {
		log.
			WithField("jwt-id", claims.Id).
			WithField("jwt-valid", tok.Valid).
			WithField("duration-ns", time.Since(now).Nanoseconds()).
			Info("validate json web token")
	}

	var reply = &sessions.TokenReply{
		Valid: tok.Valid,
		Token: tok.Raw,
	}

	if tok.Claims != nil {
		reply.ExpiredAt = claims.ExpiresAt
		reply.IssuedAt = claims.IssuedAt
	}

	return reply, err
}
