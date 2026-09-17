package services

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"atk-backend/models"
	"atk-backend/repositories"
	"atk-backend/utils"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type UserService struct {
	userRepo *repositories.UserRepository
	roleRepo *repositories.RoleRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repositories.NewUserRepository(),
		roleRepo: repositories.NewRoleRepository(),
	}
}

// ============================================================
// REQUEST STRUCTS
// ============================================================

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required,min=3"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	FullName    string `json:"full_name" binding:"required"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
	RoleID      uint   `json:"role_id"`
}

type UpdateUserRequest struct {
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	Email       string `json:"email" binding:"omitempty,email"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
	RoleID      uint   `json:"role_id"`
}

type UpdateProfileRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

type ImportResult struct {
	Imported int      `json:"imported"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors"`
}

// ============================================================
// CRUD USER
// ============================================================

func (s *UserService) Create(req CreateUserRequest) (*models.User, error) {
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("username sudah digunakan")
	}

	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email sudah digunakan")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	user := &models.User{
		UUID:        uuid.New().String(),
		Username:    req.Username,
		Email:       req.Email,
		Password:    hashedPassword,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Status:      status,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	if req.RoleID > 0 {
		if _, err := s.roleRepo.FindByID(req.RoleID); err == nil {
			s.userRepo.AssignRole(user.ID, req.RoleID)
		}
	} else {
		if role, err := s.roleRepo.FindByName("staff"); err == nil {
			s.userRepo.AssignRole(user.ID, role.ID)
		}
	}

	return s.userRepo.FindByID(user.ID)
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) GetByUsername(username string) (*models.User, error) {
	return s.userRepo.FindByUsername(username)
}

func (s *UserService) GetAll(page, limit int, search string) ([]models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.userRepo.FindAll(offset, limit, search)
}

// ============================================================
// UPDATE USER (data profil — TANPA password)
// ============================================================
func (s *UserService) Update(id uint, req UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	fmt.Println("=====================================================")
	fmt.Printf("🔍 [UPDATE] ID=%d | Username LAMA='%s' | Username BARU='%s'\n",
		id, user.Username, req.Username)

	// 1. Update Username
	if req.Username != "" && req.Username != user.Username {
		if existing, err := s.userRepo.FindByUsername(req.Username); err == nil && existing.ID != id {
			return nil, errors.New("username sudah digunakan user lain")
		}
		user.Username = req.Username
	}

	// 2. Update FullName
	if req.FullName != "" {
		user.FullName = req.FullName
	}

	// 3. Update Email
	if req.Email != "" {
		if existing, err := s.userRepo.FindByEmail(req.Email); err == nil && existing.ID != id {
			return nil, errors.New("email sudah digunakan user lain")
		}
		user.Email = req.Email
	}

	// 4. Update Phone Number
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}

	// 5. Update Status
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := s.userRepo.Update(user); err != nil {
		fmt.Printf("❌ [UPDATE] Gagal simpan: %v\n", err)
		return nil, err
	}

	// 6. Update role jika dikirim
	if req.RoleID > 0 {
		s.userRepo.RemoveAllRoles(id)
		s.userRepo.AssignRole(id, req.RoleID)
	}

	updated, _ := s.userRepo.FindByID(id)
	fmt.Printf("✅ [UPDATE] Setelah re-fetch: Username='%s'\n", updated.Username)
	fmt.Println("=====================================================")

	return updated, nil
}

func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}

// ============================================================
// ROLE ASSIGNMENT
// ============================================================

func (s *UserService) AssignRole(userID, roleID uint) error {
	if _, err := s.userRepo.FindByID(userID); err != nil {
		return errors.New("user tidak ditemukan")
	}
	if _, err := s.roleRepo.FindByID(roleID); err != nil {
		return errors.New("role tidak ditemukan")
	}
	return s.userRepo.AssignRole(userID, roleID)
}

func (s *UserService) RemoveRole(userID, roleID uint) error {
	return s.userRepo.RemoveRole(userID, roleID)
}

func (s *UserService) GetUserRoles(userID uint) ([]models.Role, error) {
	return s.userRepo.GetUserRoles(userID)
}

// ============================================================
// PROFILE
// ============================================================

func (s *UserService) UpdateProfile(userID uint, req UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if existing, err := s.userRepo.FindByEmail(req.Email); err == nil && existing.ID != userID {
		return nil, errors.New("email sudah digunakan user lain")
	}

	user.FullName = req.FullName
	user.Email = req.Email

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return s.userRepo.FindByID(userID)
}

// ============================================================
// ✅ CHANGE PASSWORD — DIPERBAIKI
// ============================================================
// Menggunakan method khusus `UpdatePassword` di repository,
// BUKAN `Update` biasa (karena Update map tidak include password).
// ============================================================
func (s *UserService) ChangePassword(userID uint, req ChangePasswordRequest) error {
	fmt.Println("=====================================================")
	fmt.Printf("🔍 [CHANGE PASSWORD] UserID=%d\n", userID)

	// Validasi konfirmasi password
	if req.NewPassword != req.ConfirmPassword {
		fmt.Println("❌ [CHANGE PASSWORD] Konfirmasi tidak sama")
		return errors.New("konfirmasi password tidak sama")
	}

	// Validasi panjang minimal
	if len(req.NewPassword) < 6 {
		return errors.New("password baru minimal 6 karakter")
	}

	// Cari user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		fmt.Printf("❌ [CHANGE PASSWORD] User tidak ditemukan: %v\n", err)
		return errors.New("user tidak ditemukan")
	}

	// Verifikasi password lama
	if !utils.CheckPasswordHash(req.CurrentPassword, user.Password) {
		fmt.Println("❌ [CHANGE PASSWORD] Password lama salah")
		return errors.New("password lama salah")
	}

	// Cegah password baru sama dengan password lama
	if req.CurrentPassword == req.NewPassword {
		fmt.Println("❌ [CHANGE PASSWORD] Password baru sama dengan lama")
		return errors.New("password baru tidak boleh sama dengan password lama")
	}

	// Hash password baru
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		fmt.Printf("❌ [CHANGE PASSWORD] Gagal hash: %v\n", err)
		return err
	}

	// ✅ PAKAI METHOD KHUSUS — bukan s.userRepo.Update()
	if err := s.userRepo.UpdatePassword(userID, hashedPassword); err != nil {
		fmt.Printf("❌ [CHANGE PASSWORD] Gagal update di DB: %v\n", err)
		return err
	}

	fmt.Println("✅ [CHANGE PASSWORD] Password berhasil diubah")
	fmt.Println("=====================================================")

	return nil
}

