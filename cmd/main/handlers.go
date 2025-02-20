package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"onepage/internal/objects"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func (e *Env) AddBoard(c *gin.Context) {
	url := c.PostForm("url")
	name := c.PostForm("name")
	description := c.Query("description")

	// validate
	if url == "" || name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// validate
	for _, cc := range url {
		if cc == '/' {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	if err := e.st.AddBoard(objects.Board{Url: url, Name: name, Descripion: description}); err != nil {
		c.AbortWithStatus(400)
		return
	}

	switch c.Request.Header.Get("Accept") {
	// Respond with JSON
	case "application/json":
		c.AbortWithStatus(http.StatusOK)
	// Respond with HTML
	default:
		// c.HTML(http.StatusOK, HTMLtemplateName, data)
		c.Redirect(http.StatusFound, "/v2/boards")
	}
}

func (e *Env) GetAllBoards(c *gin.Context) {
	boards := e.st.GetAllBoards()

	type CaptionImgPair struct {
		C string
		I string
	}
	signs := []CaptionImgPair{
		{"Добро пожаловать?", ""},
		{"Добро пожаловать?", ""},
		{"Добро пожаловать?", ""},
		{"Добро пожаловать?", ""},
		{"Добро пожаловать?", ""},
		{"Добро пожаловать?", ""},
		{"Добро пожаловать?", ""},
		{"C - Си?", ""},
		{"C - Cargo?", ""},
		{"C - Cryptocolony?", "deg_anime.jpg"},
		// {"Абстрактное мышление, сослагательное наклонение - это как толерантность к лактозе - доступно не всем.", ""},
	}

	singIndex := rand.IntN(len(signs))
	if signs[singIndex].I == "" {
		signs[singIndex].I = "white_cat.png"
	}
	render_v2(c, gin.H{"payload": boards, "Caption": signs[singIndex]}, "index.html")
}

func (e *Env) GetThreadsFromBoard(c *gin.Context) {
	boardUrl := c.Params.ByName("board")
	threads := e.st.GetThreadsFromBoard(boardUrl)

	// get board name from DB
	boardName := e.st.GetBoardName(boardUrl)
	render_v2(c, gin.H{"payload": threads, "Board": boardName, "BoardUrl": boardUrl}, "threads.html")
}

func (e *Env) GetPostFromThread(c *gin.Context) {
	boardUrl := c.Params.ByName("board")
	idUrl := c.Params.ByName("id")

	boardName := e.st.GetBoardName(boardUrl)
	OP, posts := e.st.GetPostFromThread(boardUrl, idUrl)

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(200, posts)
		// Respond with HTML
	default:
		c.HTML(http.StatusOK, "posts.html", gin.H{"payload": posts, "boardName": boardName, "OP_post": OP})
	}
}

func (e *Env) AddThread(c *gin.Context) {
	// save OP post files
	fileNameForDB, err := SaveFile("name-for-file", c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	boardUrl := c.Params.ByName("board")
	text := c.PostForm("op-post")

	if err = e.st.AddThread(boardUrl, text, fileNameForDB); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		fmt.Println("создание треда: транзакция прошла успешно")
	}

	c.Redirect(http.StatusFound, "/v2/threads/"+boardUrl)
}

func (e *Env) AddPost(c *gin.Context) {
	// boardUrl := c.Params.ByName("board")
	idUrl := c.Params.ByName("id")
	text := c.PostForm("text")

	fileNameForDB, err := SaveFile("name-for-file", c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	status := http.StatusFound
	if err = e.st.AddPost(idUrl, text, fileNameForDB); err != nil {
		// status = http.StatusInternalServerError
	} else {
		fmt.Println("создание поста: транзакция прошла успешно")
	}

	redirectUrl := (*c.Request.URL).String()
	c.Redirect(status, redirectUrl)
}

func (e *Env) FlushAllTable(c *gin.Context) {
	e.st.FlushAllTable()
	c.Redirect(http.StatusFound, "/v2/boards")
}

func (e *Env) GetTest(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}

func (e *Env) PostTest(c *gin.Context) {
	// hadle text
	var err error
	login := c.PostForm("login")
	password := c.PostForm("password")
	postText := c.PostForm("post-text")
	// handle files
	fileHeader, _ := c.FormFile("name-for-file")
	fileName := "файла нет"
	fileNameForDB := ""
	if fileHeader != nil {
		fileName = fileHeader.Filename
		ext := filepath.Ext(fileName)
		filenameAsUNIXsecond := fmt.Sprint(time.Now().UnixNano())

		fileNameForDB = filenameAsUNIXsecond + ext
		err = c.SaveUploadedFile(fileHeader, "./test_img/"+fileNameForDB)
	}
	if err != nil {
		fileName += ("\nerror: " + err.Error())
	}

	// make report
	report := fmt.Sprintf("%#v %#v %#v %#v", login, password, postText, fileName)
	c.HTML(http.StatusOK, "admin.html", gin.H{"Info": report})
}

func (e *Env) NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", nil)
}

func (e *Env) Test(c *gin.Context) {
	fmt.Println("Test handler")
}
