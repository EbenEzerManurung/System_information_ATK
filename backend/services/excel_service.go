package services

import (
    "bytes"
    "context"
    "fmt"
    "strconv"
    "sync"
    "time"
    
    "github.com/google/uuid"
    "github.com/xuri/excelize/v2"
    
    "atk-backend/models"
    "atk-backend/repositories"
)

type ExcelService struct {
    masterRepo *repositories.MasterRepository
    stockRepo  *repositories.StockRepository
}

func NewExcelService() *ExcelService {
    return &ExcelService{
        masterRepo: repositories.NewMasterRepository(),
        stockRepo:  repositories.NewStockRepository(),
    }
}

// ExportMasterData exports up to 1 million records with streaming
func (s *ExcelService) ExportMasterData(ctx context.Context) (*bytes.Buffer, int64, error) {
    f := excelize.NewFile()
    defer f.Close()
    
    sheetName := "Master Data"
    index, err := f.NewSheet(sheetName)
    if err != nil {
        return nil, 0, err
    }
    f.SetActiveSheet(index)
    f.DeleteSheet("Sheet1")
    
    // Set headers
    headers := []string{"Kode", "Nama", "Kategori", "Deskripsi", "Unit", "Min Stock", "Max Stock", "Harga", "Stok"}
    for i, header := range headers {
        cell := fmt.Sprintf("%s1", string(rune('A'+i)))
        f.SetCellValue(sheetName, cell, header)
    }
    
    // Set column widths
    f.SetColWidth(sheetName, "A", "A", 15)
    f.SetColWidth(sheetName, "B", "B", 30)
    f.SetColWidth(sheetName, "C", "C", 20)
    f.SetColWidth(sheetName, "D", "D", 40)
    f.SetColWidth(sheetName, "E", "E", 10)
    f.SetColWidth(sheetName, "F", "F", 15)
    f.SetColWidth(sheetName, "G", "G", 15)
    f.SetColWidth(sheetName, "H", "H", 15)
    f.SetColWidth(sheetName, "I", "I", 15)
    
    // Get total count
    total, err := s.masterRepo.CountAll()
    if err != nil {
        return nil, 0, err
    }
    
    // Stream data in chunks
    chunkSize := 1000
    offset := 0
    row := 2
    
    for {
        select {
        case <-ctx.Done():
            return nil, 0, ctx.Err()
        default:
        }
        
        masters, _, err := s.masterRepo.FindAll(offset, chunkSize, "", "")
        if err != nil {
            break
        }
        
        if len(masters) == 0 {
            break
        }
        
        for _, master := range masters {
            stock, _ := s.stockRepo.FindByMasterID(master.ID)
            stockQty := 0
            if stock != nil {
                stockQty = stock.Quantity
            }
            
            f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), master.Code)
            f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), master.Name)
            f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), master.Category)
            f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), master.Description)
            f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), master.Unit)
            f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), master.MinStock)
            f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), master.MaxStock)
            f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), master.Price)
            f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), stockQty)
            
            row++
        }
        
        offset += chunkSize
        if offset >= int(total) {
            break
        }
    }
    
    buf := new(bytes.Buffer)
    if err := f.Write(buf); err != nil {
        return nil, 0, err
    }
    
    return buf, total, nil
}

