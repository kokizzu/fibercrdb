package main

import (
	"context"
	"embed"
	"errors"
	"log"

	"fibercockroach/controller"
	"fibercockroach/model"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/kokizzu/gotro/S"
	"github.com/pressly/goose/v3"
)

//go:embed dbschema/*.sql
var dbSchema embed.FS

const salt = `blabla`

// https://gofiber.io
// https://www.cockroachlabs.com/docs/releases/
// https://github.com/jackc/pgx/v4/pgxpool
// https://github.com/pressly/goose
// https://github.com/air-verse/air
func main() {
	config, err := pgxpool.ParseConfig(`postgresql://root@127.0.0.1:26257/defaultdb?sslmode=disable`)
	if err != nil {
		panic(`failed to connect db`)
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		log.Println(`database connected`)
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	goose.SetBaseFS(dbSchema)
	sqlDB := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(sqlDB, "dbschema"); err != nil {
		panic(err)
	}

	adminUser := model.User{
		Email:   `root@localhost`,
		Pwdhash: S.HashPassword(salt + `testing123`),
	}
	err = adminUser.DoInsert(context.Background(), pool)
	if !errors.Is(err, pgx.ErrNoRows) {
		panic(err)
	}

	app := fiber.New(fiber.Config{
		Immutable: true,
	})

	ctr := controller.Controller{
		Pool: pool,
	}

	app.All("/guest/listAllUsers", ctr.GuestListAllUsers)

	log.Fatal(app.Listen(":4523"))
}
