package data

import (
	"backend/internal/validator"
	"context"
	"database/sql"
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
	Ping      int       `json:"ping"`
	UpdatedAt time.Time `json:"updated_at"`
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
	    RETURNING id, version
	`

	args := []any{c.Name, c.IP}

	return m.DB.QueryRowContext(ctx, stmt, args...).Scan(&c.ID, &c.Version)
}

func (m ContainerModel) Get(ctx context.Context, id int) (*Container, error) {
	//stmt := `
	//    SELECT id, name, ip, status, version
	//`

	return nil, nil
}
