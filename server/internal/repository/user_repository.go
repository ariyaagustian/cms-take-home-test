package repository

import (
	"context"
	"time"

	"cms/server/internal/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserRepository interface {
	List(ctx context.Context) ([]UserWithRoles, error)
	FindWithRolesByID(ctx context.Context, id uuid.UUID) (*UserWithRoles, error)
	GetRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error)
	SetRoles(ctx context.Context, userID uuid.UUID, roleIDs []int) error
}

type UserWithRoles struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) List(ctx context.Context) ([]UserWithRoles, error) {
	var results []UserWithRoles

	rows, err := r.db.WithContext(ctx).
		Table("users").
		Select(`users.id, users.name, users.email, users.created_at, users.updated_at, 
            COALESCE(array_agg(roles.name), '{}') AS roles`).
		Joins("LEFT JOIN user_roles ON users.id = user_roles.user_id").
		Joins("LEFT JOIN roles ON roles.id = user_roles.role_id").
		Group("users.id").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u UserWithRoles
		var roles pq.StringArray
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt, &roles); err != nil {
			return nil, err
		}
		u.Roles = roles
		results = append(results, u)
	}

	return results, nil
}

func (r *userRepository) FindWithRolesByID(ctx context.Context, id uuid.UUID) (*UserWithRoles, error) {
	row := r.db.WithContext(ctx).
		Table("users").
		Select(`users.id, users.name, users.email, users.created_at, users.updated_at,
                COALESCE(array_agg(roles.name), '{}') AS roles`).
		Joins("LEFT JOIN user_roles ur ON ur.user_id = users.id").
		Joins("LEFT JOIN roles ON roles.id = ur.role_id").
		Where("users.id = ?", id).
		Group("users.id").
		Row()

	var u UserWithRoles
	var roles pq.StringArray // gunakan pq untuk scan array
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt, &roles); err != nil {
		return nil, err
	}
	u.Roles = roles
	return &u, nil
}

func (r *userRepository) GetRoles(ctx context.Context, userID uuid.UUID) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.WithContext(ctx).
		Table("roles").
		Joins("JOIN user_roles ur ON ur.role_id = roles.id").
		Where("ur.user_id = ?", userID).
		Scan(&roles).Error
	return roles, err
}

func (r *userRepository) SetRoles(ctx context.Context, userID uuid.UUID, roleIDs []int) error {
	// Hapus role sebelumnya
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
		return err
	}

	// Bulk insert roles
	var userRoles []model.UserRole
	for _, rid := range roleIDs {
		userRoles = append(userRoles, model.UserRole{UserID: userID, RoleID: rid})
	}
	return r.db.WithContext(ctx).Create(&userRoles).Error
}
