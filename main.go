package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	_ "todo-app/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Controller struct {
	MysqlDB *gorm.DB
}

type TodoList struct {
	ID     int    `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" example:"1" `
	Title  string `gorm:"column:title;not null" json:"title" example:"tile"`
	Status bool   `gorm:"column:status;not null" json:"status" example:"true"`
}

type CreateTodoReq struct {
	Title  string `json:"title" example:"tile"`
	Status bool   `json:"status" example:"true"`
}

func NewMysqlDB() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/todoapp?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	db.AutoMigrate(&TodoList{})
	return db
}

//	@title			Todo App
//	@version		1.0
//	@description	This is a sample server Todo App server.

//	@contact.name	API Support
//	@contact.url	https://www.support.comsithiask.com
//	@contact.email	sitthisak.dev@gmail.com.com

//	@host		localhost:1323
//	@BasePath	/
func main() {
	authValidator := func(username string, password string, c echo.Context) (bool, error) {
		if username == "admin" && password == "1234" {
			return true, nil
		}
		return false, nil
	}
	basicAuth := middleware.BasicAuth(authValidator)

	db := NewMysqlDB()
	controller := Controller{MysqlDB: db}

	e := echo.New()

	e.File("/favicon.ico", "./public/favicon.ico")
	e.GET("/", func(c echo.Context) error {
		htmlFile, err := os.ReadFile("public/views/index.html")
		if err != nil {
			fmt.Print(err)
		}
		str := string(htmlFile)
		return c.HTML(http.StatusOK, str)
	}, basicAuth)

	e.GET("/todo", controller.GetTodo)
	e.POST("/todo", controller.CreateTodo)
	e.GET("/todo/:id", controller.GetTodoByID)
	e.PUT("todo/:id", controller.UpdateTodo)
	e.DELETE("todo/:id", controller.DeleteTodo)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":1323"))
}

// Todo App
//
//	@Tags			Todo
//	@Summary		Get all Todo
//	@Description	Get all Todo.
//	@Produce		json
//	@Success		200	{array}	TodoList
//	@Router			/todo [get]
func (ctrl Controller) GetTodo(c echo.Context) error {
	todoList := []TodoList{}
	err := ctrl.MysqlDB.Order("status asc").Find(&todoList).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, todoList)
}

// Todo App
//
//	@Tags			Todo
//	@Summary		Create Todo
//	@Description	Create Todo.
//	@Produce		json
//	@Success		201
//	@Param			data	body	CreateTodoReq	true	"Request payload"
//	@Router			/todo [post]
func (ctrl Controller) CreateTodo(c echo.Context) error {
	u := new(CreateTodoReq)
	if err := c.Bind(u); err != nil {
		return err
	}
	todoList := TodoList{
		Title:  u.Title,
		Status: u.Status,
	}
	err := ctrl.MysqlDB.Create(&todoList).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusCreated)
}

// Todo App
//
//	@Tags			Todo
//	@Summary		Get Todo by ID
//	@Description	Get Todo by ID.
//	@Produce		json
//	@Param			id	path		int	true	"ID"
//	@Success		200	{object}	TodoList
//	@Router			/todo/{id} [get]
func (ctrl Controller) GetTodoByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "id is required!")
	}

	todoList := TodoList{}
	err = ctrl.MysqlDB.Model(TodoList{}).Where(&TodoList{
		ID: id,
	}).Take(&todoList).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, todoList)
}

// Todo App
//
//	@Tags			Todo
//	@Summary		Update Todo by ID
//	@Description	Update Todo by ID.
//	@Produce		json
//	@Param			id	path	int	true	"ID"
//	@Success		201
//	@Param			data	body	CreateTodoReq	true	"Request payload"
//	@Router			/todo/{id} [put]
func (ctrl Controller) UpdateTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "id is required!")
	}

	todoList := TodoList{}
	err = ctrl.MysqlDB.Model(TodoList{}).Where(&TodoList{
		ID: id,
	}).Take(&todoList).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	u := new(CreateTodoReq)
	if err := c.Bind(u); err != nil {
		return err
	}
	todoList = TodoList{
		ID:     id,
		Title:  u.Title,
		Status: u.Status,
	}
	err = ctrl.MysqlDB.Save(&todoList).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusCreated)
}

// Todo App
//
//	@Tags			Todo
//	@Summary		Delete Todo by ID
//	@Description	Delete Todo by ID.
//	@Success		201
//	@Param			id	path	int	true	"ID"
//	@Router			/todo/{id} [delete]
func (ctrl Controller) DeleteTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "id is required!")
	}

	err = ctrl.MysqlDB.Model(TodoList{}).Where(&TodoList{
		ID: id,
	}).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = ctrl.MysqlDB.Delete(&TodoList{}, id).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusCreated)
}
