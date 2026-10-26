package service

import (
	"Museum/internal/repository"
	"context"
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type ReportService struct {
	repo *repository.ExhibitRepository
}

func NewReportService(repo *repository.ExhibitRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) Report(ctx context.Context) (*excelize.File, error) {
	list, total, err := s.repo.List(ctx, 0, 100)
	if err != nil {
		return nil, fmt.Errorf("failed create report: %w", err)
	}
	res := make(map[string]int, 10)
	for _, v := range list {
		res[v.Category]++
	}
	f := excelize.NewFile()
	defer f.Close()
	//sw, _ := f.NewStreamWriter("Sheet1")
	//sw.SetRow("A1", []interface{}{"Вид предметов", "Число предметов"})
	//sw.SetRow("A2", []interface{}{"Всего", total})
	//i := 3
	//for k, v := range res {
	//	sw.SetRow("A"+strconv.Itoa(i), []interface{}{k, v})
	//	i++
	//}
	//sw.Flush()
	sheetName := "Sheet1"
	f.SetCellValue(sheetName, "A1", "Вид предметов")
	f.SetCellValue(sheetName, "B1", "Число предметов")
	f.SetCellValue(sheetName, "A2", "Всего")
	f.SetCellValue(sheetName, "B2", total)
	f.SetCellValue(sheetName, "B2", total)
	i := 3
	for k, v := range res {
		f.SetCellValue(sheetName, "A"+strconv.Itoa(i), k)
		f.SetCellValue(sheetName, "b"+strconv.Itoa(i), v)
		i++
	}
	return f, nil
}
