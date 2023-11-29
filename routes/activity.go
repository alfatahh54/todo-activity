package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alfatahh54/todo-activity/db"
	"github.com/alfatahh54/todo-activity/types"
	"github.com/gin-gonic/gin"
)

func init() {
	go MainRoute.NewRoute("GET", "/activity-groups", GetActivity)
	go MainRoute.NewRoute("GET", "/activity-groups/:id", GetActivity)
	go MainRoute.NewRoute("POST", "/activity-groups", CreateActivity)
	go MainRoute.NewRoute("DELETE", "/activity-groups/:id", DeleteActivity)
	go MainRoute.NewRoute("PATCH", "/activity-groups/:id", UpdateActivity)
}

func GetActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	fmt.Println(id)
	if id == "" {
		var Activities []types.ActivityType
		params := db.QueryParams{
			Where: &db.WhereParam{
				Table: &[]string{"a"}[0],
				Field: &[]string{"deleted_at"}[0],
				Op:    &[]string{"IS"}[0],
				Value: &[]any{"NULL"}[0],
			},
		}
		db.MysqlQueryParams("activity", "a", "GetAllActivity", &Activities, params)
		if len(Activities) == 0 {
			ctx.JSON(200, gin.H{
				"status":  "Not Found",
				"message": "Data not found",
				"data":    Activities,
			})
			return
		} else {
			ctx.JSON(200, gin.H{
				"status":  "Success",
				"message": "Success",
				"data":    Activities,
			})
			return
		}
	} else {
		var Activity types.ActivityType
		db.MysqlQuerySingleRow("SELECT * FROM activity WHERE id = ?  AND deleted_at IS NULL;", "GetAllToDo", &Activity, id)
		if Activity.ID == nil {
			ctx.JSON(404, gin.H{
				"status":  "Not Found",
				"message": "Activity with ID " + id + " Not Found",
				"data":    Activity,
			})
			return
		} else {
			ctx.JSON(200, gin.H{
				"status":  "Success",
				"message": "Success",
				"data":    Activity,
			})
			return
		}
	}
}

func CreateActivity(ctx *gin.Context) {
	var body types.ActivityType

	if err := ctx.ShouldBind(&body); err != nil {
		fmt.Println("Error binding body : ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "Bad request",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	nowTime := time.Now()
	if body.CreatedAt == nil {
		body.CreatedAt = &nowTime
	}
	if body.UpdatedAt == nil {
		body.UpdatedAt = &nowTime
	}
	err := db.InsertOrUpdateStruct("activity", &body)
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

func DeleteActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "id is required",
			"data":    nil,
		})
		return
	}
	var Activity types.ActivityType
	db.MysqlQuerySingleRow("SELECT * FROM activity WHERE id = ? AND deleted_at IS NULL;", "GetActivityForDelete", &Activity, id)
	if Activity.ID == nil {
		ctx.JSON(404, gin.H{
			"status":  "Not Found",
			"message": "Activity with ID " + id + " Not Found",
			"data":    Activity,
		})
		return
	} else {
		timeNow := time.Now()
		Activity.DeletedAt = &timeNow
		db.InsertOrUpdateStruct("activity", &Activity)
		ctx.JSON(200, gin.H{
			"status":  "Success",
			"message": "Success",
			"data":    map[string]any{},
		})
		return
	}
}

func UpdateActivity(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "id is required",
			"data":    nil,
		})
		return
	}
	var Activity types.ActivityType
	db.MysqlQuerySingleRow("SELECT * FROM activity WHERE id = ? AND deleted_at IS NULL;", "GetActivityForDelete", &Activity, id)
	if Activity.ID == nil {
		ctx.JSON(404, gin.H{
			"status":  "Not Found",
			"message": "Activity with ID " + id + " Not Found",
			"data":    Activity,
		})
		return
	} else {
		var body types.ActivityUpdateType
		if err := ctx.ShouldBind(&body); err != nil {
			fmt.Println("Error binding body : ", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "Bad request",
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		if body.Email == "" && body.Title == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "Bad Request",
				"message": "title cannot be null",
				"data":    map[string]any{},
			})
			return
		}
		if body.Email != "" {
			Activity.Email = body.Email
		}
		if body.Title != "" {
			Activity.Title = body.Title
		}
		timeNow := time.Now()
		Activity.UpdatedAt = &timeNow
		db.InsertOrUpdateStruct("activity", &Activity)
		ctx.JSON(200, gin.H{
			"status":  "Success",
			"message": "Success",
			"data":    map[string]any{},
		})
	}
}
