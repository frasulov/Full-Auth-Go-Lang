package config

import (
	"fmt"
	"reflect"
	"strconv"
)

type Configurations struct {
	Server Server
	Profile Profile
	Database DB
	Password Password
	Redis Redis
	Mail Mail
	SessionDeleteScheduler string
}

type Mail struct {
	Username, Password, Host string
}

type Redis struct {
	Host string
	Port int
	Db int
	Expires int
}

type Server struct {
	Port int		`default:"8080"`
	Host string		`default:"localhost"`
}

type Profile struct {
	Active string
}

type DB struct {
	Host string
	User string
	Password string
	Dbname string
	Sslmode string
	Timezone string
	Port int
}

type Password struct {
	ResetExpire int
	MinLength	int
	ActivateAccountTokenExpire int
	ForgotPasswordTokenExpire int
	Jwt Jwt
}

type Jwt struct {
	SecretKey string
	PublicKey string
	PrivateKey string
	SaltKey string
	AccessTokenExpire int
	RefreshTokenExpire int
}


func (s  Server) FillDefaults() Server {
	typ := reflect.TypeOf(s)
	if s.Host == "" {
		f, _ := typ.FieldByName("Host")
		s.Host = f.Tag.Get("default")
	}
	if s.Port == 0 {
		f, _ := typ.FieldByName("Port")
		fmt.Println(strconv.Atoi(f.Tag.Get("default")))
		port, err := strconv.Atoi(f.Tag.Get("default"))
		if err != nil{
			return s
		}
		s.Port = port
	}
	return s
}