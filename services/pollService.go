package services

import (
	"context"
	"database/sql"
	"errors"
	db "webpolls/db/sqlc"
	sqlc "webpolls/db/sqlc"
)

// PollService encapsula la lógica de negocio para las encuestas.
type PollService struct {
	Queries *sqlc.Queries // <-- Exportado (con mayúscula)
}

// NewPollService crea una nueva instancia de PollService.
func NewPollService(queries *sqlc.Queries) *PollService {
	return &PollService{Queries: queries} // <-- actualizado
}

type OptionResponse struct {
	ID      int    `json:"id"` // <--- necesario para PUT
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
	UserID   int32 //  usar el ID del usuario autenticado
}

type PollResponse struct {
	ID      int32            `json:"id"`
	Title   string           `json:"title"`
	UserID  int32            `json:"user_id"`
	Options []OptionResponse `json:"options"`
}

type PollWithOptions struct {
	ID      int32            `json:"poll_id"`
	Title   string           `json:"title"`
	UserID  int32            `json:"user_id"`
	Options []OptionResponse `json:"options"`
}

func (s *PollService) CreatePoll(ctx context.Context, params CreatePollParams) (*PollResponse, error) {
	if params.Question == "" {
		return nil, errors.New("La pregunta no puede estar vacía")
	}
	if len(params.Options) < 2 {
		return nil, errors.New("Deben ser al menos 2 opciones")
	}
	if len(params.Options) > 4 {
		return nil, errors.New("Deben ser máximo 4 opciones")
	}

	poll, err := s.Queries.CreatePoll(ctx, sqlc.CreatePollParams{
		Title:  params.Question,
		UserID: params.UserID,
	})
	if err != nil {
		return nil, err
	}

	// Crear opciones asociadas
	var options []sqlc.CreateOptionRow
	for _, optionContent := range params.Options {
		option, err := s.Queries.CreateOption(ctx, sqlc.CreateOptionParams{
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
			ID:      int(option.ID),
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
	poll, err := s.Queries.GetPollByID(ctx, id)
	if err != nil {
		return nil, err
	}

	options, err := s.Queries.GetOptionByPollID(ctx, id)
	if err != nil {
		return nil, err
	}

	var optionsResponse []OptionResponse
	for _, option := range options {
		optionsResponse = append(optionsResponse, OptionResponse{
			ID:      int(option.ID),
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

func (s *PollService) GetPolls(ctx context.Context) ([]db.GetAllPollsRow, error) {
	return s.Queries.GetAllPolls(ctx)
}

func (s *PollService) DeletePoll(ctx context.Context, id int32) error {
	return s.Queries.DeletePoll(ctx, id)
}

func (s *PollService) GetPollsWithOptions(ctx context.Context) ([]PollWithOptions, error) {
	rows, err := s.Queries.GetAllPolls(ctx)
	if err != nil {
		return nil, err
	}

	// Vamos a agrupar las opciones por encuesta
	pollsMap := make(map[int32]*PollWithOptions)
	for _, row := range rows {
		if _, exists := pollsMap[row.PollID]; !exists {
			pollsMap[row.PollID] = &PollWithOptions{
				ID:      row.PollID,
				Title:   row.Title,
				UserID:  row.UserID,
				Options: []OptionResponse{},
			}
		}

		// Si hay opción (puede ser NULL por el LEFT JOIN)
		if row.OptionID.Valid {
			pollsMap[row.PollID].Options = append(pollsMap[row.PollID].Options, OptionResponse{
				ID:      int(row.OptionID.Int32),
				Content: row.Content.String,
				Correct: row.Correct.Bool,
			})
		}
	}

	// Convertir el mapa a slice
	var result []PollWithOptions
	for _, poll := range pollsMap {
		result = append(result, *poll)
	}

	return result, nil
}
