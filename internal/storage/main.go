package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"onepage/internal/objects"
	"os"
	"path"
)

type Storage interface {
	GetAllBoards() []objects.Board
	GetThreadsFromBoard(string) []objects.Threads
	GetPostFromThread(boardUrl, idUrl string) (OP objects.Posts, posts []objects.Posts)

	GetBoardName(boardUrl string) string

	AddBoard(b objects.Board) error
	AddThread(boardUrl, text, fileName string) error
	AddPost(boardUrl, text, fileName string) error

	FlushAllTable()
}

const (
	//!! Если доски нет не возвращает ошибку
	// Можно исправить сделав два запроса в транзакции
	GET_THREADS_SQL_QUERY = `
	select id, op_post_id, text, image_url from threads_op
	join posts_base 
	on op_post_id = post_id
	where board_id = (select id from boards where url = $1)`

	GET_POSTS_SQL_QUERY = `
	select base.post_id, text, image_url from boards as b
	join threads_op as op on b.id = op.board_id
	join thread_posts p on op.id = p.fk_thread_id
	join posts_base as base on p.post_id = base.post_id
	where b.url = $1 and op.id = $2`

	GET_BOARD_NAME_SQL_QUERY = `
	select name from boards where url = $1;`

	GET_OP_POST_BY_THREAD_ID_SQL_QUERY = `
	select post_id, text, image_url from threads_op join posts_base on op_post_id = post_id where id = $1;
	`
)

type SQLStorage struct {
	db *sql.DB
}

func New(db *sql.DB) (st *SQLStorage) {
	st = &SQLStorage{db: db}
	return
}

func (st SQLStorage) GetAllBoards() []objects.Board {
	rows, err := st.db.Query("SELECT url, name, description FROM boards")
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("No result rows", err)
	} else if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	boards := []objects.Board{}
	for rows.Next() {
		board := objects.Board{}
		err := rows.Scan(&board.Url, &board.Name, &board.Descripion)
		if err != nil {
			log.Fatal(err)
		}
		boards = append(boards, board)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return boards
}

func (st SQLStorage) GetThreadsFromBoard(boardUrl string) []objects.Threads {
	// get threads from DB
	rows, err := st.db.Query(GET_THREADS_SQL_QUERY, boardUrl)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("No result rows", err)
	} else if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var image_url sql.NullString
	var text sql.NullString

	threads := []objects.Threads{}
	for rows.Next() {
		thread := objects.Threads{}
		err := rows.Scan(&thread.ID, &thread.OP_post_id, &text, &image_url)
		if err != nil {
			log.Fatal(err)
		}
		thread.Text = text.String
		thread.ImageUrl = image_url.String
		threads = append(threads, thread)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return threads
}

func (st SQLStorage) GetPostFromThread(boardUrl, idUrl string) (OP objects.Posts, posts []objects.Posts) {
	OP = objects.Posts{}
	rowP := st.db.QueryRow(GET_OP_POST_BY_THREAD_ID_SQL_QUERY, idUrl)
	rowP.Scan(&OP.Post_id, &OP.Text, &OP.ImageUrl)

	rows, err := st.db.Query(GET_POSTS_SQL_QUERY, boardUrl, idUrl)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("No result rows", err)
	} else if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	posts = []objects.Posts{}
	for rows.Next() {
		post := objects.Posts{}
		err := rows.Scan(&post.Post_id, &post.Text, &post.ImageUrl)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return
}

func (st SQLStorage) AddBoard(b objects.Board) error {
	_, err := st.db.Exec("insert into boards (url, name, description) values ($1, $2, $3)", b.Url, b.Name, b.Descripion)
	return err
}

func (st SQLStorage) GetBoardName(boardUrl string) string {
	boardName := ""
	rowP := st.db.QueryRow(GET_BOARD_NAME_SQL_QUERY, boardUrl)
	if rowP != nil {
		rowP.Scan(&boardName)
	}
	fmt.Println(boardName)
	return boardName
}

func (st SQLStorage) AddThread(boardUrl, text, fileName string) error {
	tx, err := st.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//!! проверять текст на валидность
	// пример: пост не пуст?
	lastInsertId := 0
	err = tx.QueryRow("insert into posts_base (text, image_url) values ($1, $2) returning post_id", text, fileName).Scan(&lastInsertId)
	if err != nil {
		return err
	}

	boardId := 0
	err = tx.QueryRow("select id from boards where url = $1", boardUrl).Scan(&boardId)
	if err != nil {
		fmt.Println("ошибка при обращении к таблице boards")
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("не возврашает строк -> такой борды не существует")
		}
		return err
	}

	_, err = tx.Exec("insert into threads_op (op_post_id, board_id) values ($1, $2)", lastInsertId, boardId)
	if err != nil {
		fmt.Println("ошибка при обращении к таблице threads_op")
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (st SQLStorage) AddPost(idUrl, text, fileName string) error {
	tx, err := st.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//!! проверять текст на валидность
	// пример: пост не пуст?

	lastInsertId := 0
	err = tx.QueryRow("insert into posts_base (text, image_url) values ($1, $2) returning post_id", text, fileName).Scan(&lastInsertId)
	if err != nil {
		return err
	}

	_, err = tx.Exec("insert into thread_posts (fk_thread_id, post_id) values ($1, $2)", idUrl, lastInsertId)
	if err != nil {
		fmt.Println("ошибка при обращении к таблице thread_posts")
		fmt.Println(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (st SQLStorage) FlushAllTable() {
	delete_cmds := []string{"delete from boards", "delete from posts_base", "delete from thread_posts",
		"delete from threads_op", "delete from users"}

	var err error
	for i := 0; ; i++ {
		errorFlag := false
		for _, cmd := range delete_cmds {
			_, err = st.db.Exec(cmd)
			if err != nil {
				errorFlag = true
			}
		}
		if errorFlag {
			fmt.Printf("Итерация %d произошла с ошибкой\n", i)
		} else {
			fmt.Println("Данные удалены")
			break
		}
	}

	dir, err := ioutil.ReadDir("v2/loaded_img/")
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{"v2/loaded_img/", d.Name()}...))
	}
}
