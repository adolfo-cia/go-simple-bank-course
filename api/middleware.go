package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/adolfo-cia/go-simple-bank-course/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey     = "authorization"
	authorizationTypeBearer    = "bearer"
	authorizationPayloadCtxKey = "authPayload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				errorResponse(errors.New("missing authentication header")))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				errorResponse(errors.New("invalid authorization header format")))
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != authorizationTypeBearer {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				errorResponse(errors.New("unsupported authorization type: "+authType)))
			return
		}

		accessToken := fields[1]

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				errorResponse(err))
		}

		ctx.Set(authorizationPayloadCtxKey, payload)
		ctx.Next()
	}
}
