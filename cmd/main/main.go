package main

import (
	"database/sql"
	"fmt"
	"log"
	"onepage/internal/middleware"
	"onepage/internal/storage"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Env struct {
	db *sql.DB
	st storage.Storage
}

var INVITE_COD string

func main() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config.yml") // name of config file (without extension)
	viper.AddConfigPath("./conf/")    // path to look for the config file in
	err := viper.ReadInConfig()       // Find and read the config file
	if err != nil {                   // Handle errors reading the config file
		log.Fatal("fatal error config file: ", err)
	}

	port := viper.GetInt("server.port")
	postgres_param := viper.GetString("BD.params")
	INVITE_COD = viper.GetString("INVITE_COD")
	cookieSalt := viper.GetString("cookie_salt")

	middleware.SetCookieSalt(cookieSalt)

	// env init
	db, err := sql.Open("postgres", postgres_param)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database offline")
	}
	log.Println("DB connected")

	env := &Env{db: db, st: storage.New(db)}
	// router init
	router := gin.Default()
	router.LoadHTMLGlob("v2/html_templates/*")

	router.Static("/v2/css", "./v2/css")
	router.Static("/v2/img", "./v2/img")
	router.Static("/v2/loaded_img", "./v2/loaded_img")
	router.Static("/v2/font", "./v2/font")
	router.Static("/v2/js", "./v2/js")
	router.Static("/v2/lightbox", "./v2/lightbox")
	router.StaticFile("/favicon.ico", "./favicon_io/favicon.ico")

	initializeRoutes(router, env)

	err = router.RunTLS(fmt.Sprintf(":%d", port), "cert/cert.pem", "cert/key.pem")
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
