package users

import (
	"context"

	"github.com/hxl1513/ffgoapi/dto"
	"github.com/hxl1513/ffgoapi/model"
	"github.com/hxl1513/ffgoapi/service"
	"gorm.io/gorm"
)

// UserRepository encapsulates database operations for users.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user record.
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID retrieves a user by primary key.  Returns nil, nil when the record
// does not exist (i.e. gorm.ErrRecordNotFound).
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UserService orchestrates user-related business logic.
type UserService struct {
	service.BaseService
	// Repo is exported so that callers in other packages (e.g. tests)
	// can perform read-side operations directly.
	Repo *UserRepository
}

// NewUserService wires up a UserService with a database handle.
func NewUserService(db *model.DB) *UserService {
	svc := service.NewBaseService(db)
	return &UserService{
		BaseService: *svc,
		Repo:        NewUserRepository(db.Write),
	}
}

// hashPassword is a simplified (non-production) hashing helper.
// It prefixes the plaintext with "hashed_" so that tests can verify the
// transformation without pulling in bcrypt/argon2 dependencies.
func hashPassword(password string) string {
	return "hashed_" + password
}

// CreateUser creates a new user from a CreateUserRequest.
func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashPassword(req.Password),
		Role:     2, // default role: regular user
	}
	if err := s.Repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return dto.FromUser(user), nil
}
