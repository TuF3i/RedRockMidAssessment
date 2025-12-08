package api

import (
	"RedRockMidAssessment/core"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
	monitor "github.com/hertz-contrib/monitor-prometheus"
)

func HertzApi() {
	// 构造Url
	url := fmt.Sprintf("%v:%v", core.Config.HertzAPI.ListenAddr, core.Config.HertzAPI.ListenPort)
	monitorUrl := fmt.Sprintf("%v:%v", core.Config.HertzAPI.ListenAddr, core.Config.HertzAPI.MonitorPort)
	// 创建服务核心
	h := server.Default(server.WithHostPorts(url), server.WithTracer(monitor.NewServerTracer(monitorUrl, "/monitor")))
	// 初始化路由
	initRouter(h)
	// 启动Hertz引擎
	go func() { h.Spin() }()
}

func initRouter(h *server.Hertz) {
	// 注册公共接口路由组
	publicApi := h.Group("/v1/api/public")
	// 注册学生管理接口路由组
	stuManagerForStu := h.Group("/v1/api/stu-manager", JWTAuthMiddleWare())
	// 注册课程管理接口路由组
	courseManagerForStu := h.Group("/v1/api/class-manager", JWTAuthMiddleWare())
	// 注册带有管理员权限的学生管理接口路由组
	stuManagerForAdmin := h.Group("/v1/api/admin/stu-manager", JWTAuthMiddleWare())
	// 注册带有管理员权限的课程管理接口路由组
	courseManagerForAdmin := h.Group("/v1/api/admin/classes-manager", JWTAuthMiddleWare())

	// 注册接口
	publicApi.POST("/register", RegisterHandleFunc())
	// 登录接口
	publicApi.POST("/login", LoginHandleFunc())
	// 刷新Token接口
	publicApi.GET("/refresh", JWTRefreshMiddleWare(), RefreshTokensHandleFunc())

	// 获取学生所有信息接口
	stuManagerForStu.GET("/stu-info", GetStudentInfoForStuHandleFunc())
	// 更新学生指定字段的信息接口
	stuManagerForStu.PATCH("/stu-update", UpdateStudentInfoForStuHandleFunc())
	// 学生注销（退出登录）接口
	stuManagerForStu.GET("/stu-logout", LogoutHandleFunc())

	// 学生查看所有可选课程接口
	courseManagerForStu.GET("/get-selectable-classes", CheckIfCourseSelectionStartedMiddleWare(), GetCourseInfoForStudentHandleFunc())
	// 学生添加选课接口
	courseManagerForStu.POST("/subscribe-class", CheckIfCourseSelectionStartedMiddleWare(), SubscribeCourseForStudentHandleFunc())
	// 学生删除选课接口
	courseManagerForStu.DELETE("/del-class", CheckIfCourseSelectionStartedMiddleWare(), DropCourseForStudentHandleFunc())
	// 学生查看已选课程接口
	courseManagerForStu.GET("/get-subscribed-classes", CheckIfCourseSelectionStartedMiddleWare(), GetSelectedCourseForStudentHandleFunc())

	// 管理员查看学生列表接口
	stuManagerForAdmin.GET("/get-stu-list", CheckIfCourseSelectionStartedMiddleWareForAdmin(), GetStudentListForAdminHandleFunc())
	// 管理员获取学生信息接口
	stuManagerForAdmin.GET("/get-stu-info", CheckIfCourseSelectionStartedMiddleWareForAdmin(), GetStudentInfoForAdminHandleFunc())
	// 管理员更新学生信息
	stuManagerForAdmin.PATCH("/update-stu-info", CheckIfCourseSelectionStartedMiddleWareForAdmin(), UpdateStudentInfoForAdminHandleFunc())
	// 管理员创建学生
	stuManagerForAdmin.POST("/create-stu", CheckIfCourseSelectionStartedMiddleWareForAdmin(), AddStudentForAdminHandleFunc())
	// 管理员删除学生
	stuManagerForAdmin.DELETE("/del-stu", CheckIfCourseSelectionStartedMiddleWareForAdmin(), DelStudentForAdminHandleFunc())

	// 管理员查看选课情况
	courseManagerForAdmin.GET("/get-class-status", GetCourseInfoForAdminHandleFunc())
	// 管理员查看学生选课情况
	courseManagerForAdmin.GET("/get-stu-classes", GetStuCourseSelectionInfoForAdminHandleFunc())
	// 管理员修改学生选课
	courseManagerForAdmin.PATCH("/update-stu-classes", CheckIfCourseSelectionStartedMiddleWareForAdmin(), AddStuCourseSelectionInfoForAdminHandleFunc())
	// 管理员删除学生选课
	courseManagerForAdmin.DELETE("/update-stu-classes", CheckIfCourseSelectionStartedMiddleWareForAdmin(), DelStuCourseSelectionInfoForAdminHandleFunc())
	// 管理员添加课程
	courseManagerForAdmin.POST("/add-course", CheckIfCourseSelectionStartedMiddleWareForAdmin(), AddCourseForAdminHandleFunc())
	// 管理员删除课程
	courseManagerForAdmin.DELETE("/delete-course", CheckIfCourseSelectionStartedMiddleWareForAdmin(), DelCourseForAdminHandleFunc())
	// 管理员更新课程信息
	courseManagerForAdmin.PATCH("/edit-class-info", CheckIfCourseSelectionStartedMiddleWareForAdmin(), UpdateCourseInfoForAdminHandleFunc())
	// 管理员更新课程容量
	courseManagerForAdmin.PATCH("/edit-class-stock", CheckIfCourseSelectionStartedMiddleWareForAdmin(), UpdateCourseStockForAdminHandleFunc())
	// 管理员开始选课
	courseManagerForAdmin.GET("/start-course-select-event", StartCourseSelectionEventForAdminHandleFunc())
	// 管理员结束选课
	courseManagerForAdmin.GET("/stop-course-select-event", StopCourseSelectionEventForAdminHandleFunc())
}
