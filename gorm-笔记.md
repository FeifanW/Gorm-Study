#### 一、Gorm简介

##### 什么是ORM：

Object Relational Mapping  “对象关系映射”，解决了对象和关系型数据库之间的数据交互问题

简单理解就是**使用一个类表示一张表，类中的属性表示表的字段，类的实例化对象表示一条记录**

使用对象的方法操作数据库

和自动生成SQL语句相比，手动编写SQL语句缺点很明显：

- 对象的属性名和数据表的字段名往往不一致，在编写SQL语句时要非常小新，要逐一核对属性名和字段名，确保它们不会出错，而且彼此之间要一一对应
- 当SQL语句出错时，数据库的提示信息往往不精准，给拍错带来很多困难
- 不同的数据库，对应的sql语句也不太一样
- sql注入问题

##### ORM的缺点：

- 自动生成SQL语句会消耗计算资源，会对程序员性能产生影响
- 复杂的数据库操作，ORM难以处理，即使可以，自动生成的SQL语句在性能方面也不如原生的SQL
- 生成SQL语句的过程是自动进行的，不能人工干预，所以无法定制一些特殊的SQL语句

每一门语言都有对应的ORM框架：

| 语言   | ORM框架               |
| ------ | --------------------- |
| python | SQLAlchemy、DjangoORM |
| Java   | Hibernate、Mybatis    |
| Golang | GORM                  |

#### 二、操作

##### Gorm命名策略：

需要下载mysql的驱动：

```go
go get gorm.io/driver/mysql
go get gorm.io/gorm
```

gorm命名策略是，表名是蛇形复数，字段名是蛇形单数

```shell
#连接数据库
mysql -u root -p   # -p输入密码
#查看数据库
show databases;
#创建一个gorm数据库
create database gorm;
#使用这个数据库
use gorm;
```

在项目中连接mysql

```go

```

##### 高级配置：

###### 跳过默认事务：

为了确保数据的一致性，GORM会在事务里面执行写入操作（创建、更新、删除），如果没有这方面的要求，可以在初始化时禁用它，这样可以获得60%的性能提升

```go
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    SkipDefaultTransaction: false,
})
```

###### 显示日志：

gorm的默认日志是只打印错误和慢SQL

日志使用有三种，一种是全局、另一种是Session还有一种就是Debug

可以自己设置

```go
var mysqlLogger logger.Interface
// 要显示的日志等级
mysqlLogger = logger.Default.LogNode(logger.Info)
db,err := gorm.Open(mysql.Open(dsn),&gorm.Config{
    Logger:mysqlLogger,
})
```

部分展示日志：

```go
var model Student
session := DB.Session(&gorm.Session{Logger:newLogger})
session.First(&model)
```

如果只想某些语句显示日志：

```go
DB.Debug().First(&model)
```

```go
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB // 一般挂在这个全局变量上

func init() {
	username := "root"       //账号
	password := "" // 密码
	host := "127.0.0.1"      // 数据库地址，可以是IP或者域名
	port := 3306             // 数据库端口
	Dbname := "gorm"         // 数据库名
	timeout := "10s"         // 连接超时，10s

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
	fmt.Println(db)
}

type Student struct {
	ID   uint
	Name string
	Age  int
}

func main() {
	DB.AutoMigrate(&Student{}) // 创建一张表
}
```

##### 模型定义：

模型是标准的struct，由Go的基本数据类型，实现了Scanner和Valuer接口的自定义类型及其指针或别名组成

定义一张表

```go
type Student struct {
	ID uint //默认使用ID作为主键
    Name string
    Email *string // 使用指针是为了存空值
}
```

常识：小写属性是不会生成字段的

###### 自动生成表结构

```go
// 可以放多个
DB.AutoMigrate(&Student{})
```

`AutoMigrate`的逻辑只新增，不删除，不修改（大小会修改）

比如将Name修改为Name1，进行迁移，会多出一个name1字段

###### 修改大小：

我们可以使用gorm的标签进行修改，有两种方式

```go
Name string `gorm:"type:varchar(12)"`
Name string `gorm:"size:2"`
```

###### 字段标签：

