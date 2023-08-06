package models

import (
	"context"
	"database/sql"
	"scheduleme/frame"
	sq "scheduleme/sqlite"
	"scheduleme/toerr"
	"scheduleme/values"
)

func NewUserService(db *sq.Db) UserServiceInterface {
	return &UserService{db: db}
}

// User
func UserIDFromContext(ctx context.Context) values.ID {
	return frame.FromContext[SessionInfo](ctx).UserID
}

type User struct {
	IsAdmin bool      `json:"is_admin,omitempty"`
	ID      values.ID `json:"id,omitempty"`
	Email   string    `json:"email,omitempty"`
	Name    string    `json:"name,omitempty"`
}

type UserMutate struct {
	IsAdmin bool   `json:"is_admin,omitempty"`
	Name    string `json:"name,omitempty"`
}

type UserView struct {
	ID   values.ID `json:"id"`
	Name string    `json:"name"`
}
type UserViewPrivate struct {
	ID    values.ID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func (u User) View() UserView {
	return UserView{
		ID:   u.ID,
		Name: u.Name,
	}
}

func (u User) ViewPrivate() UserViewPrivate {
	return UserViewPrivate{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func (um *UserMutate) Validate() error {
	return nil
}

func (um *UserMutate) ModifiesBodyInfo(bi *BodyInfo, ri RouteInfo, si SessionInfo) {
	bi.User = &User{
		ID:   ri.User.ID,
		Name: um.Name,
	}
	if si.IsAdmin {
		bi.User.IsAdmin = um.IsAdmin
	}

}

func (u *User) New(email string, name string) *User {
	return &User{Email: email, Name: name}
}

type UserService struct {
	db *sq.Db
}

// Ensure UserService implements UserServiceInterface
var _ UserServiceInterface = (*UserService)(nil)

type UserServiceInterface interface {
	All() (users []*User, err error)
	GetOrCreateUserByEmail(email string, name string) (*User, error)
	CreateUser(user *User) (values.ID, error)
	GetUserByID(id values.ID) (*User, error)
	UpdateUser(user *User) (int64, error)
	DeleteUser(id values.ID) (int64, error)
	AttachRemoteByID(ID values.ID, ri *RouteInfo) error
	MeForUserRoute(ri *RouteInfo, ctx context.Context) error
}

// A method to implicitly attach a self (user) to a route info object when {UserID} is not in the route but Session User data is
func (us *UserService) MeForUserRoute(ri *RouteInfo, ctx context.Context) (err error) {
	userID := UserIDFromContext(ctx)
	u, err := us.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return toerr.NotFound(err)
		}
		return toerr.Internal(err)
	}
	ri.User = *u
	return
}

func (us *UserService) AttachRemoteByID(ID values.ID, ri *RouteInfo) (err error) {
	u, err := us.GetUserByID(ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return toerr.NotFound(err)
		}
		return toerr.Internal(err)
	}

	ri.User = *u
	return
}

// All returns all users from the database
func (us *UserService) All() ([]*User, error) {
	var users []*User
	rows, err := us.db.Query(`SELECT id, email, name FROM users`)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, toerr.NotFound(err).Msg("no users")
		}
		return nil, toerr.Internal(err)
	}
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ID, &user.Email, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// GetOrCreateUserByEmail retrieves a user by email, or creates a new one if not found.
func (us *UserService) GetOrCreateUserByEmail(email string, name string) (*User, error) {
	user := &User{}
	err := us.db.QueryRow(`SELECT id, email, name FROM users WHERE email = ?`, email).Scan(
		&user.ID, &user.Email, &user.Name)

	// If user does not exist, create a new one.
	var id values.ID
	if err == sql.ErrNoRows {
		user = &User{Email: email, Name: name}
		id, err = us.CreateUser(user)
		user.ID = id
	}

	if err != nil {
		return nil, toerr.Internal(err).Msg("get or create user failed")
	}

	return user, nil
}

// CreateUser inserts a new user into the database
func (us *UserService) CreateUser(user *User) (values.ID, error) {
	res, err := us.db.Exec(`INSERT INTO users (email, name) VALUES (?, ?)`,
		user.Email, user.Name)
	if err != nil {
		return 0, toerr.Internal(err).Msg("create user failed")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return values.ID(id), nil
}

// GetUserById finds a user by ID
func (us *UserService) GetUserByID(id values.ID) (*User, error) {
	user := &User{}
	err := us.db.QueryRow(`SELECT id, email, name FROM users WHERE id = ?`, id).Scan(
		&user.ID, &user.Email, &user.Name)
	if err == sql.ErrNoRows {
		return nil, toerr.NotFound(err)
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates an existing user in the database
func (us *UserService) UpdateUser(user *User) (int64, error) {
	return withCount(
		us.db.Exec(`UPDATE users SET email = ?, name = ? WHERE id = ?`,
			user.Email, user.Name, user.ID),
	)
}

// DeleteUser deletes a user from the database
func (us *UserService) DeleteUser(id values.ID) (int64, error) {
	return withCount(
		us.db.Exec(`DELETE FROM users WHERE id = ?`, id),
	)
}