// ImportMasterData imports data with concurrent processing
func (s *ExcelService) ImportMasterData(ctx context.Context, buf *bytes.Buffer) (int, []map[string]string, error) {
    startTime := time.Now()
    
    f, err := excelize.OpenReader(buf)
    if err != nil {
        return 0, nil, fmt.Errorf("failed to read Excel file: %v", err)
    }
    defer f.Close()
    
    sheetName := f.GetSheetName(0)
    rows, err := f.GetRows(sheetName)
    if err != nil {
        return 0, nil, fmt.Errorf("failed to get rows: %v", err)
    }
    
    if len(rows) < 2 {
        return 0, nil, fmt.Errorf("no data found")
    }
    
    totalRows := len(rows) - 1
    if totalRows > 1000000 {
        return 0, nil, fmt.Errorf("maximum 1 million records allowed (got %d)", totalRows)
    }
    
    // HAPUS: validData tidak digunakan
    // validData := make([]models.Master, 0)
    
    // Gunakan nama variabel berbeda untuk menghindari konflik dengan package errors
    errorList := make([]map[string]string, 0)
    
    var wg sync.WaitGroup
    var mu sync.Mutex
    semaphore := make(chan struct{}, 10)
    
    type ProcessResult struct {
        Masters []models.Master
        Errors  []map[string]string
    }
    
    resultChan := make(chan ProcessResult, 10)
    
    chunkSize := 100
    for i := 1; i < len(rows); i += chunkSize {
        select {
        case <-ctx.Done():
            return 0, nil, ctx.Err()
        default:
        }
        
        end := i + chunkSize
        if end > len(rows) {
            end = len(rows)
        }
        
        wg.Add(1)
        go func(start, end int) {
            defer wg.Done()
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            result := ProcessResult{
                Masters: make([]models.Master, 0),
                Errors:  make([]map[string]string, 0),
            }
            
            for rowIdx := start; rowIdx < end; rowIdx++ {
                row := rows[rowIdx]
                if len(row) < 8 {
                    continue
                }
                
                code := row[0]
                name := row[1]
                category := row[2]
                description := row[3]
                unit := row[4]
                minStock, _ := strconv.Atoi(row[5])
                maxStock, _ := strconv.Atoi(row[6])
                price, _ := strconv.ParseFloat(row[7], 64)
                
                rowErrors := make(map[string]string)
                valid := true
                
                if code == "" {
                    rowErrors["Kode"] = "kode tidak boleh kosong"
                    valid = false
                }
                if name == "" {
                    rowErrors["Nama"] = "nama tidak boleh kosong"
                    valid = false
                }
                if category == "" {
                    rowErrors["Kategori"] = "kategori tidak boleh kosong"
                    valid = false
                }
                if unit == "" {
                    rowErrors["Unit"] = "unit tidak boleh kosong"
                    valid = false
                }
                if price <= 0 {
                    rowErrors["Harga"] = "harga harus lebih dari 0"
                    valid = false
                }
                
                if valid {
                    master := models.Master{
                        UUID:        uuid.New().String(),
                        Code:        code,
                        Name:        name,
                        Category:    category,
                        Description: description,
                        Unit:        unit,
                        MinStock:    minStock,
                        MaxStock:    maxStock,
                        Price:       price,
                    }
                    result.Masters = append(result.Masters, master)
                } else {
                    rowErrors["Row"] = strconv.Itoa(rowIdx + 1)
                    result.Errors = append(result.Errors, rowErrors)
                }
            }
            
            resultChan <- result
        }(i, end)
    }
    
    go func() {
        wg.Wait()
        close(resultChan)
    }()
    
    mastersToInsert := make([]models.Master, 0)
    
    for result := range resultChan {
        mu.Lock()
        mastersToInsert = append(mastersToInsert, result.Masters...)
        errorList = append(errorList, result.Errors...)
        mu.Unlock()
    }
    
    inserted := 0
    if len(mastersToInsert) > 0 {
        chunkSizeDB := 1000
        for i := 0; i < len(mastersToInsert); i += chunkSizeDB {
            select {
            case <-ctx.Done():
                return inserted, errorList, ctx.Err()
            default:
            }
            
            end := i + chunkSizeDB
            if end > len(mastersToInsert) {
                end = len(mastersToInsert)
            }
            
            chunk := mastersToInsert[i:end]
            err := s.masterRepo.BulkInsert(chunk)
            if err != nil {
                return inserted, errorList, fmt.Errorf("failed to insert chunk: %v", err)
            }
            inserted += len(chunk)
        }
    }
    
    elapsed := time.Since(startTime)
    fmt.Printf("✅ Imported %d records in %v\n", inserted, elapsed)
    
    return inserted, errorList, nil
}

// ExportTemplate creates a template for import
func (s *ExcelService) ExportTemplate() (*bytes.Buffer, error) {
    f := excelize.NewFile()
    defer f.Close()
    
    sheetName := "Template Import Master"
    index, err := f.NewSheet(sheetName)
    if err != nil {
        return nil, err
    }
    f.SetActiveSheet(index)
    f.DeleteSheet("Sheet1")
    
    headers := []string{
        "Kode*", "Nama*", "Kategori*", "Deskripsi",
        "Unit*", "Min Stock", "Max Stock", "Harga*", "Stok Awal",
    }
    
    comments := []string{
        "Kode unik untuk item (wajib)",
        "Nama item (wajib)",
        "Kategori: alat tulis/kertas/elektronik/perlengkapan (wajib)",
        "Deskripsi item (opsional)",
        "Unit: pcs/rim/pack/set (wajib)",
        "Minimal stock (opsional, default 0)",
        "Maksimal stock (opsional, default 0)",
        "Harga per unit dalam Rupiah (wajib)",
        "Stok awal (opsional, default 0)",
    }
    
    for i, header := range headers {
        cell := fmt.Sprintf("%s1", string(rune('A'+i)))
        f.SetCellValue(sheetName, cell, header)
        f.SetCellHyperLink(sheetName, cell, comments[i], "External")
    }
    
    f.SetColWidth(sheetName, "A", "A", 15)
    f.SetColWidth(sheetName, "B", "B", 30)
    f.SetColWidth(sheetName, "C", "C", 20)
    f.SetColWidth(sheetName, "D", "D", 40)
    f.SetColWidth(sheetName, "E", "E", 10)
    f.SetColWidth(sheetName, "F", "F", 15)
    f.SetColWidth(sheetName, "G", "G", 15)
    f.SetColWidth(sheetName, "H", "H", 15)
    f.SetColWidth(sheetName, "I", "I", 15)
    
    examples := [][]interface{}{
        {"AT011", "Pulpen Gel", "alat tulis", "Pulpen gel dengan tinta biru", "pcs", 50, 500, 3000, 200},
        {"AT012", "Binder Clip", "perlengkapan", "Binder clip ukuran 25mm", "pcs", 30, 300, 5000, 100},
    }
    
    for rowIdx, example := range examples {
        for colIdx, val := range example {
            cell := fmt.Sprintf("%s%d", string(rune('A'+colIdx)), rowIdx+2)
            f.SetCellValue(sheetName, cell, val)
        }
    }
    
    buf := new(bytes.Buffer)
    if err := f.Write(buf); err != nil {
        return nil, err
    }
    
    return buf, nil
}