| 字段           | 用途             |
| -------------- | ---------------- |
| type           | 定义字段类型     |
| size           | 定义字段大小     |
| column         | 自定义列名       |
| primaryKey     | 将列定义为主键   |
| unique         | 将列定义为唯一键 |
| default        | 定义列的默认值   |
| not null       | 不可为空         |
| embedded       | 嵌套字段         |
| embeddedPrefix | 嵌套字段前缀     |
| comment        | 注释             |

多标签之前用`；`连接

```go
type Student struct {
	ID    uint    `gorm:"size:10"`
	Name  string  `gorm:"size:16"`
	Age   int     `gorm:"size:3"`
	Email *string `gorm:"size:128"`
	Type  string  `gorm:"column:_type;size:4"`           // 别名
	Date  string  `gorm:"default:2022-12-30;comment:日期"` // 默认时间
}

func main() {
	DB.AutoMigrate(&Student{}) // 创建一张表
}
```

##### 单表插入：

使用gorm对单张表进行增删改查

添加记录

```go
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
	//DB.AutoMigrate(&Student{}) // 创建一张表
	email := "123465@163.com"
	s1 := Student{
		Name:   "hello",
		Age:    21,
		Gender: true,
		Email:  &email,
	}
	err := DB.Create(&s1).Error
	fmt.Println(s1)
	fmt.Println(err)
}
```

有两个地方需要注意：

- 指针类型是为了更好的存null类型，但是传值的时候，也记得传指针
- Create接收的是一个指针，而不是值

由于我们传递的是一个指针，调用完Create之后，student这个对象上面就有该记录的信息了，如创建的id

```go
DB.Create(&student)
fmt.Printf("%#v\n",student)
```

##### 批量插入：

```go
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB // 一般挂在这个全局变量上

func init() {
	username := "root"       //账号
	password := "" // 密码
	host := "127.0.0.1"      // 数据库地址，可以是IP或者域名
	port := 3306             // 数据库端口
	Dbname := "gorm"         // 数据库名
	timeout := "10s"         // 连接超时，10s

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
	fmt.Println(db)
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
	email := "123465@163.com"
	// 批量插入
	var studentList []Student // 切片
	for i := 0; i < 10; i++ {
		studentList = append(studentList, Student{
			Name:   fmt.Sprintf("hello%d", i+1),
			Age:    21 + i + 1,
			Gender: true,
			Email:  &email,
		})
	}
	err := DB.Create(&studentList).Error
	fmt.Println(err)
	fmt.Println(err)
}
```

##### 单表查询：

```go
var mysqlLogger logger.Interface
func main() {
	DB.AutoMigrate(&Student{}) // 创建一张表
	// 单条记录的查询
	var student Student
	DB = DB.Session(&gorm.Session{
		Logger: mysqlLogger,
	})
	DB.Take(&student)
	fmt.Println(student)
	DB.First(&student)
	fmt.Println(student)
	student = Student{}
	DB.Last(&student)
	fmt.Println(student)
}
```

###### 根据主键查询：

```go
// 可以根据主键查询
DB.Take(&student, 2)
fmt.Println(student)
err := DB.Take(&student, "45").Error
fmt.Println(err)
fmt.Println(student)
```

###### 根据其他条件查询：

```go
var student Student
DB.Take(&student,"name=?","机器人")
fmt.Println(student)
```

使用？作为占位符，将查询的内容放入？

```mysql
SELECT * FROM `students` WHERE name = '机器人' LIMIT 1
```

这样可以有效的防止sql注入

他的原理就是将参数全部转义，如

```go
DB.Take(&student,"name=?","机器人' or 1=1;#" )  // 感觉这里引号有错误哦
SELECT * FROM `students` WHERE name = '机器人\' or 1=1;#' LIMIT 1
```

###### 根据结构体查询：

```go
var student Student
// 只能有一个主要值
student.ID = 2
DB.Take(&student)
fmt.Println(student)
```

###### 获取查询结果：

```go
// 获取查询的记录数
count := DB.Find(&studentList).RowsAffected
// 是否查询失败
err := DB.Find(&studentList).Error
// 查询失败有查询为空，查询条件错误，sql语法错误可以使用判断
var student Student
err := DB.Take(&student,"xx").Error
switch err {
    case gorm.ErrRecordNotFount:
    fmt.Println("没有找到")
    default:
    fmt.Println("sql错误")
}
```

