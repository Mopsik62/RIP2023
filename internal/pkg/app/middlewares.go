package app

import (
	"awesomeProject1/internal/app/ds"
	"awesomeProject1/internal/app/role"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strings"
)

const jwtPrefix = "Bearer "

func (a *Application) WithAuthCheck(assignedRoles ...role.Role) func(context *gin.Context) {
	return func(c *gin.Context) {
		isPassing := false
		for _, element := range assignedRoles {
			if element == role.Undefined {
				isPassing = true
				break
			}

		}

		jwtStr := c.GetHeader("Authorization")

		if jwtStr == "" {
			var cookieErr error
			jwtStr, cookieErr = c.Cookie("One-pot-api-token")
			if (cookieErr != nil) && (!isPassing) {
				c.AbortWithStatus(http.StatusBadRequest)
			}
		}

		if !isPassing && !strings.HasPrefix(jwtStr, jwtPrefix) {
			c.AbortWithStatus(http.StatusForbidden)

			return
		}

		if len(jwtStr) > len(jwtPrefix) {
			jwtStr = jwtStr[len(jwtPrefix):]
		}

		err := a.redis.CheckJWTInBlackList(c.Request.Context(), jwtStr)
		if err == nil && !isPassing { // значит что токен в блеклисте
			c.AbortWithStatus(http.StatusForbidden)

			return
		}

		if !isPassing && !errors.Is(err, redis.Nil) { // значит что это не ошибка отсуствия - внутренняя ошибка
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("test"), nil
		})
		if !isPassing && err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			log.Println(err)

			return
		}

		if len(jwtStr) <= len(jwtPrefix) {
			return
		}

		myClaims := token.Claims.(*ds.JWTClaims)

		isAssigned := false

		for _, oneOfAssignedRole := range assignedRoles {
			if oneOfAssignedRole == role.Undefined {
				c.Set("role", myClaims.Role)
				c.Set("userUUID", myClaims.UserUUID)
				return
			}
			if myClaims.Role == oneOfAssignedRole {
				isAssigned = true
				break
			}
		}

		if !isPassing && !isAssigned {
			c.AbortWithStatus(http.StatusForbidden)
			log.Printf("role %d is not assigned in %d", myClaims.Role, assignedRoles)
			return
		}

		c.Set("role", myClaims.Role)
		c.Set("userUUID", myClaims.UserUUID)

	}

}
