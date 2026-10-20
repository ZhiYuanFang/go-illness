package jwt

import (
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/golang-jwt/jwt/v4"
)

// JWT declared for sign-unsign
type JWT struct {
	SigningKey []byte
}

// NewJWT 新建jwt实例
func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte("com.landun.service.micros"),
	}
}

type AuthClaims struct {
	Id int64 `json:"id"`
	jwt.StandardClaims
}

func (j *JWT) CreateToken(id int64) (accessToken string, refreshToken string, err error) {
	accessTokenModel := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute + 10).Unix(),
		},
	})
	accessToken, err = accessTokenModel.SignedString(j.SigningKey)

	refreshTokenModel := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour + 2).Unix(),
		},
	})
	refreshToken, err = refreshTokenModel.SignedString(j.SigningKey)
	return accessToken, refreshToken, err
}

func (j *JWT) Parse(tokenString string) (id int64, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	}, jwt.WithJSONNumber())

	// fmt.Println(token.Header)
	// // fmt.Println(token.Payload)
	// fmt.Println(token.Signature)
	if err != nil {
		return 0, gerror.NewCode(gcode.New(1004, "token格式错误", nil), err.Error())
	} else {
		claims, ok := token.Claims.(*AuthClaims)
		if ok {
			return claims.Id, nil
		} else {
			return 0, gerror.NewCode(gcode.New(1004, "密钥解析错误", nil), err.Error())
		}
	}
}

func NewToken(id int64) (string, string, error) {
	j := NewJWT()
	return j.CreateToken(id)
}

func ParseToken(tokenString string) (id int64, err error) {
	j := NewJWT()
	return j.Parse(tokenString)
}
