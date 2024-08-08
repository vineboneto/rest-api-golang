package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/vineboneto/rest-api-golang/internal/infra/db"
)

type HandlerTenant struct {
	pool *db.DB
}

func (h *HandlerTenant) CreateTenant(c *gin.Context) {

	type Body struct {
		Nome      string `json:"nome"`
		Sobrenome string `json:"sobrenome"`
	}

	body := Body{}

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := validation.ValidateStruct(&body,
		validation.Field(&body.Nome, validation.Required, validation.Length(5, 50)),
		validation.Field(&body.Sobrenome, validation.Required, validation.Length(5, 50)),
	)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()

	ib.InsertInto("tbl_tenant").Cols("nome")
	ib.Values(body.Nome)
	ib.SQL("RETURNING id")

	sql, args := ib.Build()

	log.Println(sql, args)

	err = h.pool.WithTransaction(ctx, func(tx pgx.Tx) error {

		_, err = tx.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", "Alice")
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", "Bob")
		if err != nil {
			return err
		}

		err = tx.Commit(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return
	}

}

func NewHandlerTenant() HandlerTenant {
	return HandlerTenant{}
}
