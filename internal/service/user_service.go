package service

import (
	"context"
	"errors"

	"tung-tung-money-api/internal/model"
	"tung-tung-money-api/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	existing, _ := s.repo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, req *model.LoginRequest) (*model.User, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *UserService) Update(ctx context.Context, id uint, req *model.UpdateUserRequest) error {
	return s.repo.Update(ctx, id, req)
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
