package service

import (
	"context"
	"errors"
	"time"
	"vira-api-dev/internal/types"
	"vira-api-dev/internal/viraid"
)

type AuthService struct {
	IDClient *viraid.Client
	Repo     types.UserProfileRepository
}

func NewAuthService(idClient *viraid.Client, repo types.UserProfileRepository) *AuthService {
	return &AuthService{IDClient: idClient, Repo: repo}
}

func (s *AuthService) RegisterProxy(ctx context.Context, req types.RegisterRequest, extra types.ProfileData) (*types.DevAuthResponse, error) {
	authResp, err := s.IDClient.Register(ctx, req)
	if err != nil {
		return nil, err
	}
	uid := authResp.User.ID

	exists, err := s.Repo.Exists(ctx, uid)
	if err != nil {
		return nil, err
	}
	if !exists {
		prof := types.UserProfile{
			UserID:   uid,
			City:     extra.City,
			JoinedAt: time.Now(),
		}
		if err := s.Repo.Create(ctx, prof); err != nil {
			return nil, err
		}
	}
	prof, err := s.Repo.GetByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &types.DevAuthResponse{
		Tokens:  authResp.Tokens,
		Profile: prof,
	}, nil
}

func (s *AuthService) LoginProxy(ctx context.Context, req types.LoginRequest) (*types.DevAuthResponse, error) {
	authResp, err := s.IDClient.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	uid := authResp.User.ID

	prof, err := s.Repo.GetByUserID(ctx, uid)
	if err != nil {
		if errors.Is(err, types.ErrProfileNotFound) {
			prof = types.UserProfile{UserID: uid}
			if err := s.Repo.Create(ctx, prof); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &types.DevAuthResponse{
		Tokens:  authResp.Tokens,
		Profile: prof,
	}, nil
}
