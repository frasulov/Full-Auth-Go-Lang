package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Host string
	User string
	Password string
	Dbname string
	Port uint
	Sslmode string
	TimeZone string
}

func NewDatabase(host,user,password,dbname,sslmode,timezone string, port uint) *Database {
	return &Database{
		Host: host,
		User: user,
		Password: password,
		Dbname: dbname,
		Port: port,
		Sslmode: sslmode,
		TimeZone: timezone,
	}
}

func (db * Database) Connect() (*gorm.DB, error){
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v",
		db.Host,
		db.User,
		db.Password,
		db.Dbname,
		db.Port,
		db.Sslmode,
		db.TimeZone,
	)
	fmt.Println(dsn)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return conn, err
}
