package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"homework-5/internal/app/group"
	"homework-5/internal/app/student"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	Port = ":80"
)

type Server interface {
	Get(r *http.Request) ([]byte, int)
	Post(r *http.Request) ([]byte, int)
	Put(r *http.Request) ([]byte, int)
	Delete(r *http.Request) ([]byte, int)
	Handle(w http.ResponseWriter, r *http.Request)
}

type GroupsRepository interface {
	Add(ctx context.Context, group *group.Group) (uint64, error)                 //create
	GetById(ctx context.Context, id uint64) (*group.Group, error)                //read
	UpdateById(ctx context.Context, id uint64, group *group.Group) (bool, error) // update
	Remove(ctx context.Context, id uint64) (bool, error)                         //delete
}

type StudentsRepository interface {
	Add(ctx context.Context, student *student.Student) (uint64, error)                 //create
	GetById(ctx context.Context, id uint64) (*student.Student, error)                  //read
	UpdateById(ctx context.Context, id uint64, student *student.Student) (bool, error) // update
	Remove(ctx context.Context, id uint64) (bool, error)                               //delete
}

type server struct {
	studentRepo StudentsRepository
	groupRepo   GroupsRepository
	ctx         context.Context
}

func New(ctx context.Context, studentRepo *student.StudentsRepository, groupRepo *group.GroupsRepository) Server {
	return &server{
		ctx:         ctx,
		studentRepo: studentRepo,
		groupRepo:   groupRepo}
}

func (s *server) GetIdQuery(r *http.Request) (uint64, error) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		return 0, errors.New("no id")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *server) GetStudentFromBody(r *http.Request) (student.Student, error) {
	data, err := io.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error while reading body, method:[server/Post]")
		return student.Student{}, err
	}

	var unmarshalled student.Student

	if err = json.Unmarshal(data, &unmarshalled); err != nil {
		fmt.Printf("Error while unmarshalling body, method:[server/Post]")
		return student.Student{}, err
	}

	return unmarshalled, nil
}

func (s *server) GetGroupFromBody(r *http.Request) (group.Group, error) {
	data, err := io.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error while reading body, method:[server/Post]")
		return group.Group{}, err
	}

	var unmarshalled group.Group

	if err = json.Unmarshal(data, &unmarshalled); err != nil {
		fmt.Printf("Error while unmarshalling body, method:[server/Post]")
		return group.Group{}, err
	}

	return unmarshalled, nil
}

func (s *server) Get(r *http.Request) ([]byte, int) {
	id, err := s.GetIdQuery(r)
	if err != nil {
		return nil, http.StatusBadRequest
	}
	table := r.URL.Query().Get("table")

	var object interface{}

	switch table {
	case "student":
		object, err = s.studentRepo.GetById(s.ctx, id)
	case "group":
		object, err = s.groupRepo.GetById(s.ctx, id)
	case "":
		return nil, http.StatusBadRequest
	}

	if err != nil {
		return nil, http.StatusNotFound
	}

	marshalled, err := json.Marshal(object)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return marshalled, http.StatusOK
}

func (s *server) Post(r *http.Request) ([]byte, int) {
	table := r.URL.Query().Get("table")

	var id uint64
	var err error

	switch table {
	case "student":
		student, err := s.GetStudentFromBody(r)
		if err != nil {
			return nil, http.StatusBadRequest
		}
		id, err = s.studentRepo.Add(s.ctx, &student)
	case "group":
		group, err := s.GetGroupFromBody(r)
		if err != nil {
			return nil, http.StatusBadRequest
		}
		id, err = s.groupRepo.Add(s.ctx, &group)
	default:
		return nil, http.StatusBadRequest
	}

	if err != nil {
		fmt.Println(err)
		return nil, http.StatusInternalServerError
	}
	return []byte(fmt.Sprintf("%v with id[%v] has been added to database", table, id)), http.StatusOK
}

func (s *server) Put(r *http.Request) ([]byte, int) {
	id, err := s.GetIdQuery(r)
	if err != nil {
		return nil, http.StatusBadRequest
	}

	table := r.URL.Query().Get("table")

	var ok bool

	switch table {
	case "student":
		student, err := s.GetStudentFromBody(r)
		if err != nil {
			return nil, http.StatusBadRequest
		}
		ok, err = s.studentRepo.UpdateById(s.ctx, id, &student)
	case "group":
		group, err := s.GetGroupFromBody(r)
		if err != nil {
			return nil, http.StatusBadRequest
		}
		ok, err = s.groupRepo.UpdateById(s.ctx, id, &group)
	default:
		return nil, http.StatusBadRequest
	}

	if err != nil {
		return nil, http.StatusInternalServerError
	}

	if ok == false {
		return nil, http.StatusBadRequest
	}

	return []byte("Object has been successfully updated"), http.StatusOK
}

func (s *server) Delete(r *http.Request) ([]byte, int) {
	id, err := s.GetIdQuery(r)
	if err != nil {
		return nil, http.StatusBadRequest
	}
	table := r.URL.Query().Get("table")
	var ok bool

	switch table {
	case "student":
		ok, err = s.studentRepo.Remove(s.ctx, id)
	case "group":
		ok, err = s.groupRepo.Remove(s.ctx, id)
	default:
		return nil, http.StatusBadRequest
	}

	if err != nil {
		return nil, http.StatusInternalServerError
	}

	if ok == false {
		return nil, http.StatusNotFound
	}

	return []byte("Object has been successfully deleted"), http.StatusOK
}

func (s *server) Handle(w http.ResponseWriter, r *http.Request) {
	var (
		buf  []byte
		code int
	)
	switch r.Method {
	case http.MethodGet:
		buf, code = s.Get(r)
	case http.MethodPost:
		buf, code = s.Post(r)
	case http.MethodPut:
		buf, code = s.Put(r)
	case http.MethodDelete:
		buf, code = s.Delete(r)
	default:
		code = http.StatusBadRequest
	}
	w.WriteHeader(code)

	if code == http.StatusOK {
		w.Write(buf)
	}
}
