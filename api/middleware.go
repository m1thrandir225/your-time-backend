package api

import (
	"errors"
	"m1thrandir225/your_time/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])

		if authorizationType != authorizationTypeBearer {
			err := errors.New("unsupported authorization type")
			ctx.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		accessToken := fields[1]

		payload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil {

			ctx.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)	 
		ctx.Next()
	}
}