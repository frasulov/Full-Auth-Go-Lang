package jwt

import (
	"auth/config"
	"auth/global"
	"auth/models/redisModels"
	cache "auth/redis"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/square/go-jose.v2"
	"time"
)

func GenerateJWT(duration time.Duration, id string) (string, error) {
	fmt.Println(config.Configuration.Password.Jwt.SecretKey)
	var list []string
	list = append(list, "admin")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	payloadHasura := jwt.MapClaims{}
	payloadHasura["x-hasura-allowed-roles"] = list
	payloadHasura["x-hasura-default-role"] = "admin"
	payloadHasura["x-hasura-user-id"] = fmt.Sprintf("%v", id)
	claims["https://hasura.io/jwt/claims"] = payloadHasura
	claims["user-id"] = id
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(duration).Unix()
	tokenString, err := token.SignedString([]byte(config.Configuration.Password.Jwt.SecretKey))
	return tokenString, err
}

func GenerateJWE(token string) string {
	rcpt := jose.Recipient{
		Algorithm:  jose.PBES2_HS256_A128KW,
		Key:        config.Configuration.Password.Jwt.PublicKey,
		PBES2Count: 4096,
		PBES2Salt:  []byte(global.SALT_KEY),
	}
	encrypter, err := jose.NewEncrypter(jose.A128CBC_HS256, rcpt, nil)
	if err != nil {
		panic(err)
	}

	object, err := encrypter.Encrypt([]byte(token))
	if err != nil {
		panic("oops")
	}
	key, err := object.CompactSerialize()
	if err != nil {
		panic("oops")
	}
	return key
}

func DecryptJWE(jweToken string) string {
	jwe, err := jose.ParseEncrypted(jweToken)
	if err != nil {
		panic("oops")
	}
	decryptedKey, err := jwe.Decrypt(config.Configuration.Password.Jwt.PublicKey)
	if err != nil {
		panic("oops")
	}
	return string(decryptedKey)
}

func CreateSession(u *redisModels.UserSession) {
	redis := cache.NewRedisCache(fmt.Sprintf("%s:%v", config.Configuration.Redis.Host, config.Configuration.Redis.Port), config.Configuration.Redis.Db, time.Duration(config.Configuration.Redis.Expires))
	redis.Set(fmt.Sprintf("accessToken-of-%v", u.UserId), u)
}

func AddTokens(u *redisModels.UserSession) error {
	accessJWT, err := GenerateJWT(time.Minute*time.Duration(config.Configuration.Password.Jwt.AccessTokenExpire), u.UserId)
	if err != nil {
		return err
	}
	u.AccessToken = accessJWT
	refreshToken, err := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), 14)
	if err != nil {
		return err
	}
	u.RefreshToken = string(refreshToken)
	u.ExpiredAt = time.Now().Add(time.Minute * time.Duration(config.Configuration.Password.Jwt.RefreshTokenExpire)).Unix()
	return nil
}
