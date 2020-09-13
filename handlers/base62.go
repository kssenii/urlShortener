package handlers

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length   = int64(len(alphabet))
)

func EncodeBase62(data *Data) {

	if data.ID == 0 {
		data.ShortURL = string(alphabet[0])
	}

	for n := data.ID; n > 0; n = n / length {
		data.ShortURL = string(alphabet[n%length]) + data.ShortURL
	}
}
