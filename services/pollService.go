package services

import (
	"context"
	"database/sql"
	"errors"
	"log"
	db "webpolls/db/sqlc"
)

//encapsula la lógica de negocio para las encuestas.
type PollService struct {
	Queries *db.Queries // <-- Exportado (con mayúscula)
	DB      *sql.DB
}

//crea una nueva instancia de PollService.
func NewPollService(queries *db.Queries, db *sql.DB) *PollService {
	return &PollService{Queries: queries, DB: db} 
}

type OptionResponse = db.Option

type OptionRequest struct {
	Content string `json:"content"`
}

type PollRequest struct {
	Question string          `json:"question"`
	Options  []OptionRequest `json:"options"`
	UserID   int32           `json:"user_id"`
}

type PollResponse struct {
	ID      int32            `json:"id"`
	Title   string           `json:"title"`
	UserID  int32            `json:"user_id"`
	Options []OptionResponse `json:"options"`
}

func (s *PollService) CreatePoll(ctx context.Context, params PollRequest) (*PollResponse, error) {
	if params.Question == "" {
		return nil, errors.New("la pregunta no puede estar vacía")
	}
	if len(params.Options) < 2 {
		return nil, errors.New("deben ser al menos 2 opciones")
	}
	if len(params.Options) > 4 {
		return nil, errors.New("deben ser máximo 4 opciones")
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := s.Queries.WithTx(tx)

	log.Println(params)
	poll, err := qtx.CreatePoll(ctx, db.CreatePollParams{
		Title:  params.Question,
		UserID: params.UserID,
	})
	if err != nil {
		return nil, err
	}

	// Crear opciones asociadas
	var options []db.Option
	for _, optionContent := range params.Options {
		option, err := qtx.CreateOption(ctx, db.CreateOptionParams{
			Content: optionContent.Content,
			PollID:  poll.ID,
		})
		if err != nil {
			return nil, err
		}
		options = append(options, option)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	data := &PollResponse{
		ID:      poll.ID,
		Title:   poll.Title,
		UserID:  poll.UserID,
		Options: options,
	}
	return data, nil
}

func (s *PollService) GetPollByID(ctx context.Context, id int32) (*PollResponse, error) {
	poll, err := s.Queries.GetPollByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var options []OptionResponse

	for _, pollRow := range poll {
		options = append(options, OptionResponse{
			ID:      pollRow.OptionID,
			Content: pollRow.OptionContent,
		})
	}

	return &PollResponse{
		ID:      poll[0].ID,
		Title:   poll[0].Title,
		UserID:  poll[0].UserID,
		Options: options,
	}, nil
}

func (s *PollService) GetPolls(ctx context.Context) ([]*PollResponse, error) {
	rows, err := s.Queries.GetAllPolls(ctx)
	if err != nil {
		return nil, err
	}

	pollsMap := make(map[int32]*PollResponse)

	for _, row := range rows {
		// Si la encuesta aún no está en el mapa, la creamos.
		if _, ok := pollsMap[row.PollID]; !ok {
			pollsMap[row.PollID] = &PollResponse{
				ID:      row.PollID,
				Title:   row.Title,
				UserID:  row.UserID,
				Options: []OptionResponse{},
			}
		}

		// Agregamos la opción actual a la encuesta correspondiente.
		poll := pollsMap[row.PollID]
		poll.Options = append(poll.Options, OptionResponse{
			ID:      row.OptionID,
			Content: row.OptionContent,
		})
	}

	// Convertimos el mapa de punteros a una lista de valores.
	var result []*PollResponse
	for _, poll := range pollsMap {
		result = append(result, poll)
	}

	return result, nil
}

func (s *PollService) DeletePoll(ctx context.Context, id int32) error {
	return s.Queries.DeletePoll(ctx, id)
}

func (s *PollService) UpdateOption(ctx context.Context, params OptionResponse) (*OptionResponse, error) {
	if params.Content == "" {
		return nil, errors.New("el contenido de la opción no puede estar vacío")
	}

	updatedOption, err := s.Queries.UpdateOption(ctx, db.UpdateOptionParams{
		ID:      params.ID,
		Content: params.Content,
	})

	if err != nil {
		return nil, err
	}

	return &OptionResponse{
		ID:      updatedOption.ID,
		Content: updatedOption.Content,
		PollID:  updatedOption.PollID,
	}, nil
}

func (s *PollService) DeleteOption(ctx context.Context, id int32, poll_id int32) error {
	options, err := s.Queries.GetOptionByPollID(ctx, poll_id)
	if err != nil {
		return err
	}

	if len(options) == 2 {
		return errors.New("la encuesta debe tener al menos 2 opciones")
	}

	for _, option := range options {
		if option.ID == id {
			return s.Queries.DeleteOption(ctx, id)
		}
	}

	return errors.New("opcion no encontrada")
}