###### 查询多条记录：

```go
var studentList []Student
DB.Find(&studentList)
```

```go
// 查询多条记录
var studentList []Student
DB.Find(&studentList)
for _, student := range studentList {
    fmt.Println(student)
}

// 由于email是指针类型，所以看不到实际的内容
// 但是序列化之后，会转换为我们看得懂的方式
var studentList1 []Student
DB.Find(&studentList1)
for _, student := range studentList1 {
    data, _ := json.Marshal(student)
    fmt.Println(string(data))
}
```

###### 单表修改和删除：

更新的前提是先查询到记录

- Save保存所有字段

  用于单个记录的全字段更新

  它会保存所有字段，即使零值也会保存

  ```go
  var student Student
  DB.Take(&student,20)  // 更新主键20的
  student.Age = 23
  // 全字段更新
  DB.Save(&student)
  ```

  零值也会更新

  ```go
  var student Student
  DB.Take(&student)
  student.Age = 0
  ```

- 更新指定字段

  可以使用select选择要更新的字段

  ```go
  var student Student
  DB.Take(&student)
  student.Age = 21
  // 全字段更新
  DB.Select("age").Save(&student)
  ```

- 批量更新

  例如给年龄21的学生，更新一下邮箱

  ```go
  var studentList []Student
  DB.Find(&student{},[]int{12,13,14}).Update("eamil","is21@qq.com")
  DB.Find(&student{}).Where("age=?",21).Update("eamil","is21@qq.com")
  ```

- 更新多例

  如果是结构体，它默认不会更新零值

  ```go
  email := "xxx@qq.com"
  DB.Model(&student{}).Where("age=?",21).Updates(Student{
      Email:&email,
      Gender:false,   // 这个不会更新
  })
  ```

  如果想让它更新零值，用select就好

  ```go
  email := "xxx@qq.com"
  DB.Model(&Student{}).Where("age=?",21).Select("gender","email").Updates(Student{
      Email:&email,
      Gender:false,
  })
  ```

  如果不想多写几行代码，可以使用map

  ```go
  DB.Model(&student{}).Where("age=?",21).Updates(map[string]any{
      "email":&email,
      "gender":false
  })
  ```

- 删除

  ```go
  var student Student
  DB.Delete(&student,14)
  DB.Delete(&student,[]int{12,13})
  ```

##### Hook:

在插入一条记录到数据库的时候，希望做些事情

```go
func (user *Student) BeforeCreate(tx *gorm.DB) (err error) {
	email := "test@qq.com"
	user.Email = &email
	return nil
}

func main() {
	DB.Create(&Student{
		Name: "world",
		Age:  25,
	})
}
```

##### 高级查询：

###### Where:

等价于sql语句中的where

```go
var users []Student
// 查询用户名是test
db.Where("name = ?", "test").First(&user)
// 查询用户名不是test
db.Where("name <> ?", "test").Find(&users)
// 查询用户名包含  如燕、李元芳的
db.Where("name in ?", []string{"如燕", "李元芳"}).Find(&users)
// 查询姓李的
db.Where("name like ?", "李%").Find(&users)
// 查询年龄大于23，是qq邮箱的
db.Where("age > ? AND email like ?", "23", "%@qq.com").Find(&users)
// 查询大于某段时间的
db.Where("updated_at > ?", lastWeek).Find(&users)
// SELECT * FROM users WHERE updated_at > '2000-01-01 00:00:00';
// 查询是qq邮箱的，或者是女的
db.Where("gender = ? or email like ?", false, "%@qq.com").Find(&users)
```

###### Not:

```go
// 和where中的not等价
// 排除年龄大于23的
DB.Not("age > 23").Find(&users)
```

###### Or：

```go
// 和where中的or等价
DB.Or("gender = ?",false).Or("email like ?","%@qq.com").Find(&users)
```

###### 排序：

根据年龄倒叙

```go
var users []Student
DB.Order("age desc").Find(&users)
```

注意order的顺序

###### 分页查询：

```go
var users []Student
// 一页两条，第1条
DB.Limit(2).Offset(0).Find(&users)
// 第2页
DB.Limit(2).Offset(2).Find(&users)
// 第3页
DB.Limit(2).Offset(4).Find(&users)
```

