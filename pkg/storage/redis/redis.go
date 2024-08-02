package redis

import (
	"GoNews/pkg/storage"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	collectionAuthors = "authors" // имя коллекции в учебной БД
	collectionPosts   = "posts"   // имя коллекции в учебной БД
)

// Хранилище данных.
type Store struct {
	db *redis.Client
}

func (s *Store) GetInform() string {
	return "Redis"
}

// Конструктор объекта хранилища.
func New(constr string, password string, number int) (*Store, error) {

	db := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	s := Store{
		db: db,
	}

	fmt.Println("Loaded bd: ", s.GetInform())

	return &s, nil
}

func (s *Store) Close() {
}

func (s *Store) existsKey(key string) bool {

	exists, err := s.db.Exists(context.Background(), key).Result()
	if err != nil {
		return false
	}
	if exists == 0 {
		return false
	}

	return true
}

func (s *Store) getNameAuthorsById(post storage.Post) (string, error) {

	key := fmt.Sprintf("%s:%d", collectionAuthors, post.AuthorID)

	val, err := s.db.Get(context.Background(), key).Result()
	if err != nil {
		return fmt.Sprintf("AuthorID: %v not exist", post.AuthorID), err
	}
	var author storage.Author
	err = json.Unmarshal([]byte(val), &author)
	if err != nil {
		return fmt.Sprintf("AuthorID: %v not exist", post.AuthorID), err
	}

	return author.Name, nil
}

// Author - автор.
func (s *Store) Authors() ([]storage.Author, error) {
	var authors []storage.Author

	keyPattern := fmt.Sprintf("%s:*", collectionAuthors)
	keys, err := s.db.Keys(context.Background(), keyPattern).Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		val, err := s.db.Get(context.Background(), key).Result()
		if err != nil {
			return nil, err
		}
		var author storage.Author
		err = json.Unmarshal([]byte(val), &author)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil
}

func (s *Store) AddAuthor(author storage.Author) (int64, error) {

	key := fmt.Sprintf("%s:%d", collectionAuthors, author.ID)

	if s.existsKey(key) {
		return 0, fmt.Errorf("INSERT Id: %v exist", key)
	}

	val, err := json.Marshal(author)
	if err != nil {
		return 0, err
	}
	err = s.db.Set(context.Background(), key, string(val), 0).Err()
	if err != nil {
		return 0, err
	}

	return author.ID, nil
}

func (s *Store) UpdateAuthor(author storage.Author) (int64, error) {

	key := fmt.Sprintf("%s:%d", collectionAuthors, author.ID)

	if !s.existsKey(key) {
		return 0, fmt.Errorf("UPDATE Id: %v not exist", key)
	}

	val, err := json.Marshal(author)
	if err != nil {
		return 0, err
	}
	err = s.db.Set(context.Background(), key, string(val), 0).Err()
	if err != nil {
		return 0, err
	}

	return author.ID, nil
}

func (s *Store) DeleteAuthor(author storage.Author) (int64, error) {

	key := fmt.Sprintf("%s:%d", collectionAuthors, author.ID)

	if !s.existsKey(key) {
		return 0, fmt.Errorf("DELETE Id: %v not exist", key)
	}

	err := s.db.Del(context.Background(), key).Err()
	if err != nil {
		return 0, err
	}

	return author.ID, nil
}

func (s *Store) InsertInitDataFromFileAuthors(filename string) error {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var authors []storage.Author
	err = json.Unmarshal(data, &authors)
	if err != nil {
		return err
	}

	for _, author := range authors {
		key := fmt.Sprintf("%s:%d", collectionAuthors, author.ID)
		val, err := json.Marshal(author)
		if err != nil {
			return err
		}
		err = s.db.Set(context.Background(), key, string(val), 0).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

// Post - публикация.
func (s *Store) Posts() ([]storage.Post, error) {
	var posts []storage.Post

	keyPattern := fmt.Sprintf("%s:*", collectionPosts)
	keys, err := s.db.Keys(context.Background(), keyPattern).Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		val, err := s.db.Get(context.Background(), key).Result()
		if err != nil {
			return nil, err
		}
		var post storage.Post
		err = json.Unmarshal([]byte(val), &post)
		if err != nil {
			return nil, err
		}

		authorName, _ := s.getNameAuthorsById(post)
		post.AuthorName = authorName
		post.CreatedAtTxt = time.Unix(int64(post.CreatedAt)/1000, 0).Format("2006-01-02 15:04:05.000")
		post.PublishedAtTxt = time.Unix(int64(post.PublishedAt)/1000, 0).Format("2006-01-02 15:04:05.000")

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *Store) AddPost(post storage.Post) (int64, error) {

	key := fmt.Sprintf("%s:%d", collectionPosts, post.ID)

	if s.existsKey(key) {
		return 0, fmt.Errorf("INSERT Id: %v exist", key)
	}

	val, err := json.Marshal(post)
	if err != nil {
		return 0, err
	}
	err = s.db.Set(context.Background(), key, string(val), 0).Err()
	if err != nil {
		return 0, err
	}

	return post.ID, nil
}

func (s *Store) UpdatePost(post storage.Post) (int64, error) {

	key := fmt.Sprintf("%s:%d", collectionPosts, post.ID)

	if !s.existsKey(key) {
		return 0, fmt.Errorf("UPDATE Id: %v not exist", key)
	}

	val, err := json.Marshal(post)
	if err != nil {
		return 0, err
	}
	err = s.db.Set(context.Background(), key, string(val), 0).Err()
	if err != nil {
		return 0, err
	}

	return post.ID, nil
}

func (s *Store) DeletePost(post storage.Post) (int64, error) {

	key := fmt.Sprintf("%s:%d", collectionPosts, post.ID)

	if !s.existsKey(key) {
		return 0, fmt.Errorf("DELETE Id: %v not exist", key)
	}

	err := s.db.Del(context.Background(), key).Err()
	if err != nil {
		return 0, err
	}

	return post.ID, nil
}

func (s *Store) InsertInitDataFromFilePosts(filename string) error {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var posts []storage.Post
	err = json.Unmarshal(data, &posts)
	if err != nil {
		return err
	}

	for _, post := range posts {
		key := fmt.Sprintf("%s:%d", collectionPosts, post.ID)
		val, err := json.Marshal(post)
		if err != nil {
			return err
		}
		err = s.db.Set(context.Background(), key, string(val), 0).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
