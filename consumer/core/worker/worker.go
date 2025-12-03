package worker

import (
	"RedRockMidAssessment-Consumer/core"
	"RedRockMidAssessment-Consumer/core/models"
	"RedRockMidAssessment-Consumer/core/service"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Worker struct {
	blanket chan struct{} // 令牌桶（用于并发控制）
}

func NewConsumer(max int) *Worker {
	return &Worker{blanket: make(chan struct{}, max)}
}

func (c *Worker) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *Worker) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *Worker) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// 拿令牌
		c.blanket <- struct{}{}
		// 异步处理
		go func(m *sarama.ConsumerMessage) {
			var commander models.Commander
			// 释放令牌
			defer func() { <-c.blanket }()
			// 生成snowflakeID
			traceID := core.SnowFlake.TraceID()
			// 解析JSON
			if err := json.Unmarshal(m.Value, &commander); err != nil {
				core.Logger.Error(
					"Unmarshal JSON Byte Error",
					zap.String("snowflake", traceID),
					zap.String("detail", err.Error()),
					zap.String("raw_message", string(m.Value)),
				)
			}
			// 判断操作对象类型
			if commander.Role == "course" {
				msg, ok := commander.Msg.(models.CourseMsg)
				if !ok {
					err := errors.New("can Not Do Type Assertion")
					core.Logger.Error(
						"Type Assertion Error",
						zap.String("snowflake", traceID),
						zap.String("detail", err.Error()),
						zap.String("raw_message", string(m.Value)),
					)
				}
				//判断操作类型
				if msg.Operation == "subscribe" { // 添加选课
					if err := service.SubmitCourseForStudent(msg.StudentID, msg.CourseID); err != nil { // 调用course_service执行命令
						core.Logger.Error(
							"Submit Course Error",
							zap.String("snowflake", traceID),
							zap.String("detail", err.Error()),
							zap.String("raw_message", string(m.Value)),
						)
					}
				}
				if msg.Operation == "drop" { // 退课
					if err := service.DropCourseForStudent(msg.StudentID, msg.CourseID); err != nil { // 调用course_service执行命令
						core.Logger.Error(
							"Drop Course Error",
							zap.String("snowflake", traceID),
							zap.String("detail", err.Error()),
							zap.String("raw_message", string(m.Value)),
						)
					}
				}
			}
			if commander.Role == "selectedNum" {
				msg, ok := commander.Msg.(models.SelectedNum)
				if !ok {
					err := errors.New("can Not Do Type Assertion")
					core.Logger.Error(
						"Type Assertion Error",
						zap.String("snowflake", traceID),
						zap.String("detail", err.Error()),
						zap.String("raw_message", string(m.Value)),
					)
				}
				if err := service.UpdateSelectedStuNum(msg.CourseID, msg.SelectedNum); err != nil {
					core.Logger.Error(
						"Update SelectedStuNum Error",
						zap.String("snowflake", traceID),
						zap.String("detail", err.Error()),
						zap.String("raw_message", string(m.Value)),
					)
				}
			}
			// 输出成功日志
			core.Logger.Info(
				"Success Handle Message",
				zap.String("snowflake", traceID),
				zap.String("offset", strconv.FormatInt(msg.Offset, 10)),
				zap.String("raw_message", string(m.Value)),
			)
			// ACK
			sess.MarkMessage(m, "")
		}(msg)
	}
	return nil
}
