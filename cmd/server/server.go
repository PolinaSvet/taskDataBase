package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/memdb"
	"GoNews/pkg/storage/mongo"
	"GoNews/pkg/storage/postgres"
	"GoNews/pkg/storage/redis"

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
	// go run server.go -typebd pg -loadbd yes
	var typebd string
	var loadbd string

	flag.StringVar(&typebd, "typebd", "mem", "DataBase: pg-PostgreSQL, mem-memdb(map), mongo-MongoDB, redis-Redis")
	flag.StringVar(&loadbd, "loadbd", "yes", "Load data from json file: no/yes")
	flag.Parse()

	fmt.Println("flags: type bd->", typebd, "; preload data->", loadbd)

	// Создаём объект сервера.
	var srv server
	// Создаём объекты баз данных.
	// Инициализируем хранилище сервера конкретной БД.
	switch typebd {
	case "pg":
		// Реляционная БД PostgreSQL.
		db_pg, err := postgres.New("postgres://postgres:root@localhost:5432/prgDbStorage")
		if err != nil {
			log.Fatal(err)
		}
		srv.db = db_pg

	case "mem":
		// Не реляционная БД в памяти.
		db_mem, err := memdb.New()
		if err != nil {
			log.Fatal(err)
		}
		srv.db = db_mem

	case "mongo":
		// Не реляционная БД MongoDB.
		db_mongo, err := mongo.New("mongodb://localhost:27017/")
		if err != nil {
			log.Fatal(err)
		}
		srv.db = db_mongo

	case "redis":
		// Не реляционная БД Redis.
		db_redis, err := redis.New("localhost:6379", "", 0)
		if err != nil {
			log.Fatal(err)
		}
		srv.db = db_redis

	default:
		// Не реляционная БД в памяти.
		db_mem, err := memdb.New()
		if err != nil {
			log.Fatal(err)
		}
		srv.db = db_mem
	}

	defer srv.db.Close()

	// Загружаем данные в БД при старте из файлов, если есть необходимость.
	if loadbd == "yes" {
		err := srv.db.InsertInitDataFromFileAuthors(storage.AuthorsDb)
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
