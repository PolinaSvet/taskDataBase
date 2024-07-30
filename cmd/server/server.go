package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/postgres"
	"flag"
	"fmt"
	"log"
	"net/http"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {

	// Обрабатываем флаги при запуске программы
	// go run server.go -typebd pg -loadbd true
	var typebd string = "pg"
	var loadbd bool = false

	flag.StringVar(&typebd, "typebd", "pg", "DataBase: pg-PostgreSQL, mem-memdb(map), mongo-MongoDB")
	flag.BoolVar(&loadbd, "loadbd", false, "Load data from json file: false/true")
	flag.Parse()

	fmt.Println(typebd)
	fmt.Println(loadbd)

	// Создаём объекты баз данных.
	// Реляционная БД PostgreSQL.
	db_pg, err := postgres.New("postgres://postgres:root@localhost:5432/prgDbStorage")
	if err != nil {
		log.Fatal(err)
	}

	// БД в памяти.
	db_mem, err := postgres.New("postgres://postgres:root@localhost:5432/prgDbStorage") //err := memdb.New()
	if err != nil {
		log.Fatal(err)
	}

	// Реляционная БД MongoDB.
	db_mongo, err := postgres.New("postgres://postgres:root@localhost:5432/prgDbStorage")
	if err != nil {
		log.Fatal(err)
	}

	defer db_pg.Close()
	defer db_mem.Close()
	defer db_mongo.Close()

	// Создаём объект сервера.
	var srv server

	// Инициализируем хранилище сервера конкретной БД.
	switch typebd {
	case "pg":
		srv.db = db_pg
	case "mem":
		srv.db = db_mem
	case "mongo":
		srv.db = db_mongo
	default:
		srv.db = db_pg
	}

	// Загружаем данные в БД при старте из файлов, если есть необходимость.
	if loadbd {
		err = srv.db.InsertInitDataFromFileAuthors(storage.AuthorsDb)
		if err != nil {
			log.Fatal(err)
		}

		err = srv.db.InsertInitDataFromFilePosts(storage.PostsDb)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	fmt.Println("Запуск веб-сервера на http://127.0.0.1:8080 ...")
	http.ListenAndServe(":8080", srv.api.Router())
}
