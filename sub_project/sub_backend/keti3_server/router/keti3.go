package router

import (
	"keti3/internal/application/controller/keti3"
	"keti3/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func registerKeti3Router(router *gin.RouterGroup) {
	router.Use(middleware.CheckSignature(), middleware.CheckUserToken())

	ctrl := keti3.NewCtrl()

	router.POST("/uid/get-by-increase", ctrl.GetUIDByIncrease)
	router.POST("/log/save", ctrl.SaveLog)
	router.GET("/log/export", ctrl.ExportLog)
	router.POST("/oss/auth", ctrl.OssAuth)

	router.POST("/config/list", ctrl.ConfigList)
	router.POST("/student/login", ctrl.StudentLogin)
	router.POST("/teacher/login", ctrl.TeacherLogin)
	router.POST("/log/list", ctrl.LogList)
}
