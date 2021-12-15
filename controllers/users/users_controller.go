package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"users-api/domain/users"
	"users-api/services"
	"users-api/utils/errors"
)

func getUserId(userIdParam string) (int64, *errors.RestErr) {
	userId, userErr := strconv.ParseInt(userIdParam, 10, 64)
	if userErr != nil {
		return 0, errors.NewBadRequestError("user id should be a number")
	}
	return userId, nil
}

// @BasePath /api/v1

// Create is create user endpoint handler
// @Summary create user
// @Schemes
// @Description create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Success 201 {object} users.PrivateUser
// @Router /users/ [post]
func Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		//Handle user creation error
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshal(c.GetHeader("X-Public") == "true"))
}

// @BasePath /api/v1

// Get is get user endpoint handler
// @Summary get user
// @Schemes
// @Description return one user Based on id
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} users.PrivateUser
// @Router /users/{id} [get]
func Get(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}
	user, getErr := services.UsersService.GetUser(userId)
	if getErr != nil {
		//Handle user creation error
		c.JSON(getErr.Status, getErr)
		return
	}

	c.JSON(http.StatusOK, user.Marshal(c.GetHeader("X-Public") == "true"))
}

// @BasePath /api/v1

// GetAll is get users endpoint handler
// @Summary get all users
// @Schemes
// @Description list all the users based on filter given
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} users.PrivateUser
// @Router /users/all [get]
func GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	status := c.DefaultQuery("status", "active")

	user, getErr, total := services.UsersService.GetAllUser(page,status)

	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	users := user.Marshal(c.GetHeader("X-Public") == "true")
	res := make(map[string]interface{})
	res["total"] = total
	res["page"] = page
	res["data"] = users
	c.JSON(http.StatusOK, res)
}

// @BasePath /api/v1

// Update is update users endpoint handler
// @Summary update user
// @Schemes
// @Description update user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} users.PrivateUser
// @Router /users/ [put][patch]
func Update(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Id = userId

	isPartial := c.Request.Method == http.MethodPatch

	result, err := services.UsersService.UpdateUser(isPartial, user)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, result.Marshal(c.GetHeader("X-Public") == "true"))
}

// @BasePath /api/v1

// Delete is delete user endpoint handler
// @Summary delete user
// @Schemes
// @Description delete user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {string} string "deleted"
// @Router /users/ [delete]
func Delete(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	if err := services.UsersService.DeleteUser(userId); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// @BasePath /api/v1

// Search is search users endpoint handler
// @Summary search user
// @Schemes
// @Description list all the users based on filter given
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} users.PrivateUser
// @Router /internal/users/search [get]
func Search(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	status := c.DefaultQuery("status", "active")

	user, err , total := services.UsersService.SearchUser(page,status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	users := user.Marshal(c.GetHeader("X-Public") == "true")
	res := make(map[string]interface{})
	res["total"] = total
	res["page"] = page
	res["data"] = users
	c.JSON(http.StatusOK, res)
}
