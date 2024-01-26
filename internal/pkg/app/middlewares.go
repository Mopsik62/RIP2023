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

		jwtStr := c.GetHeader("Authorization")
		log.Println("НЕКУКИ= " + jwtStr)
		//var cookieeErr error
		//jwtStr, cookieeErr = c.Cookie("One-pot-api-token")
		//log.Println("cookie= " + c.Cookie("One-pot-api-token"))
		var cookiErr error
		jwtStr, cookiErr = c.Cookie("One-pot-api-token")
		log.Println("КУКИ= " + jwtStr)

		jwtStr = c.GetHeader("Authorization")
		if cookiErr != nil {
			log.Println("ОШИБКА КУКИ")
		}

		jwtStr = c.GetHeader("Authorization")

		if jwtStr == "" {
			var cookieErr error
			jwtStr, cookieErr = c.Cookie("One-pot-api-token")
			if cookieErr != nil {
				c.AbortWithStatus(http.StatusBadRequest)
			} else {
				//log.Println(jwtStr)
			}
		}
		log.Println("jwtStr под конец проверок = " + jwtStr)
		if !strings.HasPrefix(jwtStr, jwtPrefix) {
			c.AbortWithStatus(http.StatusForbidden)

			return
		}

		jwtStr = jwtStr[len(jwtPrefix):]

		err := a.redis.CheckJWTInBlackList(c.Request.Context(), jwtStr)
		if err == nil { // значит что токен в блеклисте
			c.AbortWithStatus(http.StatusForbidden)

			return
		}
		if !errors.Is(err, redis.Nil) { // значит что это не ошибка отсуствия - внутренняя ошибка
			c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("test"), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			log.Println(err)

			return
		}

		myClaims := token.Claims.(*ds.JWTClaims)

		isAssigned := false

		for _, oneOfAssignedRole := range assignedRoles {
			if myClaims.Role == oneOfAssignedRole {
				isAssigned = true
				break
			}
		}

		if !isAssigned {
			c.AbortWithStatus(http.StatusForbidden)
			log.Printf("role %s is not assigned in %v", myClaims.Role, assignedRoles)

			return
		}

		c.Set("role", myClaims.Role)
		c.Set("userUUID", myClaims.UserUUID)

	}

}
