package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/huandu/go-sqlbuilder"
	"github.com/vineboneto/rest-api-golang/internal/infra/db"
	pkg "github.com/vineboneto/rest-api-golang/pkg/types"
	"gorm.io/gorm"
)

type HandlerTenant struct {
	db *db.DB
}

type CreateTenantInput struct {
	Nome string `json:"nome"`
}

type CreateTenantOutput struct {
	Id string `json:"id"`
}

// CreateTenant godoc
//
//	@Summary		Cria um novo Tenant
//	@Description	Cria um novo Tenant com o nome fornecido
//	@Tags			tenant
//	@Accept			json
//	@Produce		json
//	@Param			input	body	CreateTenantInput	true	"Dados para criar um novo Tenant"
//	@Success		200		{object}	CreateTenantOutput	"Tenant criado com sucesso"
//	@Failure		400		{object}	ErrorResponse		"Erro na requisição"
//	@Router			/tenant [post]
func (h *HandlerTenant) CreateTenant(c *gin.Context) error {

	body := CreateTenantInput{}

	if err := c.BindJSON(&body); err != nil {
		AbortWithStatus(c).BadRequest(err)
		return nil
	}

	err := validation.ValidateStruct(&body,
		validation.Field(&body.Nome, validation.Required, validation.Length(5, 50)),
	)

	if err != nil {
		AbortWithStatus(c).BadRequest(err)
		return nil
	}

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()

	ib.InsertInto("tbl_tenant").Cols("nome")
	ib.Values(body.Nome)
	ib.SQL("RETURNING id")

	sql, args := ib.Build()

	var output CreateTenantOutput

	err = h.db.Conn.Transaction(func(tx *gorm.DB) error {
		err = tx.Raw(sql, args...).Scan(&output.Id).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, output)

	return nil
}

func (h *HandlerTenant) LoadAll(c *gin.Context) error {

	q := c.DefaultQuery("q", "")
	page := pkg.ParseIntOrDefault(c.DefaultQuery("page", "0"), 0)
	limit := pkg.ParseIntOrDefault(c.DefaultQuery("limit", "30"), 30)

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()

	sbCount := sqlbuilder.PostgreSQL.NewSelectBuilder()

	sb.Select("*").From("tbl_tenant")
	sbCount.Select("count(*)").From("tbl_tenant")

	if pkg.HasValue(q) {
		sb.Where(sb.ILike("nome", "%"+q+"%"))
		sbCount.Where(sbCount.ILike("nome", "%"+q+"%"))
	}

	sb.Limit(limit)
	sb.Offset(page * limit)

	type Output struct {
		Id   int    `gorm:"column:id" json:"id"`
		Nome string `gorm:"column:nome" json:"nome"`
	}

	sql, args := sb.Build()
	sqlCount, argsCount := sbCount.Build()
	output := []Output{}
	count := 0

	err := h.db.Conn.Raw(sql, args...).Scan(&output).Error

	if err != nil {
		return err
	}

	err = h.db.Conn.Raw(sqlCount, argsCount...).Scan(&count).Error

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  output,
		"count": count,
	})

	return nil
}

func (h *HandlerTenant) LoadById(c *gin.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return err
	}

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From("tbl_tenant")

	sb.Where(sb.Equal("id", id))

	type Output struct {
		Id   int    `gorm:"column:id" json:"id"`
		Nome string `gorm:"column:nome" json:"nome"`
	}

	sql, args := sb.Build()

	output := Output{}

	query := h.db.Conn.Raw(sql, args...).Scan(&output)

	if query.Error != nil {
		return query.Error
	}

	if query.RowsAffected == 0 {
		c.Status(http.StatusNoContent)
		return nil
	}

	c.JSON(http.StatusOK, output)

	return nil
}

func (h *HandlerTenant) Update(c *gin.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return err
	}

	type Body struct {
		Nome string `json:"nome"`
	}
	body := Body{}

	if err := c.BindJSON(&body); err != nil {
		AbortWithStatus(c).BadRequest(err)
		return nil
	}

	sb := sqlbuilder.PostgreSQL.NewUpdateBuilder()

	sb.Update("tbl_tenant")
	sb.Where(sb.Equal("id", id))

	if pkg.HasValue(body.Nome) {
		sb.Set(sb.Assign("nome", body.Nome))
	}

	sql, args := sb.Build()

	if len(args) < 2 {
		c.Status(http.StatusOK)
		return nil
	}

	err = h.db.Conn.Exec(sql, args...).Error

	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func NewHandlerTenant(db *db.DB) HandlerTenant {
	return HandlerTenant{db: db}
}
