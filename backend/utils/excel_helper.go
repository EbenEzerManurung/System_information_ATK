package utils

import (
    "bytes"
    "fmt"
    "strconv"
    "time"
    "github.com/xuri/excelize/v2"
)

type ExcelColumn struct {
    Header string
    Width  float64
    Field  string
}

func CreateExcelFile(sheetName string, columns []ExcelColumn, data []map[string]interface{}) (*bytes.Buffer, error) {
    f := excelize.NewFile()
    defer f.Close()
    
    // Create sheet
    index, err := f.NewSheet(sheetName)
    if err != nil {
        return nil, err
    }
    f.SetActiveSheet(index)
    f.DeleteSheet("Sheet1")
    
    // Set headers
    for col, column := range columns {
        cell := fmt.Sprintf("%s1", string(rune('A'+col)))
        f.SetCellValue(sheetName, cell, column.Header)
        if column.Width > 0 {
            f.SetColWidth(sheetName, string(rune('A'+col)), string(rune('A'+col)), column.Width)
        }
    }
    
    // Set data
    for row, item := range data {
        for col, column := range columns {
            cell := fmt.Sprintf("%s%d", string(rune('A'+col)), row+2)
            if val, ok := item[column.Field]; ok {
                f.SetCellValue(sheetName, cell, val)
            }
        }
    }
    
    // Write to buffer
    buf := new(bytes.Buffer)
    if err := f.Write(buf); err != nil {
        return nil, err
    }
    
    return buf, nil
}

func ParseExcelFile(buf *bytes.Buffer, columns []ExcelColumn) ([]map[string]interface{}, error) {
    f, err := excelize.OpenReader(buf)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    
    sheetName := f.GetSheetName(0)
    rows, err := f.GetRows(sheetName)
    if err != nil {
        return nil, err
    }
    
    if len(rows) < 2 {
        return []map[string]interface{}{}, nil
    }
    
    // Map headers to columns
    headerMap := make(map[string]int)
    for i, col := range columns {
        headerMap[col.Header] = i
    }
    
    result := make([]map[string]interface{}, 0)
    
    // Parse data rows (skip header)
    for _, row := range rows[1:] {
        item := make(map[string]interface{})
        for _, col := range columns {
            if idx, ok := headerMap[col.Header]; ok && idx < len(row) {
                val := row[idx]
                item[col.Field] = val
            }
        }
        result = append(result, item)
    }
    
    return result, nil
}

func ParseExcelFileWithValidation(buf *bytes.Buffer, columns []ExcelColumn, validators map[string]func(interface{}) error) ([]map[string]interface{}, []map[string]string, error) {
    data, err := ParseExcelFile(buf, columns)
    if err != nil {
        return nil, nil, err
    }
    
    validData := make([]map[string]interface{}, 0)
    errors := make([]map[string]string, 0)
    
    for idx, item := range data {
        rowErrors := make(map[string]string)
        valid := true
        
        for field, validator := range validators {
            if err := validator(item[field]); err != nil {
                rowErrors[field] = err.Error()
                valid = false
            }
        }
        
        if valid {
            validData = append(validData, item)
        } else {
            rowErrors["row"] = strconv.Itoa(idx + 2)
            errors = append(errors, rowErrors)
        }
    }
    
    return validData, errors, nil
}

// Common validators
func ValidateRequired(val interface{}) error {
    if val == nil || val == "" {
        return fmt.Errorf("field is required")
    }
    return nil
}

func ValidateNumeric(val interface{}) error {
    if val == nil || val == "" {
        return nil
    }
    str := fmt.Sprintf("%v", val)
    _, err := strconv.ParseFloat(str, 64)
    if err != nil {
        return fmt.Errorf("must be a number")
    }
    return nil
}

func ValidateInt(val interface{}) error {
    if val == nil || val == "" {
        return nil
    }
    str := fmt.Sprintf("%v", val)
    _, err := strconv.Atoi(str)
    if err != nil {
        return fmt.Errorf("must be an integer")
    }
    return nil
}

func ValidateDate(val interface{}) error {
    if val == nil || val == "" {
        return nil
    }
    str := fmt.Sprintf("%v", val)
    _, err := time.Parse("2006-01-02", str)
    if err != nil {
        return fmt.Errorf("invalid date format, use YYYY-MM-DD")
    }
    return nil
}