package util

import (
	"fmt"
	"net/http"
	"sub-payment-service/auth"
	"sub-payment-service/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// parse user claims in context and returns user
func GetUserFromContext(c echo.Context, db *gorm.DB) (*entity.User, error) {
	d := c.Get("user")
	token, ok := d.(*jwt.Token)
	if !ok {
		return nil, NewAppError(http.StatusInternalServerError, "failed to parse token", fmt.Sprintf("found type %T", d))
	}

	_, ok = token.Claims.(*auth.JwtAppClaims)
	if !ok {
		return nil, NewAppError(http.StatusInternalServerError, "failed to parse token", fmt.Sprintf("found type %T", token.Claims))
	}

	var user entity.User
	// err := db.Where("id=?", appClaims.UserID).First(&user).Error
	// if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, NewAppError(http.StatusBadRequest, "user not found", "")
	// } else if err != nil {
	// 	return nil, NewAppError(http.StatusInternalServerError, "server error", err.Error())
	// }

	return &user, nil
}
