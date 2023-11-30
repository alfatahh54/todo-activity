package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alfatahh54/todo-activity/db"
	"github.com/alfatahh54/todo-activity/types"
	"github.com/gin-gonic/gin"
)

func init() {
	go MainRoute.NewRoute("GET", "/todo", GetTodo)
	go MainRoute.NewRoute("GET", "/todo/:id", GetTodo)
	go MainRoute.NewRoute("POST", "/todo", CreateTodo)
	go MainRoute.NewRoute("DELETE", "/todo/:id", DeleteTodo)
	go MainRoute.NewRoute("PATCH", "/todo/:id", UpdateTodo)
}

func GetTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	fmt.Println(id)
	if id == "" {
		var Todo []types.TodoType
		params := db.QueryParams{
			Where: &db.WhereParam{
				Table: &[]string{"t"}[0],
				Field: &[]string{"deleted_at"}[0],
				Op:    &[]string{"IS"}[0],
				Value: &[]any{"NULL"}[0],
			},
		}
		db.MysqlQueryParams("todo", "t", "GetAllTodo", &Todo, params)
		if len(Todo) == 0 {
			ctx.JSON(200, gin.H{
				"status":  "Not Found",
				"message": "Data not found",
				"data":    Todo,
			})
			return
		} else {
			ctx.JSON(200, gin.H{
				"status":  "Success",
				"message": "Success",
				"data":    Todo,
			})
			return
		}
	} else {
		var Todo types.TodoType
		db.MysqlQuerySingleRow("SELECT * FROM todo WHERE id = ?  AND deleted_at IS NULL;", "GetAllToDo", &Todo, id)
		if Todo.ID == nil {
			ctx.JSON(404, gin.H{
				"status":  "Not Found",
				"message": "Todo with ID " + id + " Not Found",
				"data":    Todo,
			})
			return
		} else {
			ctx.JSON(200, gin.H{
				"status":  "Success",
				"message": "Success",
				"data":    Todo,
			})
			return
		}
	}
}

func CreateTodo(ctx *gin.Context) {
	var body types.TodoType
	if err := ctx.ShouldBind(&body); err != nil {
		fmt.Println("Error binding body : ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "Bad request",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	var Activity types.ActivityType
	db.MysqlQuerySingleRow("SELECT * FROM activity WHERE id = ? AND deleted_at IS NULL;", "GetActivityForDelete", &Activity, body.ActivityGroupID)
	if Activity.ID == nil {
		activity_id := strconv.Itoa(body.ActivityGroupID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "Bad request",
			"message": "Activiti with ID " + activity_id + " not found",
			"data":    nil,
		})
	}
	nowTime := time.Now()
	if body.CreatedAt == nil {
		body.CreatedAt = &nowTime
	}
	if body.UpdatedAt == nil {
		body.UpdatedAt = &nowTime
	}
	if body.IsActive == false {
		body.IsActive = true
	}
	err := db.InsertOrUpdateStruct("todo", &body)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"status":  "Error",
			"message": "Internal Server Error",
			"data":    nil,
		})
	} else {
		ctx.JSON(200, gin.H{
			"status":  "Success",
			"message": "Success",
			"data":    body,
		})
	}
}

func DeleteTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "id is required",
			"data":    nil,
		})
		return
	}
	var Todo types.TodoType
	db.MysqlQuerySingleRow("SELECT * FROM todo WHERE id = ? AND deleted_at IS NULL;", "GetTodoForDelete", &Todo, id)
	if Todo.ID == nil {
		ctx.JSON(404, gin.H{
			"status":  "Not Found",
			"message": "Activity with ID " + id + " Not Found",
			"data":    Todo,
		})
		return
	} else {
		timeNow := time.Now()
		Todo.DeletedAt = &timeNow
		db.InsertOrUpdateStruct("todo", &Todo)
		ctx.JSON(200, gin.H{
			"status":  "Success",
			"message": "Success",
			"data":    map[string]any{},
		})
		return
	}
}

func UpdateTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "id is required",
			"data":    nil,
		})
		return
	}
	var Todo types.TodoType
	db.MysqlQuerySingleRow("SELECT * FROM todo WHERE id = ? AND deleted_at IS NULL;", "GetTodoForUpdate", &Todo, id)
	if Todo.ID == nil {
		ctx.JSON(404, gin.H{
			"status":  "Not Found",
			"message": "Activity with ID " + id + " Not Found",
			"data":    Todo,
		})
		return
	} else {
		var body types.TodoUpdateType
		if err := ctx.ShouldBind(&body); err != nil {
			fmt.Println("Error binding body : ", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "Bad request",
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		if body.IsActive == nil && body.Title == "" && body.Priority == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "Bad Request",
				"message": "title cannot be null",
				"data":    map[string]any{},
			})
			return
		}
		if body.Priority != "" {
			Todo.Priority = body.Priority
		}
		if body.Title != "" {
			Todo.Title = body.Title
		}
		if body.IsActive != nil {
			Todo.IsActive = *body.IsActive
		}
		timeNow := time.Now()
		Todo.UpdatedAt = &timeNow
		db.InsertOrUpdateStruct("activity", &Todo)
		ctx.JSON(200, gin.H{
			"status":  "Success",
			"message": "Success",
			"data":    map[string]any{},
		})
	}
}
