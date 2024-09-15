package controller

import (
	"fibercockroach/model"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Controller struct {
	Pool *pgxpool.Pool
}

type GuestListAllUsersIn struct {
	Offset int `form:"offset" query:"offset" json:"offset"`
	Limit  int `form:"limit" query:"limit" json:"limit"`
}

type GuestListAllUsersOut struct {
	Err   string              `json:"err,omitempty"`
	Users []model.User        `json:"users"`
	In    GuestListAllUsersIn `json:"input"`
}

func (ctr Controller) GuestListAllUsers(c *fiber.Ctx) error {
	in := GuestListAllUsersIn{}
	out := GuestListAllUsersOut{}
	_ = c.QueryParser(&in)
	if len(c.Body()) != 0 {
		err := c.BodyParser(&in)
		if err != nil {
			out.Err = `failed to parse body: ` + err.Error()
			return c.JSON(out)
		}
	}
	if in.Limit <= 0 {
		in.Limit = 10
	}
	out.In = in
	user := model.User{}
	users, err := user.AllUser(c.Context(), ctr.Pool, 10, 0)
	if err != nil {
		out.Err = `failed to retrieve users: ` + err.Error()
		return c.JSON(out)
	}
	out.Users = users

	return c.JSON(out)
}
