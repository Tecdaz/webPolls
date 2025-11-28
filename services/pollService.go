package services

import (
	"context"
	"errors"
	"log"
	db "webpolls/db/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PollService encapsula la lógica de negocio para las encuestas.
type PollService struct {
	Queries *db.Queries // <-- Exportado (con mayúscula)
	DB      *pgxpool.Pool
}

// NewPollService crea una nueva instancia de PollService.
func NewPollService(queries *db.Queries, db *pgxpool.Pool) *PollService {
	return &PollService{Queries: queries, DB: db} // <-- actualizado
}

type OptionResponse struct {
	ID         int32   `json:"id"`
	Content    string  `json:"content"`
	PollID     int32   `json:"poll_id"`
	VoteCount  int64   `json:"vote_count"`
	Percentage float64 `json:"percentage"`
}

type OptionRequest struct {
	Content string `json:"content"`
}

type PollRequest struct {
	Question string          `json:"question"`
	Options  []OptionRequest `json:"options"`
	UserID   int32           `json:"user_id"`
}

type PollResponse struct {
	ID                int32            `json:"id"`
	Title             string           `json:"title"`
	UserID            int32            `json:"user_id"`
	Options           []OptionResponse `json:"options"`
	TotalVotes        int64            `json:"total_votes"`
	UserVotedOptionID *int32           `json:"user_voted_option_id"`
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

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

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
		if optionContent.Content != "" {

			option, err := qtx.CreateOption(ctx, db.CreateOptionParams{
				Content: optionContent.Content,
				PollID:  poll.ID,
			})
			if err != nil {
				return nil, err
			}
			options = append(options, option)
		}

	}

	if len(options) < 2 {
		return nil, errors.New("deben ser al menos 2 opciones")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// Convert db.Option to OptionResponse
	var responseOptions []OptionResponse
	for _, opt := range options {
		responseOptions = append(responseOptions, OptionResponse{
			ID:      opt.ID,
			Content: opt.Content,
			PollID:  opt.PollID,
		})
	}

	data := &PollResponse{
		ID:      poll.ID,
		Title:   poll.Title,
		UserID:  poll.UserID,
		Options: responseOptions,
	}
	return data, nil
}

func (s *PollService) GetPollByID(ctx context.Context, id int32, userID *int32) (*PollResponse, error) {
	poll, err := s.Queries.GetPollByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get results
	results, err := s.Queries.GetPollResults(ctx, id)
	if err != nil {
		// If no results, just continue with 0 counts
		log.Printf("Error getting results for poll %d: %v", id, err)
	}

	voteCounts := make(map[int32]int64)
	var totalVotes int64
	for _, r := range results {
		voteCounts[r.OptionID] = r.VoteCount
		totalVotes += r.VoteCount
	}

	var userVotedOptionID *int32
	if userID != nil {
		votedOption, err := s.Queries.GetUserVote(ctx, db.GetUserVoteParams{
			PollID: id,
			UserID: *userID,
		})
		if err == nil {
			userVotedOptionID = &votedOption
		} else if err != pgx.ErrNoRows {
			log.Printf("Error checking user vote: %v", err)
		}
	}

	var options []OptionResponse

	for _, pollRow := range poll {
		count := voteCounts[pollRow.OptionID]
		percentage := 0.0
		if totalVotes > 0 {
			percentage = (float64(count) / float64(totalVotes)) * 100
		}

		options = append(options, OptionResponse{
			ID:         pollRow.OptionID,
			Content:    pollRow.OptionContent,
			PollID:     pollRow.ID,
			VoteCount:  count,
			Percentage: percentage,
		})
	}

	return &PollResponse{
		ID:                poll[0].ID,
		Title:             poll[0].Title,
		UserID:            poll[0].UserID,
		Options:           options,
		TotalVotes:        totalVotes,
		UserVotedOptionID: userVotedOptionID,
	}, nil
}

func (s *PollService) Vote(ctx context.Context, pollID int32, optionID int32, userID int32) error {
	return s.Queries.VoteOneStep(ctx, db.VoteOneStepParams{
		PollID:   pollID,
		OptionID: optionID,
		UserID:   userID,
	})
}

func (s *PollService) GetPolls(ctx context.Context, userID int32) ([]*PollResponse, error) {
	rows, err := s.Queries.GetAllPolls(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.mapGetAllPollsRowsToPolls(rows)
}

func (s *PollService) GetPollsByUser(ctx context.Context, ownerID int32, viewerID int32) ([]*PollResponse, error) {
	rows, err := s.Queries.GetPollsByUserID(ctx, db.GetPollsByUserIDParams{
		OwnerID:  ownerID,
		ViewerID: viewerID,
	})
	if err != nil {
		return nil, err
	}

	pollsMap := make(map[int32]*PollResponse)

	for _, row := range rows {
		if _, ok := pollsMap[row.PollID]; !ok {
			var userVotedOptionID *int32
			if row.UserVotedOptionID.Valid {
				id := row.UserVotedOptionID.Int32
				userVotedOptionID = &id
			}

			pollsMap[row.PollID] = &PollResponse{
				ID:                row.PollID,
				Title:             row.Title,
				UserID:            row.UserID,
				Options:           []OptionResponse{},
				UserVotedOptionID: userVotedOptionID,
			}
		}

		poll := pollsMap[row.PollID]
		poll.Options = append(poll.Options, OptionResponse{
			ID:      row.OptionID,
			Content: row.OptionContent,
		})
	}

	var result []*PollResponse
	for _, poll := range pollsMap {
		result = append(result, poll)
	}

	return result, nil
}

func (s *PollService) mapGetAllPollsRowsToPolls(rows []db.GetAllPollsRow) ([]*PollResponse, error) {
	pollsMap := make(map[int32]*PollResponse)

	for _, row := range rows {
		if _, ok := pollsMap[row.PollID]; !ok {
			var userVotedOptionID *int32
			if row.UserVotedOptionID.Valid {
				id := row.UserVotedOptionID.Int32
				userVotedOptionID = &id
			}

			pollsMap[row.PollID] = &PollResponse{
				ID:                row.PollID,
				Title:             row.Title,
				UserID:            row.UserID,
				Options:           []OptionResponse{},
				UserVotedOptionID: userVotedOptionID,
			}
		}

		poll := pollsMap[row.PollID]
		poll.Options = append(poll.Options, OptionResponse{
			ID:      row.OptionID,
			Content: row.OptionContent,
		})
	}

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
