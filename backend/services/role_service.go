package services

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"atk-backend/models"
	"atk-backend/repositories"

	"github.com/xuri/excelize/v2"
)

// ============================================================
// ROLE SERVICE
// ============================================================
type RoleService struct {
	roleRepo *repositories.RoleRepository
}

func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo: repositories.NewRoleRepository(),
	}
}

// ============================================================
// REQUEST STRUCTS
// ============================================================
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

type UpdateRoleRequest struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

// ImportRoleResult — hasil import excel
type ImportRoleResult struct {
	Imported int      `json:"imported"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors"`
}

// ============================================================
// CREATE
// ============================================================
func (s *RoleService) Create(req CreateRoleRequest) (*models.Role, error) {
	// Cek duplikat
	_, err := s.roleRepo.FindByName(req.Name)
	if err == nil {
		return nil, errors.New("role name already exists")
	}

	role := &models.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Level:       req.Level,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}
	return s.roleRepo.FindByID(role.ID)
}

// ============================================================
// GET BY ID
// ============================================================
func (s *RoleService) GetByID(id uint) (*models.Role, error) {
	return s.roleRepo.FindByID(id)
}

// ============================================================
// GET ALL
// ============================================================
func (s *RoleService) GetAll() ([]models.Role, error) {
	return s.roleRepo.FindAll()
}

// ============================================================
// UPDATE
// ============================================================
func (s *RoleService) Update(id uint, req UpdateRoleRequest) (*models.Role, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("role not found")
	}

	if req.DisplayName != "" {
		role.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Level > 0 {
		role.Level = req.Level
	}

	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}
	return s.roleRepo.FindByID(id)
}

// ============================================================
// DELETE
// ============================================================
func (s *RoleService) Delete(id uint) error {
	return s.roleRepo.Delete(id)
}

// ============================================================
// IMPORT EXCEL
// Format kolom (baris 1 = header):
// A: Name | B: DisplayName | C: Description | D: Level
// ============================================================
func (s *RoleService) ImportFromExcel(fh *multipart.FileHeader) (*ImportRoleResult, error) {
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
		return nil, fmt.Errorf("file bukan excel valid: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, errors.New("file excel kosong")
	}
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("gagal baca baris: %w", err)
	}
	if len(rows) < 2 {
		return nil, errors.New("file excel tidak memiliki data")
	}

	result := &ImportRoleResult{Errors: []string{}}

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
			continue
		}
		if len(row) < 1 {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: kolom tidak lengkap", i+1))
			continue
		}

		name := strings.TrimSpace(row[0])
		displayName := ""
		if len(row) >= 2 {
			displayName = strings.TrimSpace(row[1])
		}
		description := ""
		if len(row) >= 3 {
			description = strings.TrimSpace(row[2])
		}
		level := 0
		if len(row) >= 4 && strings.TrimSpace(row[3]) != "" {
			level, _ = strconv.Atoi(strings.TrimSpace(row[3]))
		}

		if name == "" {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: nama role kosong", i+1))
			continue
		}

		// Cek duplikat
		existing, _ := s.roleRepo.FindByName(name)
		if existing != nil && existing.ID > 0 {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: role '%s' sudah ada", i+1, name))
			continue
		}

		role := &models.Role{
			Name:        name,
			DisplayName: displayName,
			Description: description,
			Level:       level,
		}

		if err := s.roleRepo.Create(role); err != nil {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("Baris %d: gagal simpan - %s", i+1, err.Error()))
			continue
		}

		result.Imported++
	}

	return result, nil
}

// ============================================================
// EXPORT EXCEL
// ============================================================
func (s *RoleService) ExportToExcel(search string) (*bytes.Buffer, error) {
	roles, err := s.roleRepo.FindAllWithSearch(0, 10000, search)
	if err != nil {
		return nil, fmt.Errorf("gagal query role: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Roles"
	idx, _ := f.NewSheet(sheet)
	f.SetActiveSheet(idx)
	_ = f.DeleteSheet("Sheet1")

	// Header style
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

	headers := []string{"No", "Name", "Display Name", "Deskripsi", "Level"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	for i, r := range roles {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), r.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), r.DisplayName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), r.Description)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), r.Level)
	}

	f.SetColWidth(sheet, "A", "A", 6)
	f.SetColWidth(sheet, "B", "B", 20)
	f.SetColWidth(sheet, "C", "C", 25)
	f.SetColWidth(sheet, "D", "D", 40)
	f.SetColWidth(sheet, "E", "E", 10)

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