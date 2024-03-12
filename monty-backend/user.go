package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Created  time.Time `json:"_"`
	Admin    bool      `json:"_"`
	Password []byte    `json:"_"`
	Salt     []byte    `json:"_"`
	Modified time.Time `json:"_"`
	Accessed time.Time `json:"_"`
	Status   int       `json:"_"`
}

type UserCreate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func hashPassword(password string) (key []byte, salt []byte) {
	salt = make([]byte, 16)
	_, err := rand.Read(salt)

	if err != nil {
		panic(err)
	}

	key = argon2.IDKey([]byte(password), salt, ARGON_TIME, ARGON_MEMORY, ARGON_THREADS, ARGON_LENGTH)

	return key, salt
}

func (user *User) VerifyPassword(password string) bool {
	key := argon2.IDKey([]byte(password), user.Salt, ARGON_TIME, ARGON_MEMORY, ARGON_THREADS, ARGON_LENGTH)

	return string(key) == string(user.Password)
}

func NewUser(user *UserCreate) User {
	password, salt := hashPassword(user.Password)

	return User{
		ID:       uuid.New(),
		Username: user.Username,
		Email:    user.Email,
		Password: password,
		Salt:     salt,
	}
}

func CreateUserHandler(context *Context, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var newUser UserCreate
	err := decoder.Decode(&newUser)

	if err != nil || newUser.Username == "" || newUser.Email == "" || newUser.Password == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(newUser.Password) < PASSWORD_MIN_LENGTH {
		http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	if len(newUser.Username) < USERNAME_MIN_LENGTH || len(newUser.Username) > USERNAME_MAX_LENGTH {
		http.Error(w, "Username must be between 3 and 32 characters", http.StatusBadRequest)
		return
	}

	user := NewUser(&newUser)

	err = context.DB.CreateUser(&user)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func LoginUserHandler(context *Context, w http.ResponseWriter, r *http.Request) {
	var user UserLogin

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dbUser, err := context.DB.VerifyUser(user.Username, user.Password)

	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	jwt, err := generateJWT(&dbUser, context.Config.JWTSecret)

	if err != nil {
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "JWTToken",
		Value:    jwt,
		Expires:  time.Now().Add(JWT_EXP),
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(w, &cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dbUser)
}
