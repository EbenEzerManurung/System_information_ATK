package repositories

import (
	"fmt"

	"atk-backend/database"
	"atk-backend/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.GetDB(),
	}
}

// ============================================================
// CRUD
// ============================================================

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").First(&user, id).Error
	return &user, err
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindAll(offset, limit int, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// ✅ Unscoped() agar user yang ter-soft-delete tetap terhitung
	// (opsional — jika Anda mau user yang sudah hard delete tidak muncul, biarkan seperti ini)
	query := r.db.Unscoped().Model(&models.User{}).Preload("Roles")

	if search != "" {
		like := "%" + search + "%"
		query = query.Where(
			"username LIKE ? OR email LIKE ? OR full_name LIKE ?",
			like, like, like,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error

	return users, total, err
}

// ============================================================
// UPDATE — untuk data profil user (TANPA password)
// ============================================================
func (r *UserRepository) Update(user *models.User) error {
	fmt.Printf("🔍 [REPO] Update ID=%d | Username='%s'\n", user.ID, user.Username)

	result := r.db.Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"username":     user.Username,
			"full_name":    user.FullName,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"status":       user.Status,
		})

	fmt.Printf("✅ [REPO] Update RowsAffected=%d | Error=%v\n",
		result.RowsAffected, result.Error)

	return result.Error
}

// ============================================================
// ✅ UPDATE PASSWORD — method KHUSUS untuk ganti password
// ============================================================
func (r *UserRepository) UpdatePassword(userID uint, hashedPassword string) error {
	fmt.Printf("🔍 [REPO] UpdatePassword ID=%d\n", userID)

	result := r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("password", hashedPassword)

	fmt.Printf("✅ [REPO] UpdatePassword RowsAffected=%d | Error=%v\n",
		result.RowsAffected, result.Error)

	return result.Error
}

// ============================================================
// ✅ DELETE — HARD DELETE dengan Unscoped()
// ============================================================
// Menggunakan Unscoped() untuk MEMAKSA hard delete, mengabaikan
// apapun terkait soft delete. Ini adalah jaring pengaman agar
// data dijamin terhapus permanen dari database.
//
// Urutan:
//   1. Hapus relasi user_roles (pivot table)
//   2. Hapus user dengan Unscoped().Delete() → DELETE FROM users
// ============================================================
func (r *UserRepository) Delete(id uint) error {
	fmt.Printf("🗑️ [REPO] HARD DELETE user ID=%d\n", id)

	// 1. Hapus relasi user_roles dulu
	//    Pakai Unscoped() juga untuk jaga-jaga
	if err := r.db.Unscoped().
		Where("user_id = ?", id).
		Delete(&models.UserRole{}).Error; err != nil {
		fmt.Printf("⚠️ [REPO] Gagal hapus user_roles: %v\n", err)
		// Lanjut saja, tidak fatal — user tetap harus dihapus
	}

	// 2. ✅ FORCE HARD DELETE dengan Unscoped()
	//    SQL yang dijalankan: DELETE FROM users WHERE id = ?
	result := r.db.Unscoped().Delete(&models.User{}, id)

	fmt.Printf("✅ [REPO] Delete RowsAffected=%d | Error=%v\n",
		result.RowsAffected, result.Error)

	return result.Error
}

// ============================================================
// ROLE ASSIGNMENT
// ============================================================

// AssignRole - tambahkan role ke user (skip jika sudah ada)
func (r *UserRepository) AssignRole(userID, roleID uint) error {
	var count int64
	if err := r.db.Model(&models.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	userRole := models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.Create(&userRole).Error
}

// RemoveRole - hapus satu role spesifik dari user
func (r *UserRepository) RemoveRole(userID, roleID uint) error {
	return r.db.
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&models.UserRole{}).Error
}

// RemoveAllRoles - hapus SEMUA role user (dipakai saat replace role)
func (r *UserRepository) RemoveAllRoles(userID uint) error {
	return r.db.
		Where("user_id = ?", userID).
		Delete(&models.UserRole{}).Error
}

// GetUserRoles - ambil semua role milik user
func (r *UserRepository) GetUserRoles(userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}