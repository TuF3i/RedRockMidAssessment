package msg

import (
	"RedRockMidAssessment-Synchronizer/core"
	"RedRockMidAssessment-Synchronizer/core/models"
	"context"
	"encoding/json"

	"go.uber.org/zap"
)

func GenUpdateSelectedStuListMsg(ctx context.Context, stuID string, courseID string) []byte {
	jsonBody := models.Commander{
		Role: "course",
		Msg: models.CourseMsg{
			Operation: "subscribe",
			CourseID:  courseID,
			StudentID: stuID,
		},
	}

	jsonBuffer, err := json.Marshal(jsonBody)
	if err != nil {
		core.Logger.Error(
			"Marshal Json Data",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil
	}

	return jsonBuffer
}

func GenUpdateDroppedStuListMsg(ctx context.Context, stuID string, courseID string) []byte {
	jsonBody := models.Commander{
		Role: "course",
		Msg: models.CourseMsg{
			Operation: "drop",
			CourseID:  courseID,
			StudentID: stuID,
		},
	}

	jsonBuffer, err := json.Marshal(jsonBody)
	if err != nil {
		core.Logger.Error(
			"Marshal Json Data",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil
	}

	return jsonBuffer
}

func GenUpdateCourseSelectedNumMsg(ctx context.Context, courseID string, selectedNum uint) []byte {
	jsonBody := models.Commander{
		Role: "selectedNum",
		Msg: models.SelectedNum{
			CourseID:    courseID,
			SelectedNum: selectedNum,
		},
	}

	jsonBuffer, err := json.Marshal(jsonBody)
	if err != nil {
		core.Logger.Error(
			"Marshal Json Data",
			zap.String("snowflake", ctx.Value("trace_id").(string)),
			zap.String("detail", err.Error()),
		)
		return nil
	}

	return jsonBuffer
}
