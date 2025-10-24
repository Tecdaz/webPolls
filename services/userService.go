package services

import (
	"context"
	"database/sql"
	"errors"
	db "webpolls/db/sqlc"
)

type UserService struct {
	Queries *db.Queries
}

type UserResponse struct {
	Id       int32  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserRequest = db.CreateUserParams

func NewUserService(queries *db.Queries) *UserService {
	return &UserService{Queries: queries}
}
func (s *UserService) CreateUser(ctx context.Context, params UserRequest) (*UserResponse, error) {
	if params.Username == "" || params.Email == "" || params.Password == "" {
		return nil, errors.New("Todos los campos son obligatorios")
	}

	_, err := s.Queries.GetUserByUsername(ctx, params.Username)
	if err == nil {
		return nil, errors.New("El nombre de usuario ya existe")
	}

	_, err = s.Queries.GetUserByEmail(ctx, params.Email)
	if err == nil {
		return nil, errors.New("El email ya existe")
	}

	createdRow, err := s.Queries.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	user := &UserResponse{
		Id:       createdRow.ID,
		Username: createdRow.Username,
		Email:    createdRow.Email,
	}
	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int32) (*UserResponse, error) {
	userRow, err := s.Queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		Id:       userRow.ID,
		Username: userRow.Username,
		Email:    userRow.Email,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) (string, error) {
	return s.Queries.DeleteUser(ctx, id)
}

type UpdateUserRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (s *UserService) UpdateUser(ctx context.Context, id int32, params UpdateUserRequest) (*UserResponse, error) {
	actualUser, err := s.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	var username, email, password sql.NullString

	if params.Username != nil && *params.Username != actualUser.Username {
		userByUsername, err := s.Queries.GetUserByUsername(ctx, *params.Username)
		if err == nil && userByUsername.ID != id {
			return nil, errors.New("el nombre de usuario ya existe")
		}
		username = sql.NullString{String: *params.Username, Valid: true}
	}

	if params.Email != nil && *params.Email != actualUser.Email {
		userByEmail, err := s.Queries.GetUserByEmail(ctx, *params.Email)
		if err == nil && userByEmail.ID != id {
			return nil, errors.New("el email ya existe")
		}
		email = sql.NullString{String: *params.Email, Valid: true}
	}

	if params.Password != nil {
		password = sql.NullString{String: *params.Password, Valid: true}
	}

	updatedRow, err := s.Queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:       id,
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		Id:       updatedRow.ID,
		Username: updatedRow.Username,
		Email:    updatedRow.Email,
	}, nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]UserResponse, error) {
	users, err := s.Queries.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	return userResponses, nil
}
