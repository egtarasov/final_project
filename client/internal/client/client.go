package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"homework-5/client/internal/group_client"
	"homework-5/client/internal/student_client"
	"homework-5/client/pb/group_repo"
	"homework-5/client/pb/student_repo"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var (
	errNoId    = errors.New("no id")
	errParseId = errors.New("cant parse id")
)

type Client struct {
	ctx           context.Context
	groupClient   group_client.Client
	studentClient student_client.Client
}

func NewClient(ctx context.Context, groupClient group_client.Client, studentClient student_client.Client) *Client {
	return &Client{
		ctx:           ctx,
		groupClient:   groupClient,
		studentClient: studentClient,
	}
}

func (s *Client) GetIdQuery(r *url.URL) (int64, error) {
	idStr := r.Query().Get("id")

	if idStr == "" {
		return 0, errNoId
	}

	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		return 0, errParseId
	}

	return id, nil
}

func (s *Client) GetStudentFromBody(r *http.Request) (student_repo.Student, error) {
	data, err := io.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error while reading body, method:[Server/Post]")
		return student_repo.Student{}, err
	}

	var unmarshalled student_repo.Student

	if err = json.Unmarshal(data, &unmarshalled); err != nil {
		fmt.Printf("Error while unmarshalling body, method:[Server/Post]")
		return student_repo.Student{}, err
	}

	return unmarshalled, nil
}

func (s *Client) GetGroupFromBody(r *http.Request) (group_repo.Group, error) {
	data, err := io.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error while reading body, method:[Server/Post]")
		return group_repo.Group{}, err
	}

	var unmarshalled group_repo.Group

	if err = json.Unmarshal(data, &unmarshalled); err != nil {
		fmt.Printf("Error while unmarshalling body, method:[Server/Post]")
		return group_repo.Group{}, err
	}

	return unmarshalled, nil
}

func (s *Client) Get(r *http.Request) ([]byte, int) {
	tp := otel.Tracer("GetHttpClient")
	ctx, span := tp.Start(s.ctx, "client received request from user")
	defer span.End()
	id, err := s.GetIdQuery(r.URL)
	if err != nil {
		return nil, http.StatusBadRequest
	}
	table := r.URL.Query().Get("table")

	span.SetAttributes(attribute.Key("id").Int64(id))
	span.SetAttributes(attribute.Key("table").String(table))

	var object interface{}

	switch table {
	case "student":
		object, err = s.studentClient.GetById(ctx, id)
	case "group":
		object, err = s.groupClient.GetById(ctx, id)
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

func (s *Client) Post(r *http.Request) ([]byte, int) {
	tp := otel.Tracer("GetHttpClient")
	ctx, span := tp.Start(s.ctx, "client received request from user")
	defer span.End()

	table := r.URL.Query().Get("table")

	span.SetAttributes(attribute.Key("table").String(table))
	var id int64
	var err error

	switch table {
	case "student":
		student, err := s.GetStudentFromBody(r)
		if err != nil {
			return nil, http.StatusBadRequest
		}
		// Add info about student to Span
		span.SetAttributes(attribute.Key("student").String(student.String()))

		id, err = s.studentClient.Create(ctx, &student)
	case "group":
		group, err := s.GetGroupFromBody(r)
		if err != nil {
			return nil, http.StatusBadRequest
		}

		// Add info about group to Span
		span.SetAttributes(attribute.Key("group").String(group.String()))

		id, err = s.groupClient.Create(ctx, &group)
	default:
		return nil, http.StatusBadRequest
	}

	if err != nil {
		fmt.Println(err)
		return nil, http.StatusInternalServerError
	}
	return []byte(fmt.Sprintf("%v with id[%v] has been added to database", table, id)), http.StatusOK
}

func (s *Client) Put(r *http.Request) ([]byte, int) {
	tp := otel.Tracer("GetHttpClient")
	ctx, span := tp.Start(s.ctx, "client received request from user")
	defer span.End()

	id, err := s.GetIdQuery(r.URL)
	if err != nil {
		return nil, http.StatusBadRequest
	}

	table := r.URL.Query().Get("table")
	// Add info about 'table' and 'id' to Span
	span.SetAttributes(attribute.Key("table").String(table))
	span.SetAttributes(attribute.Key("id").Int64(id))

	var ok bool

	switch table {
	case "student":
		student, err := s.GetStudentFromBody(r)
		if err != nil {
			return nil, http.StatusBadRequest
		}
		// Add info about student to Span
		span.SetAttributes(attribute.Key("student").String(student.String()))

		ok, err = s.studentClient.Update(ctx, id, &student)
	case "group":
		group, err := s.GetGroupFromBody(r)
		if err != nil {
			return nil, http.StatusBadRequest
		}
		// Add info about group to Span
		span.SetAttributes(attribute.Key("group").String(group.String()))

		ok, err = s.groupClient.Update(ctx, id, &group)
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

func (s *Client) Delete(r *http.Request) ([]byte, int) {
	tp := otel.Tracer("GetHttpClient")
	ctx, span := tp.Start(s.ctx, "client received request from user")
	defer span.End()

	id, err := s.GetIdQuery(r.URL)
	if err != nil {
		return nil, http.StatusBadRequest
	}
	table := r.URL.Query().Get("table")

	// Add info about 'table' and 'id' to Span
	span.SetAttributes(attribute.Key("table").String(table))
	span.SetAttributes(attribute.Key("id").Int64(id))

	var ok bool

	switch table {
	case "student":
		ok, err = s.studentClient.Delete(ctx, id)
	case "group":
		ok, err = s.groupClient.Delete(ctx, id)
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

func (s *Client) Handle(w http.ResponseWriter, r *http.Request) {
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
