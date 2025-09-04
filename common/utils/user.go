package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructs "github.com/Orange-Health/citadel/common/structures"
)

func GetUserIdFromContext(c *gin.Context) (uint, *commonStructs.CommonError) {
	claims := ExtractClaimsFromJwtPayload(c)
	i := claims["citadel_user_id"]
	if i != nil {
		return uint(i.(float64)), nil
	}

	return 0, &commonStructs.CommonError{
		StatusCode: http.StatusUnauthorized,
		Message:    commonConstants.USER_ID_NOT_FOUND_IN_JWT,
	}
}

func GetCurrentLabIdFromContext(c *gin.Context) (uint, *commonStructs.CommonError) {
	claims := ExtractClaimsFromJwtPayload(c)
	i := claims["current_lab"]
	if i != nil {
		labIdString, ok := i.(string)
		if ok {
			labIdUint, err := strconv.ParseUint(labIdString, 10, 64)
			if err != nil {
				return 0, &commonStructs.CommonError{
					StatusCode: http.StatusUnauthorized,
					Message:    commonConstants.LAB_ID_NOT_FOUND_IN_JWT,
				}
			}
			return uint(labIdUint), nil
		}

		labIdFloat, ok := i.(float64)
		if ok {
			return uint(labIdFloat), nil
		}
	}

	return 0, &commonStructs.CommonError{
		StatusCode: http.StatusUnauthorized,
		Message:    commonConstants.LAB_ID_NOT_FOUND_IN_JWT,
	}
}

func ExtractClaimsFromJwtPayload(c *gin.Context) jwt.MapClaims {
	if _, exists := c.Get("JWT_PAYLOAD"); !exists {
		emptyClaims := make(jwt.MapClaims)
		return emptyClaims
	}

	jwtClaims, _ := c.Get("JWT_PAYLOAD")

	return jwtClaims.(jwt.MapClaims)
}
