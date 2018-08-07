package main

import (
    "net/http"
    "strings"
    "reflect"
    "log"
    "html/template"
    "io"
    "fmt"
    "time"
    "path"
    "os"
    "github.com/360EntSecGroup-Skylar/excelize"
    "strconv"
    "github.com/ziutek/mymysql/mysql"
)


const (
    upload_path string = "./bianma-uploaded-"
)


func adminHandler(w http.ResponseWriter, r *http.Request) {
    // 获取cookie
    cookie, err := r.Cookie("admin_name")
    if err != nil || cookie.Value == ""{
        http.Redirect(w, r, "/login/index", http.StatusFound)
    }
    
    pathInfo := strings.Trim(r.URL.Path, "/")
    parts := strings.Split(pathInfo, "/")
    var action = ""
    if len(parts) > 1 {
        action = strings.Title(parts[1]) + "Action"
    }
    
    admin := &adminController{}
    controller := reflect.ValueOf(admin)
    method := controller.MethodByName(action)
    if !method.IsValid() {
        method = controller.MethodByName(strings.Title("index") + "Action")
    }
    requestValue := reflect.ValueOf(r)
    responseValue := reflect.ValueOf(w)
    userValue := reflect.ValueOf(cookie.Value)
    method.Call([]reflect.Value{responseValue, requestValue, userValue})
}

func bianmaHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("bianmaHandler")
    // 获取cookie
    cookie, err := r.Cookie("admin_name")
    if err != nil || cookie.Value == ""{
        http.Redirect(w, r, "/login/index", http.StatusFound)
    }

    t, err := template.ParseFiles("template/html/tables/bianma/index.html")
    if (err != nil) {
        log.Println(err)
    }
    t.Execute(w, nil)
}

func bianmaUpload(w http.ResponseWriter, r *http.Request) {
    log.Println("bianmaUpload")
    // 获取cookie
    cookie, err := r.Cookie("admin_name")
    if err != nil || cookie.Value == ""{
        http.Redirect(w, r, "/login/index", http.StatusFound)
    }

    //t, err := template.ParseFiles("template/html/tables/bianma/upload.html")
    //if (err != nil) {
    //    log.Println(err)
    //}
    //t.Execute(w, nil)
    io.WriteString(w, "<html><head><title>xlsx文件导入</title></head><body><form action='' method=\"post\" " +
        "enctype=\"multipart/form-data\"><label>上传EXCEL</label><input type=\"file\" name='file'  /><br/><label>" +
        "<input type=\"submit\" value=\"上传EXCEL\"/></label></form></body></html>")
}

func bianmaUploaded(w http.ResponseWriter, r *http.Request) {
    //获取文件内容 要这样获取
    file, head, err := r.FormFile("file")
    if err != nil {
       fmt.Println(err)
       return
    }

    defer file.Close()
    //当期时间格式化
    filename := time.Now().Format("20180102150405")
    //获取文件的后缀
    fileSuffix := path.Ext(head.Filename)

    filePath := upload_path + filename + fileSuffix
    fmt.Printf("filePaht=%s\n", filePath)
    //创建文件
    fW, err := os.Create(filePath)
    if err != nil {
       fmt.Println("文件创建失败")
       return
    }
    defer fW.Close()
    _, err = io.Copy(fW, file)
    if err != nil {
       fmt.Println("文件保存失败")
       return
    }
    if (fileSuffix == ".xlsx" || fileSuffix == ".xls") {
       fileXlsx(filePath)
    }
    //跳转到首页
    http.Redirect(w, r, "/", http.StatusFound)
}

func fileXlsx(filePath string) {
    xlsx, err := excelize.OpenFile(filePath)
    if err != nil {
       fmt.Println(err)
       os.Exit(1)
    }
    rows := xlsx.GetRows("Sheet1")

    db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "xiaodao", "webdemo")
    if err := db.Connect(); err != nil {
        log.Println(err)
        //OutputJson(w, 0, "数据库操作失败", nil)
        return
    }
    defer db.Close()
    i := 0
    var category, originalTxt, details string
    var original int
    var code float64
    for _, row := range rows[1:] {
       for _, colCell := range row {
           i++
           switch i{
               case 1:{
                   category = colCell
               }
               case 2:{
                   original,_ = strconv.Atoi(colCell)

               }
               case 3:{
                   originalTxt = colCell

               }
               case 4:{
                   details = colCell

               }
               case 5:{
                   code,_ = strconv.ParseFloat(colCell,32)
                   code = code * 10
                   code = float64(int(code))/10
               }
            stmt, err := db.Prepare("insert into bianmaindex values(?, ?, ?,?,?)")
            if err != nil {
                panic(err)
            }
            defer db.Close()
            stmt.Exec(category, original, originalTxt, details, code)

           }
       }
       fmt.Println("")
    }
}


func ajaxHandler(w http.ResponseWriter, r *http.Request) {
    pathInfo := strings.Trim(r.URL.Path, "/")
    parts := strings.Split(pathInfo, "/")
    var action = ""
    if len(parts) > 1 {
        action = strings.Title(parts[1]) + "Action"
    }

    ajax := &ajaxController{}
    controller := reflect.ValueOf(ajax)
    method := controller.MethodByName(action)
    if !method.IsValid() {
        method = controller.MethodByName(strings.Title("index") + "Action")
    }
    requestValue := reflect.ValueOf(r)
    responseValue := reflect.ValueOf(w)
    method.Call([]reflect.Value{responseValue, requestValue})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("loginHandler")
    pathInfo := strings.Trim(r.URL.Path, "/")
    parts := strings.Split(pathInfo, "/")
    var action = ""
    if len(parts) > 1 {
        action = strings.Title(parts[1]) + "Action"
    }

    login := &loginController{}
    controller := reflect.ValueOf(login)
    method := controller.MethodByName(action)
    if !method.IsValid() {
        method = controller.MethodByName(strings.Title("index") + "Action")
    }
    requestValue := reflect.ValueOf(r)
    responseValue := reflect.ValueOf(w)
    method.Call([]reflect.Value{responseValue, requestValue})
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/" {
        http.Redirect(w, r, "/login/index", http.StatusFound)
    }
    
    t, err := template.ParseFiles("template/html/404.html")
    if (err != nil) {
        log.Println(err)
    }
    t.Execute(w, nil)
}