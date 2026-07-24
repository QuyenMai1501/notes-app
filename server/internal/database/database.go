package database

import (
	"context"
	"fmt"
	"os"

	"notes-app/server/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://notes:notespass@localhost:5432/notesdb?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	query := `CREATE TABLE IF NOT EXISTS notes (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL DEFAULT '',
		tags TEXT[] NOT NULL DEFAULT '{}',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`

	if _, err := pool.Exec(ctx, query); err != nil {
		return fmt.Errorf("run migration: %w", err)
	}

	return nil
}

func CreateNote(ctx context.Context, pool *pgxpool.Pool, title, content string, tags []string) (models.Note, error) {
	if tags == nil {
		tags = []string{}
	}

	var n models.Note
	err := pool.QueryRow(ctx,
		`INSERT INTO notes (title, content, tags) VALUES ($1, $2, $3)
		 RETURNING id, title, content, tags, created_at, updated_at`,
		title, content, tags,
	).Scan(&n.ID, &n.Title, &n.Content, &n.Tags, &n.CreatedAt, &n.UpdatedAt)

	if err != nil {
		return models.Note{}, fmt.Errorf("create note: %w", err)
	}

	return n, nil
}

func GetNote(ctx context.Context, pool *pgxpool.Pool, id string) (models.Note, error) {
	var n models.Note
	err := pool.QueryRow(ctx,
		`SELECT id, title, content, tags, created_at, updated_at FROM notes WHERE id = $1`,
		id,
	).Scan(&n.ID, &n.Title, &n.Content, &n.Tags, &n.CreatedAt, &n.UpdatedAt)

	if err != nil {
		return models.Note{}, fmt.Errorf("get note: %w", err)
	}

	return n, nil
}

func ListNotes(ctx context.Context, pool *pgxpool.Pool, search, tag string) ([]models.Note, error) {
	query := `SELECT id, title, content, tags, created_at, updated_at FROM notes`
	var args []interface{}
	argIdx := 1
	var conditions []string

	if search != "" {
		conditions = append(conditions, fmt.Sprintf(`(title ILIKE $%d OR content ILIKE $%d)`, argIdx, argIdx))
		args = append(args, "%"+search+"%")
		argIdx++
	}

	if tag != "" {
		conditions = append(conditions, fmt.Sprintf(`$%d = ANY(tags)`, argIdx))
		args = append(args, tag)
		argIdx++
	}

	if len(conditions) > 0 {
		query += " WHERE "
		for i, c := range conditions {
			if i > 0 {
				query += " AND "
			}
			query += c
		}
	}

	query += " ORDER BY updated_at DESC"

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}
	defer rows.Close()

	notes := []models.Note{}
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Tags, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan note: %w", err)
		}
		notes = append(notes, n)
	}

	return notes, nil
}

func UpdateNote(ctx context.Context, pool *pgxpool.Pool, id, title, content string, tags []string) (models.Note, error) {
	if tags == nil {
		tags = []string{}
	}

	var n models.Note
	err := pool.QueryRow(ctx,
		`UPDATE notes SET title = $1, content = $2, tags = $3, updated_at = NOW()
		 WHERE id = $4
		 RETURNING id, title, content, tags, created_at, updated_at`,
		title, content, tags, id,
	).Scan(&n.ID, &n.Title, &n.Content, &n.Tags, &n.CreatedAt, &n.UpdatedAt)

	if err != nil {
		return models.Note{}, fmt.Errorf("update note: %w", err)
	}

	return n, nil
}

func DeleteNote(ctx context.Context, pool *pgxpool.Pool, id string) error {
	tag, err := pool.Exec(ctx, `DELETE FROM notes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("note not found")
	}
	return nil
}
