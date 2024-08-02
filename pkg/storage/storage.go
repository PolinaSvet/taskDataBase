package storage

const (
	AuthorsDb   string = "ui/database/tableAuthors.json"
	PostsDb     string = "ui/database/tablePosts.json"
	MongodbView string = "ui/database/mongodbView.js"
)

// Author - автор.
type Author struct {
	ID   int64  `json:"id"    bson:"_id"`
	Name string `json:"name"  bson:"name"`
}

// Post - публикация.
type Post struct {
	ID             int64  `json:"id"                bson:"_id"`
	AuthorID       int64  `json:"author_id"         bson:"author_id"`
	AuthorName     string `json:"author_name"       bson:"author_name"`
	Title          string `json:"title"             bson:"title"`
	Content        string `json:"content"           bson:"content"`
	CreatedAt      int64  `json:"created_at"        bson:"created_at"`
	CreatedAtTxt   string `json:"created_at_txt"    bson:"created_at_txt"`
	PublishedAt    int64  `json:"published_at"      bson:"published_at"`
	PublishedAtTxt string `json:"published_at_txt"  bson:"published_at_txt"`
}

type SqlResponse struct {
	ID  int64  `json:"id"`
	Err string `json:"err"`
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	GetInform() string

	Authors() ([]Author, error)                 // получение всех авторов
	AddAuthor(Author) (int64, error)            // создание нового автора
	UpdateAuthor(Author) (int64, error)         // обновление списка авторов
	DeleteAuthor(Author) (int64, error)         // удаление автора по ID
	InsertInitDataFromFileAuthors(string) error // загрузить данные из файла

	Posts() ([]Post, error)                   // получение всех публикаций
	AddPost(Post) (int64, error)              // создание новой публикации
	UpdatePost(Post) (int64, error)           // обновление публикации
	DeletePost(Post) (int64, error)           // удаление публикации по ID
	InsertInitDataFromFilePosts(string) error // загрузить данные из файла
}
