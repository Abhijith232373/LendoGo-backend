package services

import (
    "errors"

    "lendogo-backend/internal/repositories"
    "lendogo-backend/structures/dto"
    "lendogo-backend/structures/models"
    "lendogo-backend/utils"
)

type AuthService interface {
    Register(req dto.RegisterReq) error
    Login(req dto.LoginReq) (*dto.AuthRes, error)
}

type authServiceImpl struct {
    userRepo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
    return &authServiceImpl{userRepo: repo}
}

func (s *authServiceImpl) Register(req dto.RegisterReq) error {
    existingUser, _ := s.userRepo.FindByEmail(req.Email)
    if existingUser != nil {
        return errors.New("email already in use")
    }

    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return errors.New("failed to secure password")
    }

    newUser := &models.User{
        FullName: req.FullName,
        Email:    req.Email,
        Password: hashedPassword,
    }

    return s.userRepo.CreateUser(newUser)
}

func (s *authServiceImpl) Login(req dto.LoginReq) (*dto.AuthRes, error) {
    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, errors.New("invalid email or password") 
    }

    isValid := utils.CheckPasswordHash(req.Password, user.Password)
    if !isValid {
        return nil, errors.New("invalid email or password")
    }

    // FIX: Convert UUID to string for token generation
    token, err := utils.GenerateToken(user.ID.String())
    if err != nil {
        return nil, errors.New("failed to generate login token")
    }

    // FIX: Correctly map the user details, ensuring ID is a string
    res := &dto.AuthRes{
        Token: token,
        User: dto.UserRes{
            ID:       user.ID.String(), 
            FullName: user.FullName,
            Email:    user.Email,
        },
    }

    return res, nil
}