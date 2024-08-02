package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName      = "test"    // имя учебной БД
	collectionAuthors = "authors" // имя коллекции в учебной БД
	collectionPosts   = "posts"   // имя коллекции в учебной БД
)

// Хранилище данных.
type Store struct {
	db *mongo.Client
}

// Конструктор объекта хранилища.
func New(constr string) (*Store, error) {
	// подключение к СУБД MongoDB
	mongoOpts := options.Client().ApplyURI(constr)
	db, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, err
	}
	// проверка связи с БД
	err = db.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	s := Store{
		db: db,
	}

	return &s, nil
}

func (s *Store) Close() {
	s.db.Disconnect(context.Background())
}

func (s *Store) GetInform() string {
	return "MongoDB"
}

// Author - автор.
func (s *Store) Authors() ([]storage.Author, error) {

	var authors []storage.Author

	collection := s.db.Database(databaseName).Collection(collectionAuthors)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &authors); err != nil {
		return nil, err
	}

	return authors, nil
}

func (s *Store) AddAuthor(author storage.Author) (int64, error) {

	collection := s.db.Database(databaseName).Collection(collectionAuthors)
	_, err := collection.InsertOne(context.Background(), author)
	if err != nil {
		return 0, err
	}

	return author.ID, nil
}

func (s *Store) UpdateAuthor(author storage.Author) (int64, error) {
	doc := bson.M{
		"name": author.Name,
	}
	id_doc := bson.M{"_id": author.ID}

	collection := s.db.Database(databaseName).Collection(collectionAuthors)
	result, err := collection.UpdateOne(context.Background(), id_doc, bson.M{"$set": doc})
	if err != nil {
		return 0, err
	}
	if result.MatchedCount == 0 {
		return 0, fmt.Errorf("UPDATE Id: %v not exist", id_doc)
	}

	return author.ID, nil
}

func (s *Store) DeleteAuthor(author storage.Author) (int64, error) {
	id_doc := bson.M{"_id": author.ID}

	collection := s.db.Database(databaseName).Collection(collectionAuthors)
	result, err := collection.DeleteOne(context.Background(), id_doc)
	if err != nil {
		return 0, err
	}
	if result.DeletedCount == 0 {
		return 0, fmt.Errorf("DELETE Id: %v not exist", id_doc)
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

	var documents []interface{}
	for _, author := range authors {
		documents = append(documents, author)
	}

	collection := s.db.Database(databaseName).Collection(collectionAuthors)
	_, err = collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	_, err = collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	return nil
}

// Post - публикация.
func (s *Store) Posts() ([]storage.Post, error) {

	pipeline := bson.A{
		bson.M{
			"$lookup": bson.M{
				"from":         "authors",
				"localField":   "author_id",
				"foreignField": "_id",
				"as":           "author",
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path":                       "$author",
				"preserveNullAndEmptyArrays": true,
			},
		},
		bson.M{
			"$project": bson.M{

				"_id":         1,
				"author_id":   1,
				"author_name": bson.M{"$ifNull": []interface{}{"$author.name", "None"}},
				"title":       1,
				"content":     1,
				"created_at":  1,
				"created_at_txt": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$gt": []interface{}{"$created_at", nil}},
						"then": bson.M{"$dateToString": bson.M{"format": "%d.%m.%Y %H:%M:%S", "date": bson.M{"$toDate": bson.M{"$multiply": []interface{}{"$created_at", 1}}}}},
						"else": "",
					},
				},
				"published_at": 1,
				"published_at_txt": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$gt": []interface{}{"$published_at", nil}},
						"then": bson.M{"$dateToString": bson.M{"format": "%d.%m.%Y %H:%M:%S", "date": bson.M{"$toDate": bson.M{"$multiply": []interface{}{"$published_at", 1}}}}},
						"else": "",
					},
				},
			},
		},
	}
	var posts []storage.Post

	collection := s.db.Database(databaseName).Collection(collectionPosts)
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &posts); err != nil {
		return nil, err
	}

	return posts, nil

}

func (s *Store) AddPost(post storage.Post) (int64, error) {

	collection := s.db.Database(databaseName).Collection(collectionPosts)
	_, err := collection.InsertOne(context.Background(), post)
	if err != nil {
		return 0, err
	}

	return post.ID, nil
}

func (s *Store) UpdatePost(post storage.Post) (int64, error) {

	doc := bson.M{
		"author_id":    post.AuthorID,
		"title":        post.Title,
		"content":      post.Content,
		"created_at":   post.CreatedAt,
		"published_at": post.PublishedAt,
	}
	id_doc := bson.M{"_id": post.ID}

	collection := s.db.Database(databaseName).Collection(collectionPosts)
	result, err := collection.UpdateOne(context.Background(), id_doc, bson.M{"$set": doc})
	//fmt.Printf("%#v\n", result)
	if err != nil {
		return 0, err
	}
	if result.MatchedCount == 0 {
		return 0, fmt.Errorf("UPDATE Id: %v not exist", id_doc)
	}

	return post.ID, nil
}

func (s *Store) DeletePost(post storage.Post) (int64, error) {

	id_doc := bson.M{"_id": post.ID}

	collection := s.db.Database(databaseName).Collection(collectionPosts)
	result, err := collection.DeleteOne(context.Background(), id_doc)
	if err != nil {
		return 0, err
	}
	if result.DeletedCount == 0 {
		return 0, fmt.Errorf("DELETE Id: %v not exist", id_doc)
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

	var documents []interface{}
	for _, post := range posts {

		doc := bson.M{
			"_id":          post.ID,
			"author_id":    post.AuthorID,
			"title":        post.Title,
			"content":      post.Content,
			"created_at":   post.CreatedAt,
			"published_at": post.PublishedAt,
		}

		documents = append(documents, doc)
	}

	collection := s.db.Database(databaseName).Collection(collectionPosts)
	_, err = collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	_, err = collection.InsertMany(context.Background(), documents)
	if err != nil {
		return err
	}

	return nil
}
