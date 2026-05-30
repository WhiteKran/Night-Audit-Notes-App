package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	_ "modernc.org/sqlite"
)

type WindowSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Note struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	IsLocked  bool   `json:"isLocked"`
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
	a.SetAppResolution(a.GetLatestResolution())
	a.ListenForResize()
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
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_locked BOOLEAN DEFAULT 0
	);
	`
	_, err := a.db.Exec(query)
	if err != nil {
		panic(err)
	}

	resolutionQuery := `
	CREATE TABLE IF NOT EXISTS resolution (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		width INTEGER NOT NULL,
		height INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	a.db.Exec(resolutionQuery)

	// Add is_locked column if it doesn't exist (for existing databases)
	a.db.Exec("ALTER TABLE notes ADD COLUMN is_locked BOOLEAN DEFAULT 0")
}

func (a *App) AddLatestResolution(size WindowSize) error {
	_, err := a.db.Exec(
		"INSERT INTO resolution (width, height) VALUES (?, ?)",
		size.Width, size.Height,
	)
	return err
}

func (a *App) GetLatestResolution() WindowSize {
	var width, height int
	row := a.db.QueryRow("SELECT width, height FROM resolution ORDER BY created_at DESC LIMIT 1")
	err := row.Scan(&width, &height)

	if err != nil {
		return WindowSize{Width: 600, Height: 600}
	}

	return WindowSize{Width: width, Height: height}
}

func (a *App) SetAppResolution(size WindowSize) error {
	fmt.Printf("Setting app resolution to %d x %d\n", size.Width, size.Height)
	runtime.WindowSetSize(a.ctx, size.Width, size.Height)
	return a.AddLatestResolution(size)
}

func (a *App) ListenForResize() {
	fmt.Println("Listening for window resize events")
	runtime.EventsOn(a.ctx, "wails:windowResized", func(optionalData ...interface{}) {
		width, height := runtime.WindowGetSize(a.ctx)
		fmt.Printf("Window resized to %d x %d, saving to database\n", width, height)
		a.AddLatestResolution(WindowSize{Width: width, Height: height})
	})
}

func (a *App) OnWindowResized(width int, height int) error {
	fmt.Printf("Window resized to %d x %d, saving to database\n", width, height)
	return a.AddLatestResolution(WindowSize{Width: width, Height: height})
}

func (a *App) GetNotes() []Note {
	rows, err := a.db.Query("SELECT id, text, created_at, updated_at, is_locked FROM notes ORDER BY created_at ASC")
	if err != nil {
		return []Note{}
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.ID, &note.Text, &note.CreatedAt, &note.UpdatedAt, &note.IsLocked)
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

func (a *App) ToggleLock(id int) []Note {
	row := a.db.QueryRow("SELECT is_locked FROM notes WHERE id = ?", id)
	var isLocked bool
	err := row.Scan(&isLocked)
	if err != nil {
		return a.GetNotes()
	}

	_, err = a.db.Exec("UPDATE notes SET is_locked = ? WHERE id = ?", !isLocked, id)
	if err != nil {
		return a.GetNotes()
	}

	return a.GetNotes()
}
