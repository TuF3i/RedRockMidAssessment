package models

import (
	"time"

	"gorm.io/gorm"
)

//查数据库 + data字段返回

type Relation struct {
	// 关联表
	gorm.Model
	CouID string
	StuID uint
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
	ClassSelectedNum uint   `json:"class_selection" gorm:"class_selected_num; not null; type:int"`

	Students []Student `json:"students" gorm:"many2many:relation; foreignKey:ClassID; joinForeignKey:CouID; references:StudentID; joinReferences:StuID"`
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

	Name         string `json:"name" gorm:"column:name; not null; type:varchar(90)"`
	StudentID    uint   `json:"stu_id" gorm:"column:student_id; index; unique; not null; type:bigint"`
	StudentClass string `json:"stu_class" gorm:"column:student_class; not null; type:varchar(40)"`
	Password     string `json:"password" gorm:"column:password; not null; type:varchar(90)"`
	Sex          uint   `json:"sex" gorm:"column:sex; not null"`
	Grade        uint   `json:"grade" gorm:"column:grade; not null"`
	Age          uint   `json:"age" gorm:"column:age; not null"`

	Courses []Course `json:"courses" gorm:"many2many:relation; foreignKey:StudentID; joinForeignKey:StuID; references:ClassID; joinReferences:CouID"`
}

func (Student) TableName() string {
	return "student"
}

// 更新学生信息column表
type UpdateData struct {
	UpdateColumns []UpdateColumnsEntity `json:"update_columns"`
}

type UpdateColumnsEntity struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}