通用写法

```go
var users []Student
// 一页多少条
limit := 2
page := 1
offset := (page - 1)*limit
DB.Limit(limit).Offset(offset).Find(&users)
```

###### 去重：

```go
var ageList []int
DB.Table("students").Select("age").Distinct("age").Scan(&ageList)
// 或者
DB.Table("students").Select("distinct age").Scan(&ageList)
```

###### 分组查询和原生SQL:

```go
var ageList []int
// 查询男生的个数和女生的个数
DB.Table("students").Select("count(id)").Group("gender").Scan(&ageList)
```

想精确一点，具体男生的名字女生的名字

```go
type AggeGroup struct {
    Gender int
    Count int `gorm:"column:count(id)"`
    Name string `gorm:"column:group_concat(name)"`
}
var agge []AggeGroup
// 查询男生的个数和女生的个数
DB.Table("student").Select("count(id)","gender","group_concat(name)").Group("gender").Scan(&agge)
```

执行原生SQL

```go
type AggeGroup struct {
    Gender int
    Count int `gorm:"column:count(id)"`
    Name string `gorm:"column:group_concat(name)"`
}
var agge []AggeGroup
DB.Raw("SELECT count(id),gender, group_concat(name) From students GROUP BY gender").Scan(&agge)
```

###### 子查询：

```mysql
select * from students where age >(select avg(age) from students)
```

使用gorm编写

```go
var users []Student
DB.Model(Student{}).Where("age > (?)",DB.Model(Student{}).Select("avg(age)")).Find(&users)
```

###### 命名参数:

```go
var users []Student
DB.Where("name = @name and age = @age",sql.Named("name","hello"),sql.Name("age",23)).Find(&users)
DB.Where("name = @name and age = @age",map[string]any{"name":"hello","age":23}).Find(&users)
```

###### find到map:

```go
var res []map[string]any
DB.Table("students").Find(&res)
```

###### 查询引用Scope：

可以再model层写一些通用的查询方式，这样外界可以直接调用方法即可

```go
func Age23(db *gorm.DB) *gorm.DB{
    return db.Where("age > ?",23)
}
func main(){
    var users []Student
    DB.Scopes(Age23).Find(&users)
}
```

##### 一对多关系：

在gorm中，官方把一对多关系分类了两类

Belongs To  属于谁

Has Mant  我拥有的

```go
// 以用户和文章为例
// 一个用户可以发布多篇文章，一篇文章属于一个用户
type User struct {
    ID uint `gorm:"size:4"`
    Name string `gorm:"size:8"`
    Articles []Article  // 用户拥有的文章列表
}
type Article struct {
    ID uint `gorm:"size:4"`
    Title string `gorm:"size:16"`
    UserID uint // 属于  这里的类型和引用的外键类型一致，包括大小
    User user // 属于
}
```

外键命名，外键名称就是关联表名 + ID，类型是uint

##### 重写外键关联：

```go
type User struct {
    ID uint `gorm:"size:4"`
    Name string `gorm:"size:8"`
    Articles []Article `gorm:"foreignKey:UserID1"` // 用户拥有的文章列表
}
type Article struct {
    ID uint `gorm:"size:4"`
    Title string `gorm:"size:16"`
    UserID1 uint 
    User user `gorm:"foreignKey:UserID1"`// 属于
}
```

##### 重写外键引用：

```go
type user struct {
    ID uint `gorm:"size:4"`
    Name string `gorm:"szie:8"`
    Article []Atricle `gorm:"foreignKey:UserNamel;references:Name"`  // 用户拥有的文章列表
}
type Article struct {
    ID uint `gorm:"size:4"`
    Title string `gorm:"size:16"`
    UserName string
    User User `gorm:"reference:Name"`
}
```

##### 一对多的添加：

创建用户，并且创建文章

```go
a1 := Article{Title:"python"}
a2 := Article{Title:"golang"}
user := User{Name:"hello",Articles:[]Article{a1,a2}}
DB.Create(&user)
```

gorm自动创建了两篇文章，以及创建了一个用户，还将他们的关系给关联上了

创建文章，关联已有用户

