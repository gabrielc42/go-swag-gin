package helper

var jwtKey = []byte("my_secret_key")

type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateToken(_id string) (string, int64, error) {}
