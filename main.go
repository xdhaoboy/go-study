package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func getUser(c echo.Context) error {
	// User ID 来自于url `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

// e.GET("/show", show)
func show(c echo.Context) error {
	// 从请求参数里获取 team 和 member 的值
	name := c.QueryParam("name")
	age := c.QueryParam("age")
	return c.String(http.StatusOK, "name:"+name+", age:"+age)
}

func test(c echo.Context) error {
	// 从请求参数里获取 team 和 member 的值
	dirPath := "/Users/hao/privateCloud"
	files, dirs, _ := GetFilesAndDirs(dirPath)

	for _, dir := range dirs {
		fmt.Printf("获取的文件夹为[%s]\n", dir)
	}

	for _, table := range dirs {
		temp, _, _ := GetFilesAndDirs(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	for _, table1 := range files {
		fmt.Printf("获取的文件为[%s]\n", table1)
	}

	fmt.Printf("=======================================\n")
	xfiles, _ := GetAllFiles("./simplemath")
	for _, file := range xfiles {
		fmt.Printf("获取的文件为[%s]\n", file)
	}
	return c.String(http.StatusOK, " OK ")
}

//获取指定目录下的所有文件和目录
func GetFilesAndDirs(dirPth string) (files []string, dirs []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetFilesAndDirs(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".go")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	return files, dirs, nil
}

//获取指定目录下的所有文件,包含子目录下的文件
func GetAllFiles(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".go")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllFiles(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	return files, nil
}

func save(c echo.Context) error {
	// Get name
	name := c.FormValue("name")
	// Get avatar
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return err
	}

	// Source
	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(avatar.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "<b>Thank you! "+name+"</b>")
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/users/:id", getUser)
	e.GET("/show", show)
	e.GET("/test", test)
	e.POST("/save", save)
	e.Logger.Fatal(e.Start(":1323"))
}
