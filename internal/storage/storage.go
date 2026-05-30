package storage

import "github.com/siddharth7actowiz/student_api/internal/types"

// usng interfaces we can make ours apps as pluglins

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	
	
	GetAllStudents() ([]types.Student, error)
	
	
	GetStudentByID(id int64) (types.Student, error)

	UpdateStudent(id int64,name string, email string, age int) (types.Student, error)

	DeleteStudent(id int64) (types.Student, error)
	
}
