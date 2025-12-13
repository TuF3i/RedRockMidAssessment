package service

import (
	"RedRockMidAssessment-Consumer/core/dao"
	"fmt"
)

func SubmitCourseForStudent(stuID string, courseID string) error {
	// 检查学生是否存在
	ok, err := dao.CheckIfStudentExist(stuID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("student not exist")
	}

	// 检查课程是否存在
	ok, err = dao.CheckIfCourseExist(courseID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("course not exist")
	}

	// 检查选课记录是否存在
	ok, err = dao.CheckIfCourseSelected(stuID, courseID)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	// 添加选课记录
	err = dao.SubscribeCourseForStu(stuID, courseID)
	if err != nil {
		return err
	}

	return nil
}

func DropCourseForStudent(stuID string, courseID string) error {
	// 检查学生是否存在
	ok, err := dao.CheckIfStudentExist(stuID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("student not exist")
	}

	// 检查课程是否存在
	ok, err = dao.CheckIfCourseExist(courseID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("course not exist")
	}

	// 检查选课记录是否存在
	ok, err = dao.CheckIfCourseSelected(stuID, courseID)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// 删除选课记录
	err = dao.DropCourseForStu(stuID, courseID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateSelectedStuNum(courseID string, num uint) error {
	// 检查课程是否存在
	ok, err := dao.CheckIfCourseExist(courseID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("course not exist")
	}

	// 获取课程容量
	capacity, err := dao.GetCourseCapacity(courseID)
	if err != nil {
		return err
	}

	// 计算课程容量
	selectedNum := capacity - num

	// 更改已选课人数
	err = dao.UpdateSelectedStuNum(courseID, selectedNum)
	if err != nil {
		return err
	}

	return nil
}
