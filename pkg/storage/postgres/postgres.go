package postgres

import (
	"GoNews/pkg/storage"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

func (s *Store) GetInform() string {
	return "PostgreSQL"
}

// Конструктор объекта хранилища.
func New(constr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: db,
	}

	fmt.Println("Loaded bd: ", s.GetInform())

	return &s, nil
}

func (s *Store) Close() {
	s.db.Close()
}

func structToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap)
	return
}

// Author - автор.
func (s *Store) Authors() ([]storage.Author, error) {
	rows, err := s.db.Query(context.Background(), `SELECT * FROM authors_func_view('{}');`)

	if err != nil {
		return nil, err
	}

	var authors []storage.Author

	for rows.Next() {
		var t storage.Author
		err = rows.Scan(
			&t.ID,
			&t.Name,
		)
		if err != nil {
			return nil, err
		}

		authors = append(authors, t)
	}
	return authors, rows.Err()
}

func (s *Store) AddAuthor(author storage.Author) (int64, error) {
	jsonRequest, err := structToMap(author)
	if err != nil {
		return 0, err
	}

	var jsonResponse storage.SqlResponse
	err = s.db.QueryRow(context.Background(), `SELECT * FROM authors_func_insert($1);`, jsonRequest).Scan(&jsonResponse)
	if err != nil {
		return 0, err
	}
	if jsonResponse.Err != "" {
		return 0, fmt.Errorf(jsonResponse.Err)
	}
	return jsonResponse.ID, nil
}

func (s *Store) UpdateAuthor(author storage.Author) (int64, error) {
	jsonRequest, err := structToMap(author)
	if err != nil {
		return 0, err
	}

	var jsonResponse storage.SqlResponse
	err = s.db.QueryRow(context.Background(), `SELECT * FROM authors_func_update($1);`, jsonRequest).Scan(&jsonResponse)
	if err != nil {
		return 0, err
	}
	if jsonResponse.Err != "" {
		return 0, fmt.Errorf(jsonResponse.Err)
	}
	return jsonResponse.ID, nil
}

func (s *Store) DeleteAuthor(author storage.Author) (int64, error) {
	jsonRequest, err := structToMap(author)
	if err != nil {
		return 0, err
	}

	var jsonResponse storage.SqlResponse
	err = s.db.QueryRow(context.Background(), `SELECT * FROM authors_func_delete($1);`, jsonRequest).Scan(&jsonResponse)
	if err != nil {
		return 0, err
	}
	if jsonResponse.Err != "" {
		return 0, fmt.Errorf(jsonResponse.Err)
	}
	return jsonResponse.ID, nil
}

func (s *Store) InsertInitDataFromFileAuthors(filename string) error {

	var jsonData []map[string]interface{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return err
	}

	for _, item := range jsonData {

		var jsonResponse storage.SqlResponse
		err = s.db.QueryRow(context.Background(), `SELECT * FROM authors_func_insert($1);`, item).Scan(&jsonResponse)
		if err != nil {
			return err
		}
		if jsonResponse.Err != "" {
			return fmt.Errorf(jsonResponse.Err)
		}

	}

	return nil
}

// Post - публикация.
func (s *Store) Posts() ([]storage.Post, error) {

	rows, err := s.db.Query(context.Background(), `SELECT * FROM posts_func_view('{}');`)

	if err != nil {
		return nil, err
	}

	var posts []storage.Post

	for rows.Next() {
		var t storage.Post
		err = rows.Scan(
			&t.ID,
			&t.AuthorID,
			&t.AuthorName,
			&t.Title,
			&t.Content,
			&t.CreatedAt,
			&t.CreatedAtTxt,
			&t.PublishedAt,
			&t.PublishedAtTxt,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, t)
	}
	return posts, rows.Err()
}

func (s *Store) AddPost(post storage.Post) (int64, error) {

	jsonRequest, err := structToMap(post)
	if err != nil {
		return 0, err
	}

	var jsonResponse storage.SqlResponse
	err = s.db.QueryRow(context.Background(), `SELECT * FROM posts_func_insert($1);`, jsonRequest).Scan(&jsonResponse)
	if err != nil {
		return 0, err
	}
	if jsonResponse.Err != "" {
		return 0, fmt.Errorf(jsonResponse.Err)
	}
	return jsonResponse.ID, nil
}

func (s *Store) UpdatePost(post storage.Post) (int64, error) {

	jsonRequest, err := structToMap(post)
	if err != nil {
		return 0, err
	}

	var jsonResponse storage.SqlResponse
	err = s.db.QueryRow(context.Background(), `SELECT * FROM posts_func_update($1);`, jsonRequest).Scan(&jsonResponse)
	if err != nil {
		return 0, err
	}
	if jsonResponse.Err != "" {
		return 0, fmt.Errorf(jsonResponse.Err)
	}
	return jsonResponse.ID, nil
}

func (s *Store) DeletePost(post storage.Post) (int64, error) {

	jsonRequest, err := structToMap(post)
	if err != nil {
		return 0, err
	}

	var jsonResponse storage.SqlResponse
	err = s.db.QueryRow(context.Background(), `SELECT * FROM posts_func_delete($1);`, jsonRequest).Scan(&jsonResponse)
	if err != nil {
		return 0, err
	}
	if jsonResponse.Err != "" {
		return 0, fmt.Errorf(jsonResponse.Err)
	}
	return jsonResponse.ID, nil
}

func (s *Store) InsertInitDataFromFilePosts(filename string) error {

	var jsonData []map[string]interface{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return err
	}

	for _, item := range jsonData {
		var jsonResponse storage.SqlResponse
		err = s.db.QueryRow(context.Background(), `SELECT * FROM posts_func_insert($1);`, item).Scan(&jsonResponse)
		if err != nil {
			return err
		}
		if jsonResponse.Err != "" {
			return fmt.Errorf(jsonResponse.Err)
		}
	}

	return nil
}
