package student_service

import (
	"database/sql"
	"homework-5/server/internal/app/pb/student_repo"
	"homework-5/server/internal/app/student"
)

func ParseStudent(st *student_repo.Student) *student.Student {
	return &student.Student{
		Id:             st.Id,
		FirstName:      st.FirstName,
		SecondName:     st.SecondName,
		MiddleName:     sql.NullString{String: st.MiddleName, Valid: true},
		Gpa:            st.Gpa,
		AttendanceRate: st.AttendanceRate,
		GroupId:        sql.NullInt64{Int64: st.GroupId, Valid: true},
	}
}

func ParseStudentRequest(st *student.Student) *student_repo.Student {
	return &student_repo.Student{
		Id:             st.Id,
		FirstName:      st.FirstName,
		SecondName:     st.SecondName,
		MiddleName:     st.MiddleName.String,
		Gpa:            st.Gpa,
		AttendanceRate: st.AttendanceRate,
		GroupId:        st.GroupId.Int64,
	}
}
