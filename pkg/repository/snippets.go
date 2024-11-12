package repository

import (
	"context"
	"database/sql"
	"deplagene/snippetbox/pkg/models"
	"errors"
)

func(repository *PGRepository) Get(id int) (*models.Snippet, error) {
	row := repository.pool.QueryRow(context.Background(), `
		SELECT snippet_id, title, content, created, expires
		FROM snippets
		WHERE expires > NOW()
		AND snippet_id = $1;
	`, id)

	data := &models.Snippet{}

	err := row.Scan(
		&data.ID,
		&data.Title,
		&data.Content,
		&data.Created,
		&data.Expires,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		} else {
			return nil, err
		}
	}

	return data, nil
}

func(repository *PGRepository) Create(title, content, expires string) (int, error) {
	var id int
	
	err := repository.pool.QueryRow(context.Background(), `
		INSERT INTO snippets (title, content, created, expires)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + $3::INTERVAL)
		RETURNING snippet_id;
	`, 
	title, 
	content, 
	expires + " days").Scan(&id)

	if err != nil {
		return 0, err
	}
	
	return id, nil
}

func(repository *PGRepository) Latest() ([]*models.Snippet, error) {
	
	rows, err := repository.pool.Query(context.Background(), `
		SELECT snippet_id, title, content, created, expires
		FROM snippets
		WHERE expires > NOW()
		ORDER BY created DESC
		LIMIT 10;
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var data []*models.Snippet

	for rows.Next() {
		item := &models.Snippet{}
		
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Content,
			&item.Created,
			&item.Expires,
		)

		if err != nil {
			return nil, err
		}

		data = append(data, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}