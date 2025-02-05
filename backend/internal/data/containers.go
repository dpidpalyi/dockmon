package data

import (
	"backend/internal/validator"
	"context"
	"database/sql"
	"errors"
	"net/netip"
	"time"
)

type ContainerModel struct {
	DB *sql.DB
}

type Container struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Status  string `json:"status"`
	Version int    `json:"version"`
	// in ns
	Ping      int        `json:"ping"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func ValidateContainer(v *validator.Validator, container *Container) {
	v.Check(container.Name != "", "name", "must be provided")
	v.Check(container.Ping >= 0, "ping", "must be great or equal zero")
	if _, err := netip.ParseAddr(container.IP); err != nil {
		v.AddError("IP", "wrong IP specified")
	}
}

func (m ContainerModel) Insert(ctx context.Context, c *Container) error {
	stmt := `
	    INSERT INTO container(name, ip)
	    VALUES ($1, $2)
	    RETURNING id, version, status, ping, updated_at
	`

	args := []any{c.Name, c.IP}

	return m.DB.QueryRowContext(ctx, stmt, args...).Scan(
		&c.ID,
		&c.Version,
		&c.Status,
		&c.Ping,
		&c.UpdatedAt,
	)
}

func (m ContainerModel) Get(ctx context.Context, id int) (*Container, error) {
	stmt := `
	    SELECT id, name, ip, status, version, ping, updated_at
		FROM container
		WHERE id = $1`

	var c Container

	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(
		&c.ID,
		&c.Name,
		&c.IP,
		&c.Status,
		&c.Version,
		&c.Ping,
		&c.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &c, nil
}

func (m ContainerModel) Update(ctx context.Context, c *Container) error {
	stmt := `
	    UPDATE container
	    SET name = $1, ip = $2, status = $3, version = version + 1, ping = $4, updated_at = $5
	    WHERE id = $6 AND version = $7
	    RETURNING version`

	args := []any{
		c.Name,
		c.IP,
		c.Status,
		c.Ping,
		c.UpdatedAt,
		c.ID,
		c.Version,
	}

	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&c.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}

	return nil
}

func (m ContainerModel) Delete(ctx context.Context, id int) error {
	stmt := `
	    DELETE FROM container
	    WHERE id = $1`

	result, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m ContainerModel) List(ctx context.Context) ([]*Container, error) {
	stmt := `
	    SELECT id, name, ip, status, version, ping, updated_at
	    FROM container
	`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cs []*Container

	for rows.Next() {
		c := &Container{}
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.IP,
			&c.Status,
			&c.Version,
			&c.Ping,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cs, nil
}
