package mysql

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
	"github.com/siddharth7actowiz/student_api/internal/config"
	"github.com/siddharth7actowiz/student_api/internal/types"
)

type MySQLStorage struct {
	db *sql.DB
}
//This is dependency Injection
func New(cfg *config.Config) (*MySQLStorage, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db connection failed: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
        id    INTEGER PRIMARY KEY AUTO_INCREMENT,
        name  TEXT,
        email TEXT,
        age   INTEGER
    )`)
	if err != nil {
		return nil, err
	}

	return &MySQLStorage{db: db}, nil
}

//Function for Create Table
func (m *MySQLStorage) CreateStudent(name string, email string, age int) (int64, error) {
	query, err := m.db.Prepare("INSERT INTO STUDENTS (name,email,age)VALUES(?,?,?)")
	if err != nil {
		return 0, err
	}
	defer query.Close()
	result, err := query.Exec(name, email, age)

	if err != nil {
		return 0, err
	}

	last_Id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return last_Id, nil

}
//Function For get students by id
func (m *MySQLStorage) GetStudentByID(id int64) (types.Student, error) {
	query, err := m.db.Prepare("SELECT * FROM STUDENTS WHERE id=? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer query.Close()

	var student types.Student

	err = query.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("No Student with id :%s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("Query Error:", query, err)
	}
	return student, nil
}

//All students
func (m *MySQLStorage) GetAllStudents() ([]types.Student, error) {
	rows, err := m.db.Query("SELECT * FROM STUDENTS")
	if err != nil {

		return []types.Student{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var students_arr []types.Student


	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return []types.Student{}, fmt.Errorf("failed to scan student row: %w", err)
		}
		students_arr = append(students_arr, student)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return students_arr, nil
}


//update student

func(m *MySQLStorage) UpdateStudent(id int64,name string, email string, age int) (types.Student, error){
	var Student types.Student
	//update query
	_,err:=m.db.Exec(`
		UPDATE students
		set name=? ,email=?,age=?
		where id=?
	`,name,
	email,
	age,id,)
	if err !=nil{
		return types.Student{},err
	}


	//select query to show wupdated row
	err = m.db.QueryRow(`
	SELECT id, name, email, age
	FROM students
	WHERE id=?
`, id).Scan(
	&Student.Id,
	&Student.Name,
	&Student.Email,
	&Student.Age,
)
	if err !=nil{
		return types.Student{},err
	}


	return Student, nil
}


func (m *MySQLStorage) DeleteStudent(id int64) (types.Student, error) {

	slog.Warn("Deleting Student", slog.Int64("id", id))

	var student types.Student

	err := m.db.QueryRow(`
		SELECT id, name, email, age
		FROM students
		WHERE id=?
	`, id).Scan(
		&student.Id,
		&student.Name,
		&student.Email,
		&student.Age,
	)

	if err != nil {
		return types.Student{}, err
	}

	_, err = m.db.Exec(
		"DELETE FROM students WHERE id=?",
		id,
	)

	if err != nil {
		return types.Student{}, fmt.Errorf(
			"error deleting student: %w",
			err,
		)
	}

	return student, nil
}