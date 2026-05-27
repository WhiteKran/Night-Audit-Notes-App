package main

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type Note struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type App struct {
	ctx context.Context
	db  *sql.DB
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	homeDir, _ := os.UserHomeDir()
	dbDir := filepath.Join(homeDir, ".notesapp")
	os.MkdirAll(dbDir, 0755)

	dbPath := filepath.Join(dbDir, "notes.db")
	db, err := sql.Open("sqlite", "file:"+dbPath)
	if err != nil {
		panic(err)
	}

	a.db = db
	a.initDB()
}

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

func (a *App) initDB() {
	query := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := a.db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func (a *App) GetNotes() []Note {
	rows, err := a.db.Query("SELECT id, text, created_at, updated_at FROM notes ORDER BY created_at ASC")
	if err != nil {
		return []Note{}
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.ID, &note.Text, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			continue
		}
		notes = append(notes, note)
	}

	if notes == nil {
		notes = []Note{}
	}
	return notes
}

func (a *App) AddNote(text string) []Note {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := a.db.Exec(
		"INSERT INTO notes (text, created_at, updated_at) VALUES (?, ?, ?)",
		text, now, now,
	)
	if err != nil {
		return a.GetNotes()
	}

	return a.GetNotes()
}

func (a *App) UpdateNote(id int, text string) []Note {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := a.db.Exec(
		"UPDATE notes SET text = ?, updated_at = ? WHERE id = ?",
		text, now, id,
	)
	if err != nil {
		return a.GetNotes()
	}

	return a.GetNotes()
}

func (a *App) RemoveNote(id int) []Note {
	_, err := a.db.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return a.GetNotes()
	}

	return a.GetNotes()
}
