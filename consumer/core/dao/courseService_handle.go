package dao

import (
	"RedRockMidAssessment-Consumer/core"
	"RedRockMidAssessment-Consumer/core/models"
	"errors"

	"gorm.io/gorm"
)

func SubscribeCourseForStu(stuID string, courseID string) error {
	var course models.Course
	var student models.Student

	// 开启数据库事务
	tx := core.MysqlConn.Begin()
	// 查询学生
	if err := tx.Where("student_id = ?", stuID).First(&student).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 查询课程
	if err := tx.Where("class_id = ?", courseID).First(&course).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 建立关联
	if err := tx.Model(&student).Association("Courses").Append(&course); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func DropCourseForStu(stuID string, courseID string) error {
	var course models.Course
	var student models.Student

	// 开启数据库事务
	tx := core.MysqlConn.Begin()
	// 查询学生
	if err := tx.Where("student_id = ?", stuID).First(&student).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 查询课程
	if err := tx.Where("class_id = ?", courseID).First(&course).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 建立关联
	if err := tx.Model(&student).Association("Courses").Delete(&course); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func CheckIfCourseExist(courseID string) (bool, error) {
	// 开启数据库事务
	tx := core.MysqlConn.Begin()
	defer tx.Commit()
	// 查询Course是否存在
	err := tx.Where("class_id = ?", courseID).Find(&models.Course{}).Error
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}

func CheckIfCourseSelected(stuID string, courseID string) (bool, error) {
	// 开启数据库事务
	tx := core.MysqlConn.Begin()
	defer tx.Commit()
	// 查询Course是否已被选择
	err := tx.Where("StuID = ? AND CouID = ?", stuID, courseID).Find(&models.Relation{}).Error
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}

func CheckIfStudentExist(stuID string) (bool, error) {
	// 开启数据库事务
	tx := core.MysqlConn.Begin()
	defer tx.Commit()
	// 查询Student是否存在
	err := tx.Where("student_id = ?", stuID).Find(&models.Student{}).Error
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	default:
		return false, err
	}
}

func UpdateSelectedStuNum(courseID string, num uint) error {
	// 开启数据库
	tx := core.MysqlConn.Begin()
	// 更新字段
	err := tx.Model(&models.Course{}).Where("class_id = ?", courseID).Update("class_selected_num", num).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func GetCourseCapacity(courseID string) (uint, error) {
	var capacity uint
	// 开启数据库事务
	tx := core.MysqlConn.Begin()
	defer tx.Commit()
	// 查询课程容量
	if err := tx.Where("class_id = ?", courseID).Pluck("class_capacity", &capacity).Error; err != nil {
		return 0, err
	}

	return capacity, nil
}
