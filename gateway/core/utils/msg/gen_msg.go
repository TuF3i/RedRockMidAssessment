package msg

import (
	"RedRockMidAssessment/core"
	"RedRockMidAssessment/core/models"
	"context"
	"encoding/json"

	"go.uber.org/zap"
)

func GenStartCourseSelection(ctx context.Context) []byte {
	jsonBody := models.Commander{
		Role: "CourseSelectionEvent",
		Msg:  "open",
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

func GenStopCourseSelection(ctx context.Context) []byte {
	jsonBody := models.Commander{
		Role: "CourseSelectionEvent",
		Msg:  "close",
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
