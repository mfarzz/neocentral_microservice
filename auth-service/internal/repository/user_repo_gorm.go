package repository

import (
	"context"
	"errors"
	"strings"

	"neocentral-go/auth-service/internal/domain"

	"gorm.io/gorm"
)

type gormUserRepo struct {
	db *gorm.DB
}

// NewGormUserRepository returns a GORM-backed UserRepository.
func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepo{db: db}
}

func (r *gormUserRepo) preloadFull(tx *gorm.DB) *gorm.DB {
	return tx.
		Preload("UserHasRoles", "status = ?", domain.RoleStatusActive).
		Preload("UserHasRoles.Role").
		Preload("Student").
		Preload("Lecturer")
}

func (r *gormUserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.preloadFull(r.db.WithContext(ctx)).First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *gormUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.preloadFull(r.db.WithContext(ctx)).First(&user, "LOWER(email) = ?", strings.ToLower(email)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *gormUserRepo) FindByIdentityNumber(ctx context.Context, identityNumber string) (*domain.User, error) {
	var user domain.User
	err := r.preloadFull(r.db.WithContext(ctx)).First(&user, "identity_number = ?", identityNumber).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *gormUserRepo) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *gormUserRepo) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *gormUserRepo) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(fields).Error
}

func (r *gormUserRepo) GetUserRoles(ctx context.Context, userID string) ([]domain.UserHasRole, error) {
	var roles []domain.UserHasRole
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("user_id = ? AND status = ?", userID, domain.RoleStatusActive).
		Find(&roles).Error
	return roles, err
}

func (r *gormUserRepo) HasRole(ctx context.Context, userID, roleName string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.UserHasRole{}).
		Joins("JOIN user_roles ON user_roles.id = user_has_roles.role_id").
		Where("user_has_roles.user_id = ? AND LOWER(user_roles.name) = LOWER(?) AND user_has_roles.status = ?",
			userID, roleName, domain.RoleStatusActive).
		Count(&count).Error
	return count > 0, err
}

func (r *gormUserRepo) GetUsersByRole(ctx context.Context, roleName string, page, pageSize int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	subquery := r.db.WithContext(ctx).
		Model(&domain.UserHasRole{}).
		Select("user_id").
		Joins("JOIN user_roles ON user_roles.id = user_has_roles.role_id").
		Where("LOWER(user_roles.name) = LOWER(?) AND user_has_roles.status = ?",
			roleName, domain.RoleStatusActive)

	base := r.db.WithContext(ctx).Where("id IN (?)", subquery)

	base.Model(&domain.User{}).Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	err := r.preloadFull(base).
		Offset(offset).Limit(pageSize).
		Find(&users).Error

	return users, total, err
}

func (r *gormUserRepo) BatchGetByIDs(ctx context.Context, ids []string) ([]domain.User, error) {
	var users []domain.User
	err := r.preloadFull(r.db.WithContext(ctx)).
		Where("id IN ?", ids).
		Find(&users).Error
	return users, err
}
