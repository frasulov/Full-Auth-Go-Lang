package jwt

import (
	"auth/config"
	"auth/global"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/square/go-jose.v2"
	"time"
)


func GenerateJWT(duration time.Duration,id uint) (string, error){
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user-id"] = id
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(duration).Unix()
	tokenString, err := token.SignedString([]byte(config.Configuration.Password.Jwt.SecretKey))
	return tokenString, err
}


func GenerateJWE(token string) string{
	rcpt := jose.Recipient{
		Algorithm:  jose.PBES2_HS256_A128KW,
		Key:        config.Configuration.Password.Jwt.PublicKey,
		PBES2Count: 4096,
		PBES2Salt: []byte(global.SALT_KEY),
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


func DecryptJWE(jweToken string) string{
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