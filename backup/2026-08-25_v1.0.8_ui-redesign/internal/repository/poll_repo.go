package repository

import (
	"time"

	"khan/internal/database"
	"khan/internal/models"
)

// PollRepo handles polls (نظرسنجی)
type PollRepo struct {
	store *database.Store
}

func NewPollRepo(store *database.Store) *PollRepo { return &PollRepo{store: store} }

// Create inserts a poll and returns its id
func (r *PollRepo) Create(roomID, creatorID int64, question string, options []string) (int64, error) {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()

	id := r.store.NextID("polls")
	r.store.Data().Polls = append(r.store.Data().Polls, database.PollRecord{
		ID:        id,
		RoomID:    roomID,
		CreatorID: creatorID,
		Question:  question,
		Options:   options,
		CreatedAt: time.Now().Format(time.RFC3339),
	})
	return id, r.store.SaveLocked()
}

// GetByID returns a poll
func (r *PollRepo) GetByID(id int64) (*models.Poll, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	for _, rec := range r.store.Data().Polls {
		if rec.ID == id {
			p := &models.Poll{
				ID: rec.ID, RoomID: rec.RoomID, CreatorID: rec.CreatorID,
				Question: rec.Question, Options: rec.Options,
				Closed: rec.Closed,
			}
			if t, err := time.Parse(time.RFC3339, rec.CreatedAt); err == nil {
				p.CreatedAt = t
			}
			if rec.ClosedAt != "" {
				if t, err := time.Parse(time.RFC3339, rec.ClosedAt); err == nil {
					p.ClosedAt = &t
				}
			}
			// Attach votes
			votes := make([]models.PollVote, 0)
			for _, v := range r.store.Data().PollVotes {
				if v.PollID == id {
					votes = append(votes, models.PollVote{
						ID: v.ID, PollID: v.PollID, UserID: v.UserID, Option: v.Option,
					})
				}
			}
			p.Votes = votes
			return p, nil
		}
	}
	return nil, nil
}

// Vote casts or changes a user's vote on a poll option
func (r *PollRepo) Vote(pollID, userID int64, option int) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()

	// Validate poll exists + not closed + option in range
	var poll *database.PollRecord
	for i := range r.store.Data().Polls {
		if r.store.Data().Polls[i].ID == pollID {
			poll = &r.store.Data().Polls[i]
			break
		}
	}
	if poll == nil || poll.Closed {
		return nil
	}
	if option < 0 || option >= len(poll.Options) {
		return nil
	}

	// Remove existing vote, then add new
	var kept []database.PollVoteRecord
	for _, v := range r.store.Data().PollVotes {
		if !(v.PollID == pollID && v.UserID == userID) {
			kept = append(kept, v)
		}
	}
	r.store.Data().PollVotes = kept
	r.store.Data().PollVotes = append(r.store.Data().PollVotes, database.PollVoteRecord{
		ID:      r.store.NextID("poll_votes"),
		PollID:  pollID,
		UserID:  userID,
		Option:  option,
		VotedAt: time.Now().Format(time.RFC3339),
	})
	return r.store.SaveLocked()
}

// Close closes a poll
func (r *PollRepo) Close(pollID int64) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	for i := range r.store.Data().Polls {
		if r.store.Data().Polls[i].ID == pollID {
			r.store.Data().Polls[i].Closed = true
			r.store.Data().Polls[i].ClosedAt = time.Now().Format(time.RFC3339)
			return r.store.SaveLocked()
		}
	}
	return nil
}

// ListByRoom returns polls in a room
func (r *PollRepo) ListByRoom(roomID int64) ([]models.Poll, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	var out = make([]models.Poll, 0)
	for _, rec := range r.store.Data().Polls {
		if rec.RoomID == roomID {
			p := models.Poll{ID: rec.ID, RoomID: rec.RoomID, CreatorID: rec.CreatorID,
				Question: rec.Question, Options: rec.Options, Closed: rec.Closed}
			if t, err := time.Parse(time.RFC3339, rec.CreatedAt); err == nil {
				p.CreatedAt = t
			}
			votes := make([]models.PollVote, 0)
			for _, v := range r.store.Data().PollVotes {
				if v.PollID == rec.ID {
					votes = append(votes, models.PollVote{
						ID: v.ID, PollID: v.PollID, UserID: v.UserID, Option: v.Option,
					})
				}
			}
			p.Votes = votes
			out = append(out, p)
		}
	}
	return out, nil
}
