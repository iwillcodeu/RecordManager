package test

import (
	"testing"
	"github.com/xuri/excelize"
	"fmt"
	"os"
	"../../"
)


func Test(t *testing.T){
	expected := "OK!"
	actual := hello("excel2DB")
	if actual != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

func hello(input string) string {
	switch input {
	case "fileXlsx":
		fileXlsx("/Users/iwillcodeu/Downloads/表1 红线问题分类编码表.xlsx")

	case "excel2DB":
		saveExcel2DB()
	}
	return "OK"
}


func fileXlsx(filePath string) {
	xlsx, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rows := xlsx.GetRows("Sheet1")
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
	}
}