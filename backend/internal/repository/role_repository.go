package repository

import (
	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) FindAll() ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Preload("Permissions").Order("is_system desc, name asc").Find(&roles).Error
	if err != nil {
		return nil, err
	}

	for i := range roles {
		var userCount int64
		r.db.Model(&model.User{}).Where("role_id = ?", roles[i].ID).Count(&userCount)
		roles[i].UsersCount = int(userCount)
	}

	return roles, nil
}

func (r *RoleRepository) FindByID(id string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("id = ? OR code = ?", id, id).Preload("Permissions").First(&role).Error
	if err != nil {
		return nil, err
	}
	var userCount int64
	r.db.Model(&model.User{}).Where("role_id = ?", role.ID).Count(&userCount)
	role.UsersCount = int(userCount)
	return &role, nil
}

func (r *RoleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepository) UpdatePermissions(roleID string, permissionCodes []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Remove existing
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}

		if len(permissionCodes) == 0 {
			return nil
		}

		var perms []model.Permission
		if err := tx.Where("code IN ?", permissionCodes).Find(&perms).Error; err != nil {
			return err
		}

		for _, p := range perms {
			rp := model.RolePermission{RoleID: roleID, PermissionID: p.ID}
			if err := tx.Create(&rp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RoleRepository) FindAllPermissions() ([]model.Permission, error) {
	var perms []model.Permission
	err := r.db.Order("module asc, action asc").Find(&perms).Error
	return perms, err
}