// ============================================================
// IMPORT EXCEL
// ============================================================

func (s *UserService) ImportFromExcel(fh *multipart.FileHeader) (*ImportResult, error) {
	if fh == nil {
		return nil, errors.New("file tidak valid")
	}
	if fh.Size > 5*1024*1024 {
		return nil, errors.New("ukuran file maksimal 5MB")
	}

	src, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file: %w", err)
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("file bukan format excel valid: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, errors.New("file excel kosong")
	}
	sheet := sheets[0]

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca baris: %w", err)
	}
	if len(rows) < 2 {
		return nil, errors.New("file excel tidak memiliki data (hanya header)")
	}

	result := &ImportResult{Errors: []string{}}

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
			continue
		}
		if len(row) < 4 {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: kolom tidak lengkap", i+1))
			continue
		}

		username := strings.TrimSpace(row[0])
		password := strings.TrimSpace(row[1])
		fullName := strings.TrimSpace(row[2])
		email := strings.TrimSpace(row[3])

		if len(username) < 3 {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: username '%s' minimal 3 karakter", i+1, username))
			continue
		}
		if len(password) < 6 {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: password minimal 6 karakter", i+1))
			continue
		}
		if !strings.Contains(email, "@") {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: email '%s' tidak valid", i+1, email))
			continue
		}
		if _, err := s.userRepo.FindByUsername(username); err == nil {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: username '%s' sudah terdaftar", i+1, username))
			continue
		}
		if _, err := s.userRepo.FindByEmail(email); err == nil {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: email '%s' sudah terdaftar", i+1, email))
			continue
		}

		hashed, err := utils.HashPassword(password)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: gagal hash password", i+1))
			continue
		}

		phone := ""
		if len(row) >= 5 {
			phone = strings.TrimSpace(row[4])
		}

		var roleID uint = 0
		if len(row) >= 6 && strings.TrimSpace(row[5]) != "" {
			var rid int
			fmt.Sscanf(strings.TrimSpace(row[5]), "%d", &rid)
			roleID = uint(rid)
		}

		status := "active"
		if len(row) >= 7 {
			v := strings.ToLower(strings.TrimSpace(row[6]))
			if v != "" {
				status = v
			}
		}

		user := &models.User{
			UUID:        uuid.New().String(),
			Username:    username,
			Password:    hashed,
			FullName:    fullName,
			Email:       email,
			PhoneNumber: phone,
			Status:      status,
		}

		if err := s.userRepo.Create(user); err != nil {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: gagal simpan - %s", i+1, err.Error()))
			continue
		}

		if roleID > 0 {
			s.userRepo.AssignRole(user.ID, roleID)
		} else if role, err := s.roleRepo.FindByName("staff"); err == nil {
			s.userRepo.AssignRole(user.ID, role.ID)
		}

		result.Imported++
	}

	return result, nil
}

// ============================================================
// EXPORT EXCEL
// ============================================================

func (s *UserService) ExportToExcel(search string) (*bytes.Buffer, error) {
	users, _, err := s.userRepo.FindAll(0, 10000, search)
	if err != nil {
		return nil, fmt.Errorf("gagal query user: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Users"
	index, err := f.NewSheet(sheet)
	if err != nil {
		return nil, fmt.Errorf("gagal buat sheet: %w", err)
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"2563EB"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "1E40AF", Style: 1},
			{Type: "right", Color: "1E40AF", Style: 1},
			{Type: "top", Color: "1E40AF", Style: 1},
			{Type: "bottom", Color: "1E40AF", Style: 1},
		},
	})

	headers := []string{"ID", "Username", "Nama Lengkap", "Email", "Phone", "Status", "Dibuat"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	for i, u := range users {
		row := i + 2
		createdAt := ""
		if !u.CreatedAt.IsZero() {
			createdAt = u.CreatedAt.Format("2006-01-02 15:04:05")
		}
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), u.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), u.Username)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), u.FullName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), u.Email)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), u.PhoneNumber)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), u.Status)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), createdAt)
	}

	f.SetColWidth(sheet, "A", "A", 8)
	f.SetColWidth(sheet, "B", "B", 18)
	f.SetColWidth(sheet, "C", "C", 25)
	f.SetColWidth(sheet, "D", "D", 30)
	f.SetColWidth(sheet, "E", "E", 15)
	f.SetColWidth(sheet, "F", "F", 12)
	f.SetColWidth(sheet, "G", "G", 20)

	f.SetPanes(sheet, &excelize.Panes{
		Freeze: true, Split: false, XSplit: 0, YSplit: 1,
		TopLeftCell: "A2", ActivePane: "bottomLeft",
	})

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("gagal generate file: %w", err)
	}
	return buf, nil
}