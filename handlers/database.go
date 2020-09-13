package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DBStorage struct {
	log *log.Logger
	db  *sql.DB
}

const (
	TABLE_NAME = "URLdata"
	URL        = "url"
	SURL       = "short-url"
	ID         = "id"
)

func SetupDB() (*DBStorage, error) {

	dbLogger := log.New(os.Stdout, "DBLog ", log.LstdFlags)
	dbLogger.Println("[TRACE] Setting up database")

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		dbLogger.Println("[ERROR] connecting db", err)
		return nil, err
	}

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id serial NOT NULL, url VARCHAR, short_url VARCHAR)", TABLE_NAME))
	if err != nil {
		dbLogger.Printf("[ERROR] creating table %s. Reason: %s", TABLE_NAME, err)
		return nil, err
	}

	dbLogger.Println("[TRACE] Database successfully created")

	return &DBStorage{
		log: dbLogger,
		db:  db,
	}, nil
}

func (ds *DBStorage) InsertData(data *Data) error {

	err := ds.db.QueryRow(fmt.Sprintf("SELECT MAX(id) + 1 FROM %s", TABLE_NAME)).Scan(&data.ID)

	/// Generate url if custom url not added
	if data.ShortURL == "" {
		EncodeBase62(data)
	}

	_, err = ds.db.Query(fmt.Sprintf("INSERT INTO %s VALUES (%d, '%s', '%s')", TABLE_NAME, data.ID, data.URL, data.ShortURL))

	if err != nil {
		log.Printf("[ERROR] Unable to insert data. Reason: %s", err)
		return err
	}

	return nil
}

func (ds *DBStorage) SelectData(data *Data, search string) (bool, error) {

	var query string
	if search == URL {
		log.Printf("[DEBUG] Searching for url %s", data.URL)
		query = fmt.Sprintf("SELECT id, url, short_url FROM %s WHERE url = '%s'", TABLE_NAME, data.URL)
	} else if search == ID {
		log.Printf("[DEBUG] Searching for id %d", data.ID)
		query = fmt.Sprintf("SELECT id, url, short_url FROM %s WHERE id = %d", TABLE_NAME, data.ID)
	} else if search == SURL {
		log.Printf("[DEBUG] Searching for url %s", data.ShortURL)
		query = fmt.Sprintf("SELECT id, url, short_url FROM %s WHERE short_url = '%s'", TABLE_NAME, data.ShortURL)
	}

	rows, err := ds.db.Query(query)

	if err != nil {
		return false, err
	}

	if rows.Next() {
		err = rows.Scan(&data.ID, &data.URL, &data.ShortURL)
	} else {
		return false, nil
	}

	if err != nil {
		return false, nil
	}

	return true, err
}
