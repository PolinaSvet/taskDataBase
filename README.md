
# DataBases + Go

<div align="center">
	<img src="https://i.ibb.co/hFfTZh9/1.jpg">
</div>


## Структура программы:

**1) Пакет "storage" содержит Interface, который задаёт контракт на работу с БД.**
***pkg\storage\storage.go***<br>

Программа использует две схемы таблиц и методы для работы с ними:


***Author - автор***<br>
type Author struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

Authors() ([]Author, error)                 // получение всех авторов
AddAuthor(Author) (int64, error)            // создание нового автора
UpdateAuthor(Author) (int64, error)         // обновление списка авторов
DeleteAuthor(Author) (int64, error)         // удаление автора по ID
InsertInitDataFromFileAuthors(string) error // загрузить данные из файла

***Post - публикация***<br>
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

Posts() ([]Post, error)                   // получение всех публикаций
AddPost(Post) (int64, error)              // создание новой публикации
UpdatePost(Post) (int64, error)           // обновление публикации
DeletePost(Post) (int64, error)           // удаление публикации по ID
InsertInitDataFromFilePosts(string) error // загрузить данные из файла

Есть возможность предварительной загрузки данных из файлов:
- Author:
***cmd\server\ui\database\tableAuthors.json***<br>
- Post:
***cmd\server\ui\database\tablePosts.json***<br>

**2) Пакет "api" реализует характерную для REST API схему запросов для работы с БД.**
***pkg\api\api.go***<br>
Запросы приходят на URL, соответствующий коллекции ресурсов:
- коллекции авторов "/authors"
- коллекции статей "/posts"

Для обозначения действий над коллекцией используются методы протокола HTTP: 
- POST для создания ресурса 
	api.router.HandleFunc("/authors", api.addAuthorHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/posts", api.addPostHandler).Methods(http.MethodPost, http.MethodOptions)

- DELETE для удаления
	api.router.HandleFunc("/authors", api.deleteAuthorHandler).Methods(http.MethodDelete, http.MethodOptions)
	api.router.HandleFunc("/posts", api.deletePostHandler).Methods(http.MethodDelete, http.MethodOptions)

- PUT для обновления
	api.router.HandleFunc("/authors", api.updateAuthorHandler).Methods(http.MethodPut, http.MethodOptions)
	api.router.HandleFunc("/posts", api.updatePostHandler).Methods(http.MethodPut, http.MethodOptions)

- GET для получения данных
	api.router.HandleFunc("/authors", api.authorsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/posts", api.postsHandler).Methods(http.MethodGet, http.MethodOptions)

**3) Для визуализации и организации REST API схемы запросов используется HTML+Javascript:**
***cmd\server\ui\html\base.html***<br>
***cmd\server\ui\html\routes.html***<br>

**4) Сервер хранит всю информацию в базе данных.**
Сервер предоставляет следующие реализации хранилища данных: 

- **postgres:** По аналогии с пакетом "memdb" разработан пакет "postgres" для поддержки базы данных под управлением СУБД PostgreSQL.
***pkg\storage\postgres\postgres.go***<br>
***cmd\server\ui\database\schema.sql*** - схема БД PostgreSQL в форме SQL-запроса<br>

-  **memdb:** Модернизирован пакет "memdb" реализована hash-структура для хранения данных.
***pkg\storage\memdb\memdb.go***<br>
type Store struct {
	AuthorsDB map[int64]storage.Author
	PostsDB   map[int64]storage.Post
}



## Требования к системе:

1) Сервер GoNews, который мы разрабатываем, должен предоставлять REST API, позволяющий выполнять следующие операции:
- Получение списка всех статей из БД,
- Добавление статьи в БД,
- Обновление статьи в БД,
- Удаление статьи из БД.

2) Сервер должен предоставлять данные в ответ на запросы по протоколу HTTP.
3) Сервер должен использовать характерную для REST API схему запросов:
- Запросы должны приходить на URL, соответствующий коллекции ресурсов. Например, коллекция статей.
- Для обозначения действий над коллекцией должны использоваться методы протокола HTTP: POST для создания ресурса, DELETE для удаления, PUT для обновления и GET для получения данных.
4) Сервер должен хранить всю информацию в базе данных.
5) Сервер должен предоставить как минимум две реализации хранилища данных: одну для реляционной СУБД и одну для документной.
- Объекты статьи должны содержать следующую информацию:
- Идентификатор,
- Имя автора,
- Заголовок,
- Текст,
- Время создания.

Для решения задачи от вас требуется следующее:

1) Разработать схему БД PostgreSQL в форме SQL-запроса. Запрос должен быть помещён в файл schema.sql в корневой каталог проекта.
2) По аналогии с пакетом "memdb" разработать пакет "postgres" для поддержки базы данных под управлением СУБД PostgreSQL.
3) По аналогии с пакетом "memdb" разработать пакет "mongo" для поддержки базы данных под управлением СУБД MongoDB.

Все выше перечисленные требования к системе учтены в проекте.
 

## Revision
 

- 1: package postgres: add tables authors, posts and functions for working with it
- 2: package memdb: modified package, added hash structure for tables authors, posts


## Usage:

**1.Enter this command to start the program:**

**go run server.go -typebd pg -loadbd yes**

1) typebd: This parameter is responsible for selecting the database.
- pg - PostgreSQL
- mem - memdb(map)
- mongo - MongoDB

2) loadbd: This parameter determines whether to preload the database from a file or not.
- yes - preload the database from a file
- no - not

**go run server.go**

defualt value (-typebd mem -loadbd yes)


**2.Open the web browser and go to:**

```sh

http://127.0.0.1:8080/ or  localhost:8080

```

## Authors:

@PolinaSvet

**!!! It is for test now !!!**