package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gr "homework-5/internal/app/group"
	st "homework-5/internal/app/student"
	"net/http"
	"net/url"
	"testing"
)

func TestServer_GetIdQuery(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var (
			ctx = context.Background()
		)
		s := NewServer(ctx, nil, nil)

		tt := []struct {
			url    *url.URL
			expect uint64
		}{
			{&url.URL{RawQuery: "id=213"}, 213},
			{&url.URL{RawQuery: "name=23&id=23024"}, 23024},
			{&url.URL{RawQuery: "id=1&asd=asd"}, 1},
			{&url.URL{RawQuery: "id=0"}, 0},
		}
		for _, tc := range tt {
			id, err := s.GetIdQuery(tc.url)
			require.NoError(t, err)
			assert.Equal(t, tc.expect, id)
		}
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		var (
			ctx = context.Background()
		)
		s := NewServer(ctx, nil, nil)

		tt := []struct {
			url    *url.URL
			expect error
		}{
			{&url.URL{RawQuery: "name=23"}, errNoId},
			{&url.URL{RawQuery: "id=name&asd=asd"}, errParseId},
			{&url.URL{RawQuery: "id=35435489578934858345784845788345745934578458457348597345834758347589375898"}, errParseId},
		}
		for _, tc := range tt {
			_, err := s.GetIdQuery(tc.url)
			assert.Equal(t, tc.expect, err)
		}
	})
}

func TestServer_GetGroupFromBody(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var (
			ctx     = context.Background()
			student = gr.DefaultGroup().V()
		)
		buffer, err := json.Marshal(student)
		s := NewServer(ctx, nil, nil)

		request, _ := http.NewRequest(http.MethodGet, "local?id=21", bytes.NewReader(buffer))

		actual, err := s.GetGroupFromBody(request)
		require.NoError(t, err)
		assert.ObjectsAreEqual(student, actual)
	})
}

func TestServer_GetStudentFromBody(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var (
			ctx     = context.Background()
			student = st.DefaultStudent().V()
		)
		buffer, err := json.Marshal(student)
		s := NewServer(ctx, nil, nil)

		request, _ := http.NewRequest(http.MethodGet, "local?id=21", bytes.NewReader(buffer))

		actual, err := s.GetStudentFromBody(request)
		require.NoError(t, err)
		assert.ObjectsAreEqual(student, actual)
	})
}

func TestServer_Get(t *testing.T) {
	t.Run("success_student", func(t *testing.T) {
		t.Parallel()
		var (
			m       = setUp(t)
			student = st.DefaultStudent().P()
		)
		defer m.teraDown()
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("localhost?id=%v&table=student", student.Id), bytes.NewReader([]byte{}))
		m.studentRepo.EXPECT().GetById(gomock.Any(), uint64(student.Id)).Return(student, nil)

		marshalled, statusCode := m.s.Get(request)

		require.Equal(t, http.StatusOK, statusCode)

		var actual st.Student
		_ = json.Unmarshal(marshalled, &actual)

		assert.ObjectsAreEqualValues(student, actual)
	})
	t.Run("success_group", func(t *testing.T) {
		t.Parallel()
		var (
			m     = setUp(t)
			group = gr.DefaultGroup().P()
		)
		defer m.teraDown()
		m.groupRepo.EXPECT().GetById(gomock.Any(), uint64(group.Id)).Return(group, nil)

		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("localhost?id=%v&table=group", group.Id), bytes.NewReader([]byte{}))
		marshalled, statusCode := m.s.Get(request)

		require.Equal(t, http.StatusOK, statusCode)

		var actual gr.Group
		_ = json.Unmarshal(marshalled, &actual)

		assert.ObjectsAreEqualValues(group, actual)
	})
	t.Run("fail_table", func(t *testing.T) {
		t.Parallel()
		var (
			ctx = context.Background()
			s   = NewServer(ctx, nil, nil)
		)
		tt := []struct {
			r    *http.Request
			code int
		}{
			{&http.Request{URL: &url.URL{RawQuery: "id=12"}}, http.StatusBadRequest},
			{&http.Request{URL: &url.URL{RawQuery: "table=12"}}, http.StatusBadRequest},
			{&http.Request{URL: &url.URL{RawQuery: "table=students"}}, http.StatusBadRequest},
		}
		for _, tc := range tt {
			_, code := s.Get(tc.r)
			assert.Equal(t, tc.code, code)
		}
	})
}

func TestServer_Post(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var (
			student = st.DefaultStudent().V()
			group   = gr.DefaultGroup().V()
		)
		m := setUp(t)
		defer m.teraDown()

		studentBuff, _ := json.Marshal(student)
		groupBuff, _ := json.Marshal(group)
		requestStudent, _ := http.NewRequest(http.MethodPost, "localhost?table=student", bytes.NewReader(studentBuff))
		requestGroup, _ := http.NewRequest(http.MethodPost, "localhost?table=group", bytes.NewReader(groupBuff))

		m.studentRepo.EXPECT().Add(gomock.Any(), &student).Return(uint64(student.Id), nil)

		tt := []struct {
			request      *http.Request
			resultExpect []byte
			codeExpect   int
		}{
			{request: requestGroup, resultExpect: []byte(fmt.Sprintf("group with id[%v] has been added to database", group.Id)), codeExpect: http.StatusOK},
			{request: requestStudent, resultExpect: []byte(fmt.Sprintf("student with id[%v] has been added to database", student.Id)), codeExpect: http.StatusOK},
			{request: &http.Request{Method: http.MethodPost, URL: &url.URL{RawQuery: "table=sdaf"}}, codeExpect: http.StatusBadRequest},
			{request: &http.Request{Method: http.MethodPost, URL: &url.URL{RawQuery: "id=123"}}, codeExpect: http.StatusBadRequest},
		}

		for _, tc := range tt {
			result, code := m.s.Post(tc.request)
			assert.Equal(t, tc.codeExpect, code)
			assert.Equal(t, tc.resultExpect, result)
		}
	})
}

func TestServer_Put(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var (
			student = st.DefaultStudent().V()
			group   = gr.DefaultGroup().V()
		)
		m := setUp(t)
		defer m.teraDown()

		studentBuff, _ := json.Marshal(student)
		groupBuff, _ := json.Marshal(group)
		requestStudent, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("localhost?table=student&id=%v", student.Id), bytes.NewReader(studentBuff))
		requestGroup, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("localhost?table=group&id=%v", group.Id), bytes.NewReader(groupBuff))

		m.studentRepo.EXPECT().UpdateById(gomock.Any(), uint64(student.Id), &student).Return(true, nil)
		m.groupRepo.EXPECT().UpdateById(gomock.Any(), uint64(group.Id), &group).Return(true, nil)

		tt := []struct {
			request      *http.Request
			resultExpect []byte
			codeExpect   int
		}{
			{request: requestGroup, resultExpect: []byte("Object has been successfully updated"), codeExpect: http.StatusOK},
			{request: requestStudent, resultExpect: []byte("Object has been successfully updated"), codeExpect: http.StatusOK},
		}

		for _, tc := range tt {
			result, code := m.s.Put(tc.request)
			assert.Equal(t, tc.codeExpect, code)
			assert.Equal(t, tc.resultExpect, result)
		}
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		var (
			ctx = context.Background()
		)
		s := NewServer(ctx, nil, nil)

		tt := []struct {
			request      *http.Request
			resultExpect []byte
			codeExpect   int
		}{
			{request: &http.Request{Method: http.MethodPut, URL: &url.URL{RawQuery: "table=2443"}}, resultExpect: nil, codeExpect: http.StatusBadRequest},
			{request: &http.Request{Method: http.MethodPut, URL: &url.URL{RawQuery: ""}}, resultExpect: nil, codeExpect: http.StatusBadRequest},
		}

		for _, tc := range tt {
			result, code := s.Put(tc.request)
			assert.Equal(t, tc.codeExpect, code)
			assert.Equal(t, tc.resultExpect, result)
		}
	})
}

func TestServer_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var (
			idGroup   = uint64(12)
			idStudent = uint64(32)
		)
		m := setUp(t)
		defer m.teraDown()

		m.studentRepo.EXPECT().Remove(gomock.Any(), idStudent).Return(true, nil)
		m.groupRepo.EXPECT().Remove(gomock.Any(), idGroup).Return(true, nil)

		tt := []struct {
			request      *http.Request
			resultExpect []byte
			codeExpect   int
		}{
			{request: &http.Request{Method: http.MethodDelete, URL: &url.URL{RawQuery: fmt.Sprintf("id=%v&table=student", idStudent)}}, resultExpect: []byte("Object has been successfully deleted"), codeExpect: http.StatusOK},
			{request: &http.Request{Method: http.MethodDelete, URL: &url.URL{RawQuery: fmt.Sprintf("id=%v&table=group", idGroup)}}, resultExpect: []byte("Object has been successfully deleted"), codeExpect: http.StatusOK},
		}
		for _, tc := range tt {
			result, code := m.s.Delete(tc.request)
			assert.Equal(t, tc.codeExpect, code)
			assert.Equal(t, tc.resultExpect, result)
		}
	})
}
