package memdb

import (
	"GoNews/pkg/storage"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

// Хранилище данных.
type Store struct {
	AuthorsDB map[int64]storage.Author
	PostsDB   map[int64]storage.Post
}

func (s *Store) GetInform() string {
	return "MemDB"
}

// Конструктор объекта хранилища.
func New() (*Store, error) {
	s := Store{
		AuthorsDB: map[int64]storage.Author{},
		PostsDB:   map[int64]storage.Post{},
	}

	fmt.Println("Loaded bd: ", s.GetInform())

	return &s, nil
}

func (s *Store) Close() {
}

// Author - автор.
func (s *Store) Authors() ([]storage.Author, error) {
	var data []storage.Author

	for _, v := range s.AuthorsDB {
		data = append(data, v)
	}
	return data, nil
}

func (s *Store) AddAuthor(author storage.Author) (int64, error) {
	if _, ok := s.AuthorsDB[author.ID]; ok {
		return 0, fmt.Errorf("Id: %v already exist", author.ID)
	} else {
		s.AuthorsDB[author.ID] = author
		return author.ID, nil
	}
}

func (s *Store) UpdateAuthor(author storage.Author) (int64, error) {
	if _, ok := s.AuthorsDB[author.ID]; !ok {
		return 0, fmt.Errorf("Id: %v not exist", author.ID)
	} else {
		s.AuthorsDB[author.ID] = author
		return author.ID, nil
	}
}

func (s *Store) DeleteAuthor(author storage.Author) (int64, error) {
	if _, ok := s.AuthorsDB[author.ID]; !ok {
		return 0, fmt.Errorf("Id: %v not exist", author.ID)
	} else {
		delete(s.AuthorsDB, author.ID)
		return author.ID, nil
	}
}

func (s *Store) InsertInitDataFromFileAuthors(filename string) error {

	file, _ := ioutil.ReadFile(filename)
	data := []storage.Author{}
	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		return err
	}

	for i := 0; i < len(data); i++ {
		s.AuthorsDB[data[i].ID] = data[i]
	}

	return nil
}

// Post - публикация.
func (s *Store) Posts() ([]storage.Post, error) {
	var data []storage.Post

	for _, v := range s.PostsDB {

		if _, ok := s.AuthorsDB[v.AuthorID]; ok {
			v.AuthorName = s.AuthorsDB[v.AuthorID].Name
		}

		dt_CreatedAt := time.Unix(v.CreatedAt/1000, 0)
		v.CreatedAtTxt = dt_CreatedAt.Format("2006-01-02 15:04:05.000")

		dt_PublishedAtTxt := time.Unix(v.PublishedAt/1000, 0)
		v.PublishedAtTxt = dt_PublishedAtTxt.Format("2006-01-02 15:04:05.000")

		data = append(data, v)
	}
	return data, nil
}

func (s *Store) AddPost(post storage.Post) (int64, error) {
	if _, ok := s.PostsDB[post.ID]; ok {
		return 0, fmt.Errorf("Id: %v already exist", post.ID)
	} else {
		if _, ok := s.AuthorsDB[post.AuthorID]; ok {
			s.PostsDB[post.ID] = post
			return post.ID, nil

		} else {
			return 0, fmt.Errorf("Author with id: %v not exist", post.AuthorID)
		}
	}
}

func (s *Store) UpdatePost(post storage.Post) (int64, error) {
	if _, ok := s.PostsDB[post.ID]; !ok {
		return 0, fmt.Errorf("Id: %v not exist", post.ID)
	} else {
		if _, ok := s.AuthorsDB[post.AuthorID]; ok {
			s.PostsDB[post.ID] = post
			return post.ID, nil

		} else {
			return 0, fmt.Errorf("Author with id: %v not exist", post.AuthorID)
		}
	}
}

func (s *Store) DeletePost(post storage.Post) (int64, error) {
	if _, ok := s.PostsDB[post.ID]; !ok {
		return 0, fmt.Errorf("Id: %v not exist", post.ID)
	} else {
		delete(s.PostsDB, post.ID)
		return post.ID, nil
	}
}

func (s *Store) InsertInitDataFromFilePosts(filename string) error {
	file, _ := ioutil.ReadFile(filename)
	data := []storage.Post{}
	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		return err
	}

	for i := 0; i < len(data); i++ {
		s.PostsDB[data[i].ID] = data[i]
	}

	return nil
}
