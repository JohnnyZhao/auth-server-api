package app

import (
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/johnnyzhao/auth-server-api/internal/domain"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

const userKey = "user"

type Api struct {
	userRepo domain.UserStorage
}

func NewApi(repo domain.UserStorage) *Api {
	return &Api{userRepo: repo}
}

func (a *Api) CreateUser(c *gin.Context) {
	var payload CreateUserPayload
	_ = c.BindJSON(&payload)
	err := payload.ValidateRequired()
	if err != nil {
		a.handleBadRequest(c, "Account creation failed", err.Error())
		return
	}
	hashedPassword := hashPassword(payload.Password)
	if err != nil {
		a.handleInternalError(c)
		return
	}
	user := domain.User{
		UserID:         payload.UserID,
		HashedPassword: hashedPassword,
	}
	if err := a.userRepo.Create(c, &user); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			a.handleBadRequest(c, "Account creation failed", "already same user_id is used")
			return
		}
		a.handleInternalError(c)
		return
	}
	c.JSON(http.StatusOK, CreateUserResponse{
		Message: "Account successfully created!",
		User:    user,
	})
}

func (a *Api) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.Replace(c.GetHeader("Authorization"), "Basic ", "", -1)
		log.Println("...", header)

		headerBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(header))
		if err != nil {
			a.handleUnauthorized(c)
			c.Abort()
			return
		}
		decodedHeader := strings.TrimSpace(string(headerBytes))
		splits := strings.Split(decodedHeader, ":")
		if len(splits) != 2 {
			a.handleUnauthorized(c)
			c.Abort()
			return
		}
		userID, password := strings.TrimSpace(splits[0]), strings.TrimSpace(splits[1])
		user, err := a.userRepo.GetByUserID(c, userID)

		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			a.handleUnauthorized(c)
			c.Abort()
			return
		case err != nil:
			a.handleInternalError(c)
			c.Abort()
			return
		}

		if hashPassword(password) != user.HashedPassword {
			a.handleForbidden(c)
			c.Abort()
			return
		}
		c.Set(userKey, user.UserID)
		c.Next()
	}
}

func (a *Api) GetByUserID(c *gin.Context) {
	userID := c.Param("user_id")

	user, err := a.userRepo.GetByUserID(c, userID)

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		a.handleNotFound(c)
		return
	case err != nil:
		a.handleInternalError(c)
		return
	}

	c.JSON(http.StatusOK, GetUserResponse{
		Message: "User details by user_id",
		User:    user,
	})
}

func (a *Api) PatchByID(c *gin.Context) {
	var payload PatchUserPayload
	_ = c.BindJSON(&payload)
	if err := payload.ValidateRequired(); err != nil {
		a.handleBadRequest(c, "User updation failed", err.Error())
		return
	}

	userID := c.Param("user_id")
	user, err := a.userRepo.GetByUserID(c, userID)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		a.handleNotFound(c)
		return
	case err != nil:
		a.handleInternalError(c)
		return
	}

	contextUserID, _ := c.Get(userKey)
	if user.UserID != contextUserID.(string) {
		a.handleForbidden(c)
		return
	}

	values := make(map[string]interface{})
	if payload.Comment != nil {
		values["comment"] = *payload.Comment
	}
	if payload.Nickname != nil {
		values["nickname"] = *payload.Nickname
	}
	_, err = a.userRepo.UpdateByUserID(c, userID, values)
	if err != nil {
		a.handleInternalError(c)
		return
	}

	user, err = a.userRepo.GetByUserID(c, userID)
	if err != nil {
		a.handleInternalError(c)
		return
	}

	c.JSON(http.StatusOK, ItemResponse{
		Message: "User successfully updated!",
		Recipe:  user,
	})
}

func (a *Api) DeleteByID(c *gin.Context) {
	userID, _ := c.Get(userKey)
	_, err := a.userRepo.GetByUserID(c, userID.(string))
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		a.handleNotFound(c)
		return
	case err != nil:
		a.handleInternalError(c)
		return
	}

	if err = a.userRepo.DeleteByUserID(c, userID.(string)); err != nil {
		a.handleInternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account and user successfully removed!"})
}

func (a *Api) handleInternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Message: "Internal Server Error",
	})
}

func (a *Api) handleBadRequest(c *gin.Context, message, cause string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Message: message,
		Cause:   cause,
	})
	return
}

func (a *Api) handleUnauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Message: "Authorization Failed",
	})
	return
}

func (a *Api) handleForbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Message: "No Permission for Update",
	})
	return
}

func (a *Api) handleNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Message: "No User found",
	})
}
