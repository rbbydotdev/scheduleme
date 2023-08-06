package mock

import (
	"context"
	"scheduleme/models"
	"scheduleme/values"
)

var _ models.UserServiceInterface = (*UserService)(nil)

type UserService struct {
	AllFn                    func() ([]*models.User, error)
	GetOrCreateUserByEmailFn func(email string, name string) (*models.User, error)
	CreateUserFn             func(user *models.User) (values.ID, error)
	GetUserByIDFn            func(id values.ID) (*models.User, error)
	UpdateUserFn             func(user *models.User) (int64, error)
	DeleteUserFn             func(id values.ID) (int64, error)
	AttachRemoteByIDFn       func(id values.ID, routeInfo *models.RouteInfo) error
	MeForUserRouteFn         func(ri *models.RouteInfo, ctx context.Context) error
}

func (s *UserService) MeForUserRoute(ri *models.RouteInfo, ctx context.Context) error {
	return s.MeForUserRouteFn(ri, ctx)
}

func (s *UserService) AttachRemoteByID(ID values.ID, routeInfo *models.RouteInfo) error {
	return s.AttachRemoteByIDFn(ID, routeInfo)
}

func (s *UserService) All() ([]*models.User, error) {
	return s.AllFn()
}

func (s *UserService) GetOrCreateUserByEmail(email string, name string) (*models.User, error) {
	return s.GetOrCreateUserByEmailFn(email, name)
}

func (s *UserService) CreateUser(user *models.User) (values.ID, error) {
	return s.CreateUserFn(user)
}

func (s *UserService) GetUserByID(id values.ID) (*models.User, error) {
	return s.GetUserByIDFn(id)
}

func (s *UserService) UpdateUser(user *models.User) (int64, error) {
	return s.UpdateUserFn(user)
}

func (s *UserService) DeleteUser(id values.ID) (int64, error) {
	return s.DeleteUserFn(id)
}
