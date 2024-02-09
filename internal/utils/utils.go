package utils

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/h3th-IV/mysticMerch/internal/models"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// Logger represents a custom logger.
type Logger struct {
	InfoLogger *log.Logger
	ErrLogger  *log.Logger
}

// NewLogger creates a new instance of Logger with customizable log formatting.
func NewLogger(infoOutput, errorOutput *os.File) *Logger {
	return &Logger{
		InfoLogger: log.New(infoOutput, "INFO\t", log.Ldate|log.Ltime),
		ErrLogger:  log.New(errorOutput, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// LogInfo logs an informational message.
func (l *Logger) LogInfo(message string) {
	l.InfoLogger.Println(message)
}

// LogError logs an error message.
func (l *Logger) LogError(message string) {
	l.ErrLogger.Println(message)
}

// Load env variables
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

// Middleware to log requests to the server ##
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// logger := NewLogger(os.Stdout, os.Stderr)
		ReplaceLogger.Info((fmt.Sprintf("%v - %v %v %v", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())))
		next.ServeHTTP(w, r)
	})
}

// func HashPassword
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", nil
	}
	return string(hash), nil
}

var (
	ErrNoRecord = errors.New("err: no matching record found")

	ErrInvalidCredentials = errors.New("err: invalid credentials")

	ErrExsistingCrednetials       = errors.New("err: dupliacte Credentials")
	MySQLErr                      *mysql.MySQLError
	ErrMismatchedCryptAndPassword = errors.New("err: Paswword does not Match registered password")
)

// Middleware to recover panic ##
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				//if panic close connection
				w.Header().Set("Connection", "Close")
				//write internal server error
				ServerError(w, "Connection Closed inabruptly", fmt.Errorf("%v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type mapKey string

const (
	UserIDkey mapKey = "user_id"
)

func GenerateToken(user *models.User, expiry time.Duration, issuer, secret string) (string, error) {
	//set expiry date
	bestBefore := time.Now().Add(expiry)

	//ceate jwt tkeo with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		"epx":     bestBefore.Unix(),
		"iss":     issuer,
	})

	//generate token str and sign with seceret key
	JWToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return JWToken, nil
}

func AdminToken(user *models.User, expiry time.Duration, issuer, secret string) (string, error) {
	bestBefore := time.Now().Add(time.Hour * 10)

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"user": user.UserID,
		"exp":  bestBefore,
		"iss":  issuer,
	})

	ADMINToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return ADMINToken, nil
}

// Middleware to Auth specific routes
func JWTAuthRoutes(next http.Handler, secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//get JWToken from request
		JWToken := r.Header.Get("Authorization")
		token, err := jwt.Parse(JWToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized Operation", http.StatusUnauthorized)
			return
		}
		tokenClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid Token claims", http.StatusBadRequest)
			return
		}

		userID, ok := tokenClaims["user_id"]
		if !ok {
			http.Error(w, "User is not Authorized", http.StatusBadRequest)
			return
		}

		//store user_id in context
		ctx := context.WithValue(r.Context(), UserIDkey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthRoute(next http.Handler) http.Handler {
	LoadEnv()
	return JWTAuthRoutes(next, os.Getenv("MYSTIC"))
}

// auth route for admin
func AdminRoute(next http.Handler) http.Handler {
	LoadEnv()
	return JWTAuthRoutes(next, os.Getenv("MYTH"))
}

// used for all internal server Error
func ServerError(w http.ResponseWriter, errMsg string, err error) {
	errTrace := fmt.Sprintf("%v\n%v", err.Error(), debug.Stack())
	ReplaceLogger.Error(errTrace)
	http.Error(w, errMsg, http.StatusInternalServerError)
}

func ValidateSignUpDetails(details []models.ValidAta) bool {
	email := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	password := regexp.MustCompile("^[a-zA-Z0-9!@#$%^&*()-_=+{}[]|;:'\",.<>?/`~]{8,15}$")
	fName := regexp.MustCompile("^[A-Za-z]+$")
	lName := regexp.MustCompile("^[A-Za-z]+$")

	for i := 0; i < len(details); i++ {
		switch details[i].Validator {
		case "email":
			if !email.MatchString(details[i].Value) {
				return false
			}
		case "password":
			if !password.MatchString(details[i].Value) {
				return false
			}
		case "f_name":
			if !fName.MatchString(details[i].Value) {
				return false
			}
		case "lName":
			if !lName.MatchString(details[i].Value) {
				return false
			}
		}
	}
	return true
}

func GenerateUUID(elemenType string) (string, error) {
	//generate new uuuid
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	//convert string and add prefix
	uuid := id.String()
	switch elemenType {
	case "user":
		uuid = "usr" + uuid
	case "product":
		uuid = "prd" + uuid
	}
	return uuid, nil
}

// EncryptPass encrypts password using AES.
func EncryptPass(password []byte) (string, error) {
	if err := LoadEnv(); err != nil {
		return "", err
	}
	key := []byte(os.Getenv("HADESKEY"))
	//create aes block with provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	//make a cuipher text to store encrypted passwrd
	cipherText := make([]byte, aes.BlockSize)
	iv := cipherText[:aes.BlockSize] //prepend initialization vector to cipher slice
	//initialization vector for randomness
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	//use cipher to create new ecncrypter stream  used to ecrypt plain text data
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], password) //use stream to encrypt password

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptPass(cipherText string) (string, error) {
	if err := LoadEnv(); err != nil {
		return "", err
	}
	key := []byte(os.Getenv("HADESKEY"))
	//create a new block with key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// decode cipherText.
	decipherText, err := base64.RawStdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", nil
	}

	// pop initialization vector.
	iv := decipherText[:aes.BlockSize]
	decipherText = decipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decipherText, decipherText)

	return string(decipherText), nil
}

func CompareCryptedAndPassword(password string, user *models.User) error {
	decrypted, err := DecryptPass(user.Password)
	if err != nil {
		return err
	}

	//mitigate time @++cks constant time compare
	if subtle.ConstantTimeCompare([]byte(decrypted), []byte(password)) != 1 {
		//password deos not match
		return ErrMismatchedCryptAndPassword
	}
	return nil
}