```
a1 := Article{Title:"go基础入门",UserID:1}
DB.Create(&a1)
var user User
DB.Take(&user,1)
DB.Create(&Article{Title:"python基础入门",User:user})
```

##### 外键添加：

给现有用户绑定文章

```go
var user User
DB.Take(&user,2)
var article Article
DB.Take(&article,5)
user.Articles = []Article{article}
DB.Save(&user)
```

也可以用Append方法

```go
var user User
DB.Take(&user,2)
var article Article
DB.Take(&article,5)
DB.Model(&user).Association("Articles").Append(&article)
```

##### 一对多关系的查询和删除：

###### 自定义预加载：

```go
var user User
DB.Preload("Articles",func(db *gotm.DB)) *gorm.DB {
    return db.Where("id in ?",[]int{1,2})
}.Take(&user,1)
```

###### 删除：

- 级联删除

  删除用户，与用户关联的文章也会删除

  ```go
  var user User
  DB.Take(&user,1)
  DB.Select("Articles").Delete(&user)
  ```

- 清除外键关系

  删除用户，与将用户关联的文章，外键设置为null

  ```go
  var user User
  DB.Preload("Articles").Take(&user,2)
  DB.Model(&user).Association("Article").Delete(&user.Articles)
  ```

##### 一对一关系：

一对一关系比较少，一般用于表的扩展

比如一张用户表，有很多字段，可以拆成两张表，常用字段放主表，不常用的字段放详情表

##### 多对多的关系

需要第三张表存储两张表的关系

表结构搭建：

```go
type Tag struct {
    ID uint
    Name string
    Articles []Article `gorm:"many2many:article_tags;"`   // 用于反向引导
}
type Article struct {
    ID uint
    Title string
    Tags []Tag `gorm:"many2many:article_tags;"`
}
```

###### 多对多添加

添加文章并创建多标签

```go
DB.Create(&Article{
    Title:"python基础课程"
    Tags:[]Tag{
        {Name:"python"},
        {Name:"基础课程"},
    },
})
```

添加文章，选择标签

```go
var Tags []Tag
DB.Find(&tags,"name = ?","基础课程")
DB.Create(&Article{
	Title:"Go基础",
	Tags:tags,
})
```

###### 多对多查询

查询文章，显示文章的标签列表

```go
var article Article
DB.Preload("Tags").Take(&article,1)
```

查询标签，显示文章列表

```go
var tag Tag
DB.Preload("Acticles").Take(&tag,2)
```

多对多更新

```
var article Article
DB.Preload("Tags").Take(&article,1)
DB.Model(&article).Association("Tags").Delete(article.Tags)
```

更新文章的标签

```go
var article Article
var tags []Tag
DB.Find(&tags,[]int{2,6,7})
DB.Preload("Tags").Take(&article,2)
DB.Model(&article).Association("Tags").Replace(tags)
```

###### 自定义连接表：

默认的链接表，只有双方的主键id，展示不了更多信息

```go
type Article struct {
    ID uint
    Title string
    Tags []Tag `gorm:"many2many:article_tags"`
}
type Tag struct {
    ID uint
    Name string
}
type Article struct {
    ArticleID uint `gorm:"primaryKey"`
    TagID uint `gorm:"primaryKey"`
    CreatedAt time.Time
}
```

生成表结构

```go
// 设置Article的Tags为ArticleTag
DB.SetupJoinTable(&Article{},"Tags",&ArticleTag{})
// 如果Tag要反向应用Article,那么也得加上
err := DB.AutoMigrate(&Article{},&Tag{}，&ArticleTag{})
```

###### 自定义连接表主键：

主要是修改的这两项

joinForeignKey 连接主键的id

joinReferences  关联的主键的id

###### 查询多对多链接表：

##### 自定义数据类型：

存储json需要定义一个结构体，在入库的时候，把它转换成[]byte类型，查询的时候再转换成结构体

存储数据，最简单的方式就是存json

或者用字符串拼接

##### 枚举：

1.0：

很多时候，我们会对一些状态进行判断，而这些状态都是有限的，在主机管理中状态有Running、OffLine、Except异常

如果存储字符串，不仅是浪费空间，每次判断还要多复制很多字符串，主要是后期维护麻烦

2.0

想到使用数字表示状态























