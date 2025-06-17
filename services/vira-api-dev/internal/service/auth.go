package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"vira-api-dev/internal/events"
	"vira-api-dev/internal/types"
	"vira-api-dev/internal/viraid"

	kafka "github.com/skrolikov/vira-kafka"
	log "github.com/skrolikov/vira-logger"
)

type AuthService struct {
	IDClient *viraid.Client
	Repo     types.UserProfileRepository
	Producer *kafka.Producer
	Logger   *log.Logger
}

func NewAuthService(idClient *viraid.Client, repo types.UserProfileRepository, producer *kafka.Producer, logger *log.Logger) *AuthService {
	return &AuthService{
		IDClient: idClient,
		Repo:     repo,
		Producer: producer,
		Logger:   logger,
	}
}

func (s *AuthService) RegisterProxy(ctx context.Context, req types.RegisterRequest, extra types.ProfileData, ip, userAgent string) (*types.DevAuthResponse, error) {
	authResp, err := s.IDClient.Register(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка регистрации через vira-id: %w", err)
	}

	uid := authResp.User.ID

	exists, err := s.Repo.Exists(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки профиля: %w", err)
	}

	prof := types.UserProfile{
		UserID:   uid,
		City:     extra.City,
		JoinedAt: time.Now(),
	}

	if !exists {
		if err := s.Repo.Create(ctx, prof); err != nil {
			return nil, fmt.Errorf("ошибка создания профиля: %w", err)
		}
	} else {
		prof, err = s.Repo.GetByUserID(ctx, uid)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения профиля: %w", err)
		}
	}

	go s.emitDevRegisterEvent(ctx, uid, extra.City)

	return &types.DevAuthResponse{
		Tokens:  authResp.Tokens,
		Profile: prof,
	}, nil
}

func (s *AuthService) LoginProxy(ctx context.Context, req types.LoginRequest, ip, userAgent string) (*types.DevAuthResponse, error) {
	authResp, err := s.IDClient.Login(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка авторизации через vira-id: %w", err)
	}

	uid := authResp.User.ID

	prof, err := s.Repo.GetByUserID(ctx, uid)
	if err != nil {
		if errors.Is(err, types.ErrProfileNotFound) {
			prof = types.UserProfile{
				UserID:   uid,
				JoinedAt: time.Now(),
			}
			if err := s.Repo.Create(ctx, prof); err != nil {
				return nil, fmt.Errorf("ошибка создания профиля: %w", err)
			}
		} else {
			return nil, fmt.Errorf("ошибка получения профиля: %w", err)
		}
	}

	go s.emitDevLoginEvent(ctx, uid, ip, userAgent)

	return &types.DevAuthResponse{
		Tokens:  authResp.Tokens,
		Profile: prof,
	}, nil
}

func (s *AuthService) emitDevRegisterEvent(ctx context.Context, userID, city string) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := events.EmitDevUserRegisteredEvent(ctxTimeout, s.Producer, s.Logger, userID, city, time.Now())
	if err != nil {
		s.Logger.Error("Ошибка отправки события регистрации: %v", err)
	}
}

func (s *AuthService) emitDevLoginEvent(ctx context.Context, userID, ip, userAgent string) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := events.EmitDevUserLoggedInEvent(ctxTimeout, s.Producer, s.Logger, userID, ip, userAgent, time.Now())
	if err != nil {
		s.Logger.Error("Ошибка отправки события входа: %v", err)
	}
}
