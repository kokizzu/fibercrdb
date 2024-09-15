package model

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id      int64  `db:"id" json:"id"`
	Email   string `db:"email" json:"email"`
	Pwdhash string `db:"pwdhash" json:"-"`
}

func WrapErr(err error, label string) error {
	if err != nil {
		return fmt.Errorf(label+`: %w`, err)
	}
	return nil
}

func (u *User) DoInsert(ctx context.Context, pool *pgxpool.Pool) error {
	query := `-- User) DoInsert
INSERT INTO users(email, pwdhash) 
VALUES ($1, $2)
ON CONFLICT DO NOTHING
RETURNING id
`
	row := pool.QueryRow(ctx, query, u.Email, u.Pwdhash)
	err := row.Scan(&u.Id)
	return WrapErr(err, `User.DoInsert row.Scan`)
}

func (u User) AllUser(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]User, error) {
	query := `-- User) AllUser
SELECT id, email 
FROM users
LIMIT $1
OFFSET $2`
	row, err := pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, WrapErr(err, `User.AllUser pool.Query`)
	}
	defer row.Close()
	res := make([]User, 0, limit)
	for row.Next() {
		line := User{}
		err := row.Scan(&line.Id, &line.Email)
		if err != nil {
			return nil, WrapErr(err, `User.AllUser row.Scan`)
		}
		res = append(res, line)
	}
	return res, nil
}
