package interfaces

import (
	"fmt"
	"food-app/application"
	"food-app/domain/entity"
	"food-app/infrastructure/auth"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Authenticate struct {
	us application.UserAppInterface
	rd auth.AuthInterface
	tk auth.TokenInterface
}

func NewAuthenticate(
	uApp application.UserAppInterface,
	rd auth.AuthInterface,
	tk auth.TokenInterface,
) *Authenticate {
	return &Authenticate{
		us: uApp,
		rd: rd,
		tk: tk,
	}
}

func (au *Authenticate) Login(c *gin.Context) {
	var user *entity.User
	var tokenErr = make(map[string]string)

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	// 验证请求
	validateUser := user.Validate("login")
	if len(validateUser) > 0 {
		c.JSON(http.StatusUnprocessableEntity, validateUser)
		return
	}
	u, userErr := au.us.GetUserByEmailAndPassword(user)
	if userErr != nil {
		c.JSON(http.StatusInternalServerError, userErr)
		return
	}
	ts, tErr := au.tk.CreateToken(u.ID)
	if tErr != nil {
		tokenErr["token_error"] = tErr.Error()
		c.JSON(http.StatusUnprocessableEntity, tErr.Error())
		return
	}
	saveErr := au.rd.CreateAuth(u.ID, ts)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr.Error())
		return
	}
	userData := make(map[string]interface{})
	userData["access_token"] = ts.AccessToken
	userData["refresh_token"] = ts.RefreshToken
	userData["id"] = u.ID
	userData["first_name"] = u.FirstName
	userData["last_name"] = u.LastName

	c.JSON(http.StatusOK, userData)
}

func (au *Authenticate) Logout(c *gin.Context) {
	metadata, err := au.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	// 如果 access token 存在并且仍然有效，则删除 access token 以及 refresh token
	deleteErr := au.rd.DeleteTokens(metadata)
	if deleteErr != nil {
		c.JSON(http.StatusUnauthorized, deleteErr.Error())
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

// 用于刷新 token 来产生新的 refresh 和 access token
func (au *Authenticate) Refresh(c *gin.Context) {
	mapToken := make(map[string]string)
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	// 验证 token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// 确保 token 方法为
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	// 其它由 token 过期而导致的错误
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	// token 是否有效
	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	// 因为 token 有效，获取 uuid
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) // 转换 interface 为 string
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, "Cannot get uuid")
			return
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Error occurred")
			return
		}
		// 删除之前的 refresh token
		delErr := au.rd.DeleteRefresh(refreshUuid)
		if delErr != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		// 创建一对新的 refresh 和 access tokens
		ts, createErr := au.tk.CreateToken(userId)
		if createErr != nil {
			c.JSON(http.StatusForbidden, createErr.Error())
			return
		}
		// 保存 tokens metadata 至 redis
		saveErr := au.rd.CreateAuth(userId, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, saveErr.Error())
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh token expired")
	}
}
