package service

import (
	dao "RedRockMidAssessment/core/dao/mysql"
	"RedRockMidAssessment/core/models"
	"RedRockMidAssessment/core/utils/response"
	"context"
)

func GetStudentsList(ctx context.Context, page int, resNum int) (models.Students, response.Response) {
	/* 数据校验 */
	if page <= 0 {
		page = 1
	}
	if resNum <= 0 {
		resNum = 10
	}

	/* 计算偏移量 */
	offset := (page - 1) * resNum

	/* 查MySQL */
	return dao.GetStudentList(ctx, resNum, offset, page)
}
