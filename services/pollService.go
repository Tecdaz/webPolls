package services

import (
	"context"
	"database/sql"
	"errors"
	sqlc "webpolls/db/sqlc"
)

// PollService encapsula la lógica de negocio para las encuestas.
type PollService struct {
	queries *sqlc.Queries
}

// NewPollService crea una nueva instancia de PollService.
func NewPollService(queries *sqlc.Queries) *PollService {
	return &PollService{queries: queries}
}

type OptionResponse struct {
	Content string `json:"content"`
	Correct bool   `json:"correct"`
}

type OptionCreatePoll struct {
	Content string `json:"content"`
	Correct bool   `json:"correct"`
}

type CreatePollParams struct {
	Question string
	Options  []OptionCreatePoll
	UserID   int32 // TODO: Usar el ID del usuario autenticado
}

type PollResponse struct {
	ID      int32            `json:"id"`
	Title   string           `json:"title"`
	UserID  int32            `json:"user_id"`
	Options []OptionResponse `json:"options"`
}

func (s *PollService) CreatePoll(ctx context.Context, params CreatePollParams) (*PollResponse, error) {
	if params.Question == "" {
		return nil, errors.New("La pregunta no puede estar vacia")
	}
	if len(params.Options) < 2 {
		return nil, errors.New("Deben ser al menos 2 opciones")
	}
	if len(params.Options) > 4 {
		return nil, errors.New("Deben ser maximo 4 opciones")
	}

	poll, err := s.queries.CreatePoll(ctx, sqlc.CreatePollParams{
		Title:  params.Question,
		UserID: params.UserID,
	})
	if err != nil {
		return nil, err
	}

	// TODO: validar opciones repetidas y solo una opcion correcta, tambien hacer rollback si va mal la creacion secuencial
	var options []sqlc.CreateOptionRow
	for _, optionContent := range params.Options {
		option, err := s.queries.CreateOption(ctx, sqlc.CreateOptionParams{
			Content: optionContent.Content,
			Correct: sql.NullBool{Bool: optionContent.Correct, Valid: true},
			PollID:  poll.ID,
		})
		if err != nil {
			return nil, err
		}
		options = append(options, option)
	}

	var optionsResponse []OptionResponse
	for _, option := range options {
		optionsResponse = append(optionsResponse, OptionResponse{
			Content: option.Content,
			Correct: option.Correct.Bool,
		})
	}

	data := &PollResponse{
		ID:      poll.ID,
		Title:   poll.Title,
		UserID:  poll.UserID,
		Options: optionsResponse,
	}
	return data, nil
}

func (s *PollService) GetPollByID(ctx context.Context, id int32) (*PollResponse, error) {
	poll, err := s.queries.GetPollByID(ctx, id)
	if err != nil {
		return nil, err
	}

	options, err := s.queries.GetOptionByPollID(ctx, id)
	if err != nil {
		return nil, err
	}
	var optionsResponse []OptionResponse
	for _, option := range options {
		optionsResponse = append(optionsResponse, OptionResponse{
			Content: option.Content,
			Correct: option.Correct.Bool,
		})
	}
	return &PollResponse{
		ID:      poll.ID,
		Title:   poll.Title,
		UserID:  poll.UserID,
		Options: optionsResponse,
	}, nil
}

func (s *PollService) DeletePoll(ctx context.Context, id int32) error {
	return s.queries.DeletePoll(ctx, id)
}
