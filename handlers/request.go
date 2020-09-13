package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
)

type Request struct {
	log       *log.Logger
	dbHandler *DBStorage
	idCounter int64
}

type Error struct {
	info   error
	status int
}

func NewRequest(dh *DBStorage) *Request {
	return &Request{
		log:       log.New(os.Stdout, "RequestLog ", log.LstdFlags),
		dbHandler: dh,
		idCounter: 0,
	}
}

func (request *Request) EncodeURL(rw http.ResponseWriter, r *http.Request) {

	data := r.Context().Value(KeyData{}).(*Data)

	if !strings.Contains(data.URL, "http") {
		data.URL = "http://" + data.URL
	}

	exists, err := request.dbHandler.SelectData(data, URL)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		request.log.Println("URL encoding already exists")
	} else {
		err = request.dbHandler.InsertData(data)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	data.URL = data.ShortURL
	errRet := data.ToJSON(rw)
	if errRet != nil {
		request.log.Println("[ERROR] Unable to parse JSON. Reason: ", errRet)
		http.Error(rw, "Unable to parse data", http.StatusInternalServerError)
		return
	}
}

func (request *Request) DecodeURL(rw http.ResponseWriter, r *http.Request) {

	data := r.Context().Value(KeyData{}).(*Data)
	var err error
	var exists bool

	data.ShortURL = strings.TrimPrefix(data.URL, "http://localhost:9090/")
	data.URL = ""

	err = DecodeBase62(data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	exists, err = request.dbHandler.SelectData(data, ID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(rw, "URL not found", http.StatusInternalServerError)
		return
	}

	err = data.ToJSON(rw)
	if err != nil {
		request.log.Println("[ERROR] Unable to parse JSON. Reason: ", err)
		http.Error(rw, "Unable to add data", http.StatusInternalServerError)
		return
	}
}

func (request *Request) Redirect(rw http.ResponseWriter, r *http.Request) {

	var data Data
	data.ShortURL = r.URL.Path[len("/"):]

	err := DecodeBase62(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	exists, err := request.dbHandler.SelectData(&data, ID)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(rw, "URL not found", http.StatusInternalServerError)
		return
	}

	request.log.Println("Redircting to", data.URL)
	http.Redirect(rw, r, data.URL, http.StatusMovedPermanently)
}

type KeyData struct{}

func (request Request) MiddlewareValidateData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		data := &Data{}
		err := data.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		task := r.URL.Path[len("/"):]
		if task == "decode" {
			err = data.ValidateShort()
			if err != nil {
				http.Error(rw, "Unable to validate data", http.StatusBadRequest)
				return
			}
		} else {
			err = data.Validate()
			if err != nil {
				http.Error(rw, "Unable to validate data", http.StatusBadRequest)
				return
			}
		}

		context := context.WithValue(r.Context(), KeyData{}, data)
		r = r.WithContext(context)

		next.ServeHTTP(rw, r)
	})
}
