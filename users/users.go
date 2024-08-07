package users

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type contextKey string

const UserContextKey contextKey = "user"

func GetUserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(UserContextKey).(User)
	if !ok {
		// Handle the case when the user is not found in the context
		return nil
	}
	return &user
}

func GetUserFromToken(tkn string) User {
	return User{UUID: uuid.New(), FirstName: "Gary", LastName: "Barnett", Email: "gary@thinkovation.com"}
}
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	fmt.Println(user.Email)

}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	fmt.Println(user.Email)

}

func getUserRights(userID uuid.UUID, entityID uuid.UUID) []string {
	fmt.Println(userID, entityID)
	return []string{}
}

// User is the struct to hold the user info. For safety, we expressly exclude the Password from the user struct
type User struct {
	UUID           uuid.UUID `json:"uuid"`
	OrganisationID uuid.UUID `json:"organisation_id"`
	Username       string    `json:"username"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func CreateUser(db *sql.DB, user User) error {
	query := `
        INSERT INTO users (uuid, organisation_id, username, first_name, last_name, email, phone)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	_, err := db.Exec(query, user.UUID, user.OrganisationID, user.Username, user.FirstName, user.LastName, user.Email)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}
func UpdateUserPwd(db *sql.DB, userID uuid.UUID, pwd string) error {
	query := `
        UPDATE users
        SET password = $2
        WHERE uuid = $1
    `
	_, err := db.Exec(query, userID, pwd)
	if err != nil {
		return fmt.Errorf("failed to update user password: %v", err)
	}
	return nil

}

func UpdateUser(db *sql.DB, user User) error {
	query := `
        UPDATE users
        SET organisation_id = $2, username = $3, first_name = $4, last_name = $5, email = $6
        WHERE uuid = $1
    `
	_, err := db.Exec(query, user.UUID, user.OrganisationID, user.Username, user.FirstName, user.LastName, user.Email)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

// DeleteUser deletes a user from the database by setting the deleted flag to true
func DeleteUser(db *sql.DB, uuid string) error {
	query := `
        UPDATE users
		SET deleted = true
        WHERE uuid = $1
    `
	_, err := db.Exec(query, uuid)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	// Generate a bcrypt hash of the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hash), nil
}

// ComparePassword compares a hashed password with a plain password
func ComparePassword(hashedPassword, password string) error {
	// Compare the hashed password with the plain password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("passwords do not match: %v", err)
	}
	return nil
}

// GetYserByUsername returns a user with a given username
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	query := `
        SELECT uuid, organisation_id, username, first_name, last_name, email, created_at, updated_at
        FROM users
        WHERE username = $1 AND deleted = false
    `
	var nullableFirstName sql.NullString
	var nullableLastName sql.NullString
	var nullablePhone sql.NullString

	user := &User{}
	err := db.QueryRow(query, username).Scan(&user.UUID, &user.OrganisationID, &user.Username, &nullableFirstName, &nullableLastName, &user.Email, &nullablePhone, &user.CreatedAt, &user.UpdatedAt)
	if nullableFirstName.Valid {
		user.FirstName = nullableFirstName.String
	}
	if nullableLastName.Valid {
		user.LastName = nullableLastName.String
	}
	if nullablePhone.Valid {
		user.Phone = nullablePhone.String
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by username: %v", err)
	}

	return user, nil
}

// GetUserByUUID returns a user with a given UUID
func GetUserByUUID(db *sql.DB, uuid uuid.UUID) (*User, error) {
	query := `
        SELECT uuid, organisation_id, username, first_name, last_name, email, phone, created_at, updated_at
        FROM users
        WHERE uuid = $1 AND deleted = false
    `
	var nullableFirstName sql.NullString
	var nullableLastName sql.NullString
	var nullablePhone sql.NullString

	user := &User{}
	err := db.QueryRow(query, uuid).Scan(&user.UUID, &user.OrganisationID, &user.Username, &nullableFirstName, &nullableLastName, &user.Email, &nullablePhone, &user.CreatedAt, &user.UpdatedAt)
	if nullableFirstName.Valid {
		user.FirstName = nullableFirstName.String
	}
	if nullableLastName.Valid {
		user.LastName = nullableLastName.String
	}
	if nullablePhone.Valid {
		user.Phone = nullablePhone.String
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by uuid: %v", err)
	}

	return user, nil
}
