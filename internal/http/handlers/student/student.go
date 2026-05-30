package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/siddharth7actowiz/student_api/internal/storage"
	"github.com/siddharth7actowiz/student_api/internal/types"
	"github.com/siddharth7actowiz/student_api/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc { //dependenccy injection
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a New Student")

		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body ")))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//validation shoud be done
		//Request validation
		//Run validator.New().Struct(student) first, then extract the error — not before.
		//Use errors.As() instead of a direct type assertion .(validator.ValidationErrors) — a direct assertion panics if the error isn't that type, which is what you were seeing.
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}
		//creating a stdent
		lastid, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age)
		slog.Info("Student Created Succesull", slog.String("student_id", fmt.Sprintf("%d", lastid)))
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastid})

	}
}

// Function for get student by id
func GetByID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id") //this is inbuild method of golang request object it will return the value of this parmas written in endpoint
		//params name should be same as you passed in endpoint
		s_id, e := strconv.ParseInt(id, 10, 64)
		if e != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(e))
			return
		}
		slog.Info("Gettig a student", slog.String("student_id", id))
		student, err := storage.GetStudentByID(s_id)

		if err != nil {
			slog.Error("Error getting user", slog.String("ID:->", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}

}

func GetAll(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Gettig all students record")
		students, err := storage.GetAllStudents()
		if err != nil {
			slog.Error("Error getting all students")
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, students)
	}
}

//update 

func Update(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e_id:=r.PathValue("id")
		slog.Info("Updating  student with id: ",e_id)
		id, err := strconv.ParseInt(e_id, 10, 64)
			if err != nil {
				response.WriteJson(
					w,
					http.StatusBadRequest,
					response.GeneralError(err),
				)
				return
			}
		//	
		var student types.Student	
		err=json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body ")))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, err = storage.UpdateStudent(id,student.Name,student.Email,	student.Age,)
		if err!=nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))


		}

		response.WriteJson(w,http.StatusAccepted,student)


	}
}

//delete 
func Delete(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		e_id:=r.PathValue("id")
		slog.Info("Deletiny student with id: ",e_id)
		id,err:=strconv.ParseInt(e_id,10,64)
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
		}
		student,err:=storage.DeleteStudent(id)
		if err!=nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
		}

		response.WriteJson(w,http.StatusAccepted,student)
	}
}
