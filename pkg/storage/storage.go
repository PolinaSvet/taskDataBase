package storage

const (
	AuthorsDb string = "ui/database/tableAuthors.json"
	PostsDb   string = "ui/database/tablePosts.json"
)

// Author - автор.
type Author struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Post - публикация.
type Post struct {
	ID             int64  `json:"id"`
	AuthorID       int64  `json:"author_id"`
	AuthorName     string `json:"author_name"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	CreatedAt      int64  `json:"created_at"`
	CreatedAtTxt   string `json:"created_at_txt"`
	PublishedAt    int64  `json:"published_at"`
	PublishedAtTxt string `json:"published_at_txt"`
}

type SqlResponse struct {
	ID  int    `json:"id"`
	Err string `json:"err"`
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Authors() ([]Author, error)                 // получение всех авторов
	AddAuthor(Author) (int, error)              // создание нового автора
	UpdateAuthor(Author) (int, error)           // обновление списка авторов
	DeleteAuthor(Author) (int, error)           // удаление автора по ID
	InsertInitDataFromFileAuthors(string) error // загрузить данные из файла

	Posts() ([]Post, error)                   // получение всех публикаций
	AddPost(Post) (int, error)                // создание новой публикации
	UpdatePost(Post) (int, error)             // обновление публикации
	DeletePost(Post) (int, error)             // удаление публикации по ID
	InsertInitDataFromFilePosts(string) error // загрузить данные из файла
}
