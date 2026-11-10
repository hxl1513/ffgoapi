package users

import (
	"github.com/gin-gonic/gin"
	"github.com/hxl1513/ffgoapi/controller"
	"github.com/hxl1513/ffgoapi/dto"
)

// UserController handles HTTP requests for user resources.
type UserController struct {
	controller.BaseController
	svc *UserService
}

// NewUserController builds a UserController with the given service.
func NewUserController(svc *UserService) *UserController {
	return &UserController{svc: svc}
}

// RegisterRoutes mounts the user endpoints on the provided gin engine
// by creating a root router group.
func (c *UserController) RegisterRoutes(r *gin.Engine) {
	group := r.Group("")
	group.POST("/users", c.CreateUser)
	group.GET("/users/:id", c.GetUser)
}

// CreateUser handles POST /users.
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.BindJSON(ctx, &req); err != nil {
		c.BadRequest(ctx, "invalid request body", err)
		return
	}
	if err := req.Validate(); err != nil {
		c.BadRequest(ctx, "validation failed", err)
		return
	}

	resp, err := c.svc.CreateUser(ctx.Request.Context(), &req)
	if err != nil {
		c.InternalError(ctx, "failed to create user", err)
		return
	}
	c.Created(ctx, resp)
}

// GetUser handles GET /users/:id.
func (c *UserController) GetUser(ctx *gin.Context) {
	id := c.ParamInt(ctx, "id")
	if id == 0 {
		c.BadRequest(ctx, "invalid user id", nil)
		return
	}

	user, err := c.svc.Repo.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		c.InternalError(ctx, "failed to fetch user", err)
		return
	}
	if user == nil {
		c.NotFound(ctx, "user not found")
		return
	}

	c.Success(ctx, dto.FromUser(user))
}
