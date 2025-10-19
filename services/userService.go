package services

import (
	"context"
	"database/sql"
	"errors"
	sqlc "webpolls/db/sqlc"
)

type UserService struct {
	queries *sqlc.Queries
}

type UserResponse struct {
	Id       int32  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewUserService(queries *sqlc.Queries) *UserService {
	return &UserService{queries: queries}
}
func (s *UserService) CreateUser(ctx context.Context, params sqlc.CreateUserParams) (*UserResponse, error) {
	if params.Username == "" || params.Email == "" || params.Password == "" {
		return nil, errors.New("Todos los campos son obligatorios")
	}

	_, err := s.queries.GetUserByUsername(ctx, params.Username)
	if err == nil {
		return nil, errors.New("El nombre de usuario ya existe")
	}

	_, err = s.queries.GetUserByEmail(ctx, params.Email)
	if err == nil {
		return nil, errors.New("El email ya existe")
	}

	createdRow, err := s.queries.CreateUser(ctx, params)
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
	userRow, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user := &UserResponse{
		Id:       userRow.ID,
		Username: userRow.Username,
		Email:    userRow.Email,
	}
	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) (string, error) {
	return s.queries.DeleteUser(ctx, id)
}

type UpdateUserParams struct {
	Username *string
	Email    *string
	Password *string
}

func (s *UserService) UpdateUser(ctx context.Context, id int32, params UpdateUserParams) (*UserResponse, error) {
	_, err := s.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.New("Usuario no encontrado")
	}

	var username, email, password sql.NullString

	if params.Username != nil {
		userByUsername, err := s.queries.GetUserByUsername(ctx, *params.Username)
		if err == nil && userByUsername.ID != id {
			return nil, errors.New("El nombre de usuario ya existe")
		}
		username = sql.NullString{String: *params.Username, Valid: true}
	}

	if params.Email != nil {
		userByEmail, err := s.queries.GetUserByEmail(ctx, *params.Email)
		if err == nil && userByEmail.ID != id {
			return nil, errors.New("El email ya existe")
		}
		email = sql.NullString{String: *params.Email, Valid: true}
	}

	if params.Password != nil {
		password = sql.NullString{String: *params.Password, Valid: true}
	}

	updatedRow, err := s.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:       id,
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	user := &UserResponse{
		Id:       updatedRow.ID,
		Username: updatedRow.Username,
		Email:    updatedRow.Email,
	}
	return user, nil
}
