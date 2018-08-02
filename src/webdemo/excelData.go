package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
	"fmt"
	"reflect"
	"log"
	"strconv"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jinzhu/gorm"
)

type ExcelData interface {
	CreateMap(arr []string) map[string]interface{}
	ChangeTime(source string) time.Time
}

type ExcelStrcut struct {
	temp  [][]string
	Model interface{}
	Info  []map[string]string
}



type Temp struct {
	id         uint64 `gorm:"primary_key"`
	category   string
	origin     uint64
	origintxt  string
	details    string
	code       string
}

var DB *gorm.DB

func saveExcel2DB() {
	db, err := gorm.Open("mysql", "root:xiaodao@/bianma?charset=utf8&parseTime=True&loc=Local")
	//db.SingularTable(true)
	DB = db
	defer db.Close()
	if err != nil {
		panic(err)
	}

	e:=ExcelStrcut{}
	temp := Temp{}
	e.Model=temp
	e.ReadExcel("/Users/iwillcodeu/Downloads/表1 红线问题分类编码表.xlsx").CreateMap().SaveDb(&temp)
}

func (excel *ExcelStrcut)ReadExcel(file string) *ExcelStrcut{

	xlsx, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Println("read excel:",err)
	}

	rows := xlsx.GetRows("Sheet1")
	excel.temp = rows

	return excel

}

func (excel *ExcelStrcut)CreateMap() *ExcelStrcut{

	fmt.Println("ExcelStruct.CreatMap() is called")
	//利用反射得到字段名
	for _,v:=range excel.temp{
		//将数组  转成对应的 map
		var info = make(map[string]string)
		for i:=0;i<reflect.ValueOf(excel.Model).NumField();i++{

			obj:=reflect.TypeOf(excel.Model).Field(i)
			fmt.Printf("key:%s--val:%s\n",obj.Name,v[i])
			info[obj.Name] = v[i]

		}
		excel.Info=append(excel.Info, info)

	}
	return excel
}

func (excel *ExcelStrcut)ChangeTime(source string) time.Time{
	ChangeAfter,err:=time.Parse("2006-01-02", source)
	if err!=nil {
		log.Fatalf("转换时间错误:%s",err)
	}
	return ChangeAfter
}

func (excel *ExcelStrcut)SaveDb(temp *Temp) *ExcelStrcut{

	//忽略标题行
	for i:=1 ;i<len(excel.Info);i++{

		t:=reflect.ValueOf(temp).Elem()
		for k,v:=range excel.Info[i]{

			//fmt.Println(t.FieldByName(k).t.FieldByName(k).Kind())
			fmt.Println("key:%v---val:%v",t.FieldByName(k),t.FieldByName(k).Kind())

			switch t.FieldByName(k).Kind(){
			case reflect.String:
				t.FieldByName(k).Set(reflect.ValueOf(v))
			case reflect.Float64:
				tempV,err:= strconv.ParseFloat(v,64)
				if err != nil{
					log.Printf("string to float64 err：%v",err)
				}

				t.FieldByName(k).Set(reflect.ValueOf(tempV))
			case reflect.Uint64:
				reflect.ValueOf(v)
				tempV, err := strconv.ParseUint(v, 0, 64)
				if err != nil{
					log.Printf("string to uint64 err：%v",err)
				}
				t.FieldByName(k).Set(reflect.ValueOf(tempV))

			case reflect.Struct:
				tempV,err:=time.Parse("2006-01-02", v)
				if err!=nil {
					log.Fatalf("string to time err:%v",err)
				}
				t.FieldByName(k).Set(reflect.ValueOf(tempV))
			default:
				fmt.Println("type err")

			}


		}
		err:=DB.Create(&temp).Error
		if err != nil{
			log.Fatalf("save temp table err:%v",err)
		}
		fmt.Printf("导入临时表成功")

	}
	return excel
}
