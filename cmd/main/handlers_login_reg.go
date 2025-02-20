package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"onepage/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (e *Env) Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func (e *Env) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (e *Env) PostRegister(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")
	invite_code := c.PostForm("invite-code")

	fmt.Println(login, password, invite_code)

	if !(len(login) > 0 && len(password) > 0 && len(invite_code) > 0) {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"loginErr": "ошибка валидации логина/пароля"})
		return
	}

	if invite_code != INVITE_COD {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"loginErr": "неверный инвайт"})
		return
	}

	lastInsertId := 0
	result := e.db.QueryRow("insert into users (login, password) values ($1, $2) returning id", login, password)
	err := result.Scan(&lastInsertId)

	if err != nil {
		log.Println(err)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Constraint == "users_unique" {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"loginErr": "такой пользователь уже существует"})
		} else {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"loginErr": "ошибка базы данных: " + err.Error()})
		}
		return
	}

	middleware.GinSetCookieID(c, lastInsertId)
	fmt.Println("Аккаунт создан успешно")

	c.Redirect(http.StatusFound, "/v2/boards")
}

func (e *Env) PostLogin(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	fmt.Println(login, password)

	if !(len(login) > 0 && len(password) > 0) {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"loginErr": "ошибка валидации логина/пароля"})
		return
	}

	id := 0
	err := e.db.QueryRow("SELECT id FROM users where login = $1 and password = $2", login, password).Scan(&id)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"loginErr": "пользователя с таким логином/паролем не существует"})
		} else {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"loginErr": "ошибка базы данных: " + err.Error()})
		}
		return
	}

	middleware.GinSetCookieID(c, id)

	fmt.Println("Вход выполнен")

	c.Redirect(http.StatusFound, "/v2/boards")
}
