package repositories

import (
	"atk-backend/database"
	"atk-backend/models"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{
		db: database.GetDB(),
	}
}

// ============================================================
// CREATE
// ============================================================
func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

// ============================================================
// FIND BY ID
// ============================================================
func (r *RoleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// ============================================================
// FIND BY NAME
// ============================================================
func (r *RoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// ============================================================
// FIND ALL — basic (backward compat)
// ============================================================
func (r *RoleRepository) FindAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.
		Order("level ASC, id ASC").
		Find(&roles).Error
	return roles, err
}

// ============================================================
// FIND ALL WITH SEARCH — untuk export & pagination
// ✅ METHOD BARU
// ============================================================
func (r *RoleRepository) FindAllWithSearch(
	offset, limit int,
	search string,
) ([]models.Role, error) {
	var roles []models.Role

	query := r.db.Model(&models.Role{})

	if search != "" {
		like := "%" + search + "%"
		query = query.Where(
			"name LIKE ? OR display_name LIKE ? OR description LIKE ?",
			like, like, like,
		)
	}

	err := query.
		Order("level ASC, id ASC").
		Offset(offset).
		Limit(limit).
		Find(&roles).Error

	return roles, err
}

// ============================================================
// UPDATE
// ============================================================
func (r *RoleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

// ============================================================
// DELETE
// ============================================================
func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}