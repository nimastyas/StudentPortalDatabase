package repository

import (
	"a21hc3NpZ25tZW50/model"
	"database/sql"
	"errors"
)

type SessionsRepository interface {
	AddSessions(session model.Session) error
	DeleteSession(token string) error
	UpdateSessions(session model.Session) error
	SessionAvailName(name string) error
	SessionAvailToken(token string) (model.Session, error)

	FetchByID(id int) (*model.Session, error)
}

type sessionsRepoImpl struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *sessionsRepoImpl {
	return &sessionsRepoImpl{db}
}

func (u *sessionsRepoImpl) AddSessions(session model.Session) error {
	row := u.db.QueryRow("INSERT INTO sessions (token, username, expiry) VALUES ($1, $2, $3) RETURNING id",
		session.Token, session.Username, session.Expiry)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (u *sessionsRepoImpl) DeleteSession(token string) error {
	stmt, err := u.db.Prepare("DELETE FROM sessions WHERE token = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(token)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("No session found with the specified token")
	}

	return nil

}

func (u *sessionsRepoImpl) UpdateSessions(session model.Session) error {
	_, err := u.db.Exec("UPDATE sessions SET token = $1, expiry = $2 WHERE username = $3",
		session.Token, session.Expiry, session.Username)
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (u *sessionsRepoImpl) SessionAvailName(name string) error {
	var session model.Session
	err := u.db.QueryRow("SELECT id, token, username, expiry FROM sessions WHERE username = $1", name).
		Scan(&session.ID, &session.Token, &session.Username, &session.Expiry)
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (u *sessionsRepoImpl) SessionAvailToken(token string) (model.Session, error) {
	var session model.Session
	err := u.db.QueryRow("SELECT id, token, username, expiry FROM sessions WHERE token = $1", token).
		Scan(&session.ID, &session.Token, &session.Username, &session.Expiry)
	if err != nil {
		return model.Session{}, err
	}
	return session, nil // TODO: replace this
}

func (u *sessionsRepoImpl) FetchByID(id int) (*model.Session, error) {
	row := u.db.QueryRow("SELECT id, token, username, expiry FROM sessions WHERE id = $1", id)

	var session model.Session
	err := row.Scan(&session.ID, &session.Token, &session.Username, &session.Expiry)
	if err != nil {
		return nil, err
	}

	return &session, nil
}
