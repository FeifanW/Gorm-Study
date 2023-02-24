package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB // 一般挂在这个全局变量上
var mysqlLogger logger.Interface

func init() {
	username := "root"  //账号
	password := ""      // 密码
	host := "127.0.0.1" // 数据库地址，可以是IP或者域名
	port := 3306        // 数据库端口
	Dbname := "gorm"    // 数据库名
	timeout := "10s"    // 连接超时，10s

	//var mysqllogger logger.Interface
	//mysqllogger = logger.Default.LogMode(logger.Info)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username,
		password, host, port, Dbname, timeout)
	// 连接MYSQL，获得DB类型实例，用于后面的数据库读写操作
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "f_", //表名前缀
			SingularTable: true, // 是否单数表名
			NoLowerCase:   true, // 不要小写转换
		},
	})

	if err != nil {
		panic("连接数据库失败，err=" + err.Error())
	}
	DB = db
	// 连接成功
	//fmt.Println(db)
}

type Student struct {
	ID     uint    `gorm:"size:10"`
	Name   string  `gorm:"size:16"`
	Age    int     `gorm:"size:3"`
	Email  *string `gorm:"size:128"`
	Gender bool
	//Type  string  `gorm:"column:_type;size:4"`           // 别名
	//Date  string  `gorm:"default:2022-12-30;comment:日期"` // 默认时间
}

func main() {
	DB.AutoMigrate(&Student{}) // 创建一张表
	//email := "123465@163.com"
	//s1 := Student{
	//	Name:   "hello",
	//	Age:    21,
	//	Gender: true,
	//	Email:  &email,
	//}
	//err := DB.Create(&s1).Error
	//fmt.Println(s1)
	//fmt.Println(err)

	// 批量插入
	//var studentList []Student // 切片
	//for i := 0; i < 10; i++ {
	//	studentList = append(studentList, Student{
	//		Name:   fmt.Sprintf("hello%d", i+1),
	//		Age:    21 + i + 1,
	//		Gender: true,
	//		Email:  &email,
	//	})
	//}
	//err := DB.Create(&studentList).Error
	//fmt.Println(err)
	//fmt.Println(err)

	// 单条记录的查询
	var student Student
	DB = DB.Session(&gorm.Session{
		Logger: mysqlLogger,
	})
	//DB.Take(&student)
	//fmt.Println(student)
	//DB.First(&student)
	//fmt.Println(student)
	//student = Student{}
	//DB.Last(&student)
	//fmt.Println(student)

	// 可以根据主键查询
	DB.Take(&student, 2)
	fmt.Println(student)
	err := DB.Take(&student, "45").Error
	fmt.Println(err)
	fmt.Println(student)

}
