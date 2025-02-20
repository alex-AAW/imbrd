/* Тут объекты БД */
package objects

type User struct {
	Login    string
	Password string `json:"-"`
}

type Board struct {
	Url        string
	Name       string
	Descripion string
}

type Threads struct {
	ID         int
	OP_post_id int
	Text       string
	ImageUrl   string
}

type Posts struct {
	Post_id  int
	Text     string
	ImageUrl string
}
