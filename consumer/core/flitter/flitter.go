package flitter

import (
	"RedRockMidAssessment-Consumer/core"
	"RedRockMidAssessment-Consumer/core/models"
	"RedRockMidAssessment-Consumer/core/service"
	"context"
	"errors"

	"go.uber.org/zap"
)

type HandleFunc func(ctx context.Context, commander models.Commander)

type Business struct {
	Do HandleFunc
}

func DefaultHandleFunc(ctx context.Context, commander models.Commander) {
	core.Logger.Error(
		"Unknow Operation Object",
		zap.String("snowflake", ctx.Value("traceID").(string)),
		zap.String("detail", commander.Role),
	)
}

func CourseHandleFunc(ctx context.Context, commander models.Commander) {
	// 类型断言
	msg, ok := commander.Msg.(models.CourseMsg)
	if !ok {
		err := errors.New("can Not Do Type Assertion")
		core.Logger.Error(
			"Type Assertion Error",
			zap.String("snowflake", ctx.Value("traceID").(string)),
			zap.String("detail", err.Error()),
		)
		return
	}
	// 判断操作类型
	switch msg.Operation {
	case "subscribe":
		// 调用course_service执行操作
		if err := service.SubmitCourseForStudent(msg.StudentID, msg.CourseID); err != nil { // 调用course_service执行命令
			core.Logger.Error(
				"Submit Course Error",
				zap.String("snowflake", ctx.Value("traceID").(string)),
				zap.String("detail", err.Error()),
			)
			return
		}
		core.Logger.Info(
			"Success Handle Message",
			zap.String("snowflake", ctx.Value("traceID").(string)),
			zap.String("op_obj", commander.Role),
			zap.String("op", msg.Operation),
			zap.Any("data", msg),
		)
	case "drop":
		// 调用course_service执行操作
		if err := service.DropCourseForStudent(msg.StudentID, msg.CourseID); err != nil { // 调用course_service执行命令
			core.Logger.Error(
				"Drop Course Error",
				zap.String("snowflake", ctx.Value("traceID").(string)),
				zap.String("detail", err.Error()),
			)
			return
		}
		core.Logger.Info(
			"Success Handle Message",
			zap.String("snowflake", ctx.Value("traceID").(string)),
			zap.String("op_obj", commander.Role),
			zap.String("op", msg.Operation),
			zap.Any("data", msg),
		)
	default:
		// 未知操作类型
		core.Logger.Error(
			"Unknow Operation",
			zap.String("snowflake", ctx.Value("traceID").(string)),
			zap.String("detail", msg.Operation),
			zap.Any("data", msg),
		)
	}
}

func SelectedNumHandleFunc(ctx context.Context, commander models.Commander) {
	// 类型断言
	msg, ok := commander.Msg.(models.SelectedNum)
	if !ok {
		err := errors.New("can Not Do Type Assertion")
		core.Logger.Error(
			"Type Assertion Error",
			zap.String("snowflake", ctx.Value("traceID").(string)),
			zap.String("detail", err.Error()),
		)
		return
	}
	// 调用course_service执行操作
	if err := service.UpdateSelectedStuNum(msg.CourseID, msg.SelectedNum); err != nil {
		core.Logger.Error(
			"Update SelectedStuNum Error",
			zap.String("snowflake", ctx.Value("TraceID").(string)),
			zap.String("detail", err.Error()),
		)
		return
	}
	core.Logger.Info(
		"Success Handle Message",
		zap.String("snowflake", ctx.Value("traceID").(string)),
		zap.String("op_obj", commander.Role),
		zap.Any("data", msg),
	)
}

func GetRelatedHandleFunc(commander models.Commander) *Business {
	// 判断操作对象类型
	if commander.Role == "course" {
		return &Business{Do: CourseHandleFunc}
	}

	if commander.Role == "selectedNum" {
		return &Business{Do: SelectedNumHandleFunc}
	}

	return &Business{Do: DefaultHandleFunc}
}
