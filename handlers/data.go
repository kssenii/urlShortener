package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

type Data struct {
	ID       int64  `json:"-"`
	URL      string `json:"url" validate:"url-validate"`
	ShortURL string `json:"short,omitempty"`
}

func (data *Data) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(data)
}

func (data *Data) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(data)
}

func (data *Data) Validate() error {

	validate := validator.New()
	validate.RegisterValidation("url-validate", IsReachableURL)
	return validate.Struct(data)
}

func (data *Data) ValidateShort() error {

	validate := validator.New()
	validate.RegisterValidation("url-validate", checkShortURL)
	return validate.Struct(data)
}

func checkShortURL(fl validator.FieldLevel) bool {

	url := fl.Field().String()
	if strings.HasPrefix(url, fmt.Sprintf("http://localhost:%s/", os.Getenv("APP_PORT"))) {
		return true
	}

	return true
}

func IsReachableURL(fl validator.FieldLevel) bool {

	url := fl.Field().String()
	response, errors := http.Get(url)

	if errors != nil {
		return false
	}

	if response.StatusCode == 200 {
		return true
	}

	return false
}
