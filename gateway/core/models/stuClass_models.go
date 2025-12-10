package models

import (
	"time"

	"gorm.io/gorm"
)

//查数据库 + data字段返回

type Relation struct {
	// 关联表
	gorm.Model
	CouID string `gorm:"column:CouID"`
	StuID uint   `gorm:"column:StuID"`
}

func (Relation) TableName() string {
	return "relation"
}

type Course struct {
	// 课程表
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `json:"create_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	ClassName        string `json:"class_name" gorm:"column:class_name; not null; type:varchar(90)"`
	ClassID          string `json:"class_id" gorm:"index; unique; not null; type:varchar(60); column:class_id"`
	ClassLocation    string `json:"class_location" gorm:"column:class_location; not null; type:varchar(90)"`
	ClassTime        string `json:"class_time" gorm:"column:class_time; not null; type:varchar(50)"`
	ClassTeacher     string `json:"class_teacher" gorm:"column:class_teacher; not null; type:varchar(30)"`
	ClassCapacity    uint   `json:"class_capacity" gorm:"column:class_capacity; not null; type:int"`
	ClassSelectedNum uint   `json:"class_selection" gorm:"column:class_selected_num; not null; type:int"`

	Students []Student `json:"students" gorm:"many2many:relation; foreignKey:ClassID; joinForeignKey:CouID; references:StudentID; joinReferences:StuID; constraint:OnDelete:CASCADE"`
}

func (Course) TableName() string {
	return "course"
}

type Student struct {
	// 学生表
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `json:"create_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Role         bool   `json:"role" gorm:"column:role; not null; type:tinyint; default:1"` // 0为admin, 1为学生
	Name         string `json:"name" gorm:"column:name; not null; type:varchar(90)"`
	StudentID    string `json:"stu_id" gorm:"column:student_id; index; unique; not null; type:bigint"`
	StudentClass string `json:"stu_class" gorm:"column:student_class; not null; type:varchar(40)"`
	Password     string `json:"password" gorm:"column:password; not null; type:varchar(90)"`
	Sex          uint   `json:"sex" gorm:"column:sex; not null"`
	Grade        uint   `json:"grade" gorm:"column:grade; not null"`
	Age          uint   `json:"age" gorm:"column:age; not null"`

	Courses []Course `json:"courses" gorm:"many2many:relation; foreignKey:StudentID; joinForeignKey:StuID; references:ClassID; joinReferences:CouID; constraint:OnDelete:CASCADE"`
}

func (Student) TableName() string {
	return "student"
}

// 更新学生信息column表
type UpdateDataForStu struct {
	StudentID     string                      `json:"stu_id"`
	UpdateColumns []UpdateColumnsEntityForStu `json:"update_columns"`
}

type UpdateColumnsEntityForStu struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

type Students struct {
	Total        int64                `json:"total"`
	Page         int                  `json:"page"`
	PageSize     int                  `json:"page_size"`
	StudentsList []StudentsListEntity `json:"students_list"`
}

type StudentsListEntity struct {
	StuId    uint   `json:"stu_id"`
	StuName  string `json:"stu_name"`
	StuClass string `json:"stu_class"`
	Grade    string `json:"grade"`
}

type DelStudentForm struct {
	StuID string `json:"stu_id"`
}

type LoginForm struct {
	StuID    string `json:"stu_id"`
	Password string `json:"password"`
}

type LoginRsp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SelectableClasses struct {
	Info []Course `json:"selectable_classes"`
}

type SelectedClasses struct {
	Info []Course `json:"selected_classes"`
}

type CourseForm struct {
	ClassID string `json:"class_id"`
}

type UpdateCourseData struct {
	StuId         string `json:"stu_id"`
	UpdateClassId string `json:"update_class_id"`
}

type UpdateDataForCourse struct {
	ClassID       string                         `json:"class_id"`
	UpdateColumns []UpdateColumnsEntityForCourse `json:"update_columns"`
}

type UpdateColumnsEntityForCourse struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

type UpdateDataForCourseStock struct {
	ClassID string `json:"class_id"`
	Stock   uint   `json:"stock"`
}

type DeleteDataForCourse struct {
	ClassID string `json:"class_id"`
}
