package data

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var db *sql.DB

const dbTimeout = time.Second * 3

// Models is the type for this package. Note that any data that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the data is also added in the New function.
type Models struct {
	User User
}

// User is the structure which holds one user from the database.
type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {

	db = dbPool
	return Models{User: User{}}
}

// GetByEmail returns one user by email
func (u User) GetByEmail(email string) (*User, error) {

	var user User
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where email = $1`

	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *User) PasswordMatches(plainText string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
