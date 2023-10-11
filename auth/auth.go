package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"com.serve_volt/kv"
	"com.serve_volt/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var store = kv.NewKVStore[string, any]()
var secret string

type UserClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func init() {
	loadEnv()
	secret = os.Getenv("SECRET_KEY")
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading the .env file")
	}
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if has := store.Has(user.Username); has {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	errs := store.Put(user.Username, user)

	if errs != nil {
		log.Fatalln(errs)
	}

	w.WriteHeader(http.StatusCreated)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storedUser, exists := store.Get(user.Username)
	if exists != nil || storedUser != user {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, errt := token.SignedString(secret)
	if errt != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{} ,func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*UserClaims);
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(AddUserToContext(r.Context(), claims.Username))

		next.ServeHTTP(w, r)

	})
}


func AddUserToContext(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, "username", username)
} // aditionally if need data from token can access adding in context.WithValue


func GetUserFromContext(ctx context.Context) string {
	username, _ := ctx.Value("username").(string)
	return username
} //
