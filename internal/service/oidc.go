//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package service

import (
	"context"
	"net/url"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/internal/logging"
	"go.uber.org/zap"
)

type AuthProofs struct {
	State        string
	CodeVerifier string
	Nonce        string
}

type Token string

type OIDCUserInfo struct {
	Subject           string `json:"sub"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
}

type OIDCClient interface {
	NewAuthCodeURL(ctx context.Context) (*url.URL, *AuthProofs, error)
	ExchangeCodeToToken(
		ctx context.Context, code, state string, authProofs *AuthProofs,
	) (Token, error)
	GetUserInfo(ctx context.Context, token Token) (*OIDCUserInfo, error)
}

type UserSubjectRepository interface {
	GetUserBySubject(ctx context.Context, subject string) (*User, error)
	RegisterUserSubject(ctx context.Context, subject string, userID uuid.UUID) error
}

func (s *Service) NewAuthCodeURL(ctx context.Context) (*url.URL, *AuthProofs, error) {
	logger := logging.GetLogger(ctx)
	url, authProofs, err := s.oidcClient.NewAuthCodeURL(ctx)
	if err != nil {
		logger.Error("failed to create auth code URL", zap.Error(err))
		return nil, nil, err
	}
	return url, authProofs, nil
}

func (s *Service) ExchangeCodeToToken(
	ctx context.Context, code, state string, authProofs *AuthProofs,
) (Token, error) {
	logger := logging.GetLogger(ctx)
	token, err := s.oidcClient.ExchangeCodeToToken(ctx, code, state, authProofs)
	if err != nil {
		logger.Error("failed to exchange code to token", zap.Error(err))
		return "", err
	}
	return token, nil
}

// DecodeToken 与えられたトークンに対応するユーザー情報を返す。
// ユーザーが存在しない場合は新たにユーザーを作成して返す(JIT Provisioning)。
func (s *Service) DecodeToken(ctx context.Context, token Token) (*User, error) {
	logger := logging.GetLogger(ctx)
	userInfo, err := s.oidcClient.GetUserInfo(ctx, token)
	if err != nil {
		logger.Error("failed to get user info from OIDC client", zap.Error(err))
		return nil, err
	}
	user, err := s.repository.GetUserBySubject(ctx, userInfo.Subject)
	if err == nil {
		return user, nil
	} else if _, ok := err.(*NotFoundError); !ok { //nolint:errorlint
		logger.Error("failed to get user by subject from repository", zap.Error(err))
		return nil, err
	}
	// ユーザーが存在しない場合は新規作成する
	// TODO(db-transaction): ユーザーの作成とSubjectの登録は同一トランザクションで実行するべき
	newUser, err := s.repository.CreateUser(ctx, userInfo.Name, userInfo.PreferredUsername, false)
	if err != nil {
		logger.Error("failed to create user in repository", zap.Error(err))
		return nil, err
	}
	err = s.repository.RegisterUserSubject(ctx, userInfo.Subject, newUser.ID)
	if err != nil {
		logger.Error("failed to register user subject in repository", zap.Error(err))
		return nil, err
	}
	return newUser, nil
}
