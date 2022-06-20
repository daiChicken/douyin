package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// 然后我们定义JWT的过期时间，这里以2小时为例：(优化成写入配置文件）
const TokenExpireDuration = time.Hour * 2

//定义密钥
var MySecret = []byte("douyin")

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
type MyClaims struct {
	UserId             int    `json:"user_id"`  //自定义字段
	Username           string `json:"username"` //自定义字段
	jwt.StandardClaims        //官方字段
}

// GenToken 生成JWT
func GenToken(userId int, username string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		userId,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(), // 过期时间
			Issuer:    "my-project",                                                                      // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
