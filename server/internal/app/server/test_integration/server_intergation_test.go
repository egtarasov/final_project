package test_integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	group2 "homework-5/server/internal/app/group"
	student2 "homework-5/server/internal/app/student"
	"net/http"
	"net/url"
	"testing"
)

func Test_Get(t *testing.T) {
	t.Run("success_student", func(t *testing.T) {
		t.Parallel()
		var (
			st = student2.DefaultStudent().P()
			gr = group2.DefaultGroup().P()
		)
		s := NewTestHandler()
		s.setUp(st, gr)
		defer s.tearDown()

		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("localhost?id=%v&table=student", st.Id), bytes.NewReader([]byte{}))

		marshalled, statusCode := s.server.Get(request)

		require.Equal(t, http.StatusOK, statusCode)
		var actual student2.Student
		_ = json.Unmarshal(marshalled, &actual)
		assert.ObjectsAreEqualValues(st, actual)
	})
	t.Run("success_group", func(t *testing.T) {
		t.Parallel()
		var (
			st = student2.DefaultStudent().P()
			gr = group2.DefaultGroup().P()
		)
		s := NewTestHandler()
		s.setUp(st, gr)
		defer s.tearDown()
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("localhost?id=%v&table=group", gr.Id), bytes.NewReader([]byte{}))

		marshalled, statusCode := s.server.Get(request)

		require.Equal(t, http.StatusOK, statusCode)

		var actual group2.Group
		_ = json.Unmarshal(marshalled, &actual)

		assert.ObjectsAreEqualValues(gr, actual)
	})
}

func TestServer_Post(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var (
			st = student2.DefaultStudent().P()
			gr = group2.DefaultGroup().P()
		)
		s := NewTestHandler()
		s.setUp(st, gr)

		studentBuff, _ := json.Marshal(st)
		groupBuff, _ := json.Marshal(gr)
		requestStudent, _ := http.NewRequest(http.MethodPost, "localhost?table=student", bytes.NewReader(studentBuff))
		requestGroup, _ := http.NewRequest(http.MethodPost, "localhost?table=group", bytes.NewReader(groupBuff))

		tt := []struct {
			request    *http.Request
			codeExpect int
		}{
			{request: requestGroup, codeExpect: http.StatusOK},
			{request: requestStudent, codeExpect: http.StatusOK},
			{request: &http.Request{Method: http.MethodPost, URL: &url.URL{RawQuery: "table=sdaf"}}, codeExpect: http.StatusBadRequest},
			{request: &http.Request{Method: http.MethodPost, URL: &url.URL{RawQuery: "id=123"}}, codeExpect: http.StatusBadRequest},
		}

		for _, tc := range tt {
			_, code := s.server.Post(tc.request)
			assert.Equal(t, tc.codeExpect, code)
		}
	})
}

func TestServer_Put(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var (
			st = student2.DefaultStudent().P()
			gr = group2.DefaultGroup().P()
		)
		s := NewTestHandler()
		s.setUp(st, gr)

		// Change structs
		st.FirstName = "Bob"
		gr.Name = sql.NullString{String: "Bob_group", Valid: true}

		studentBuff, _ := json.Marshal(st)
		groupBuff, _ := json.Marshal(gr)
		requestStudent, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("localhost?table=student&id=%v", st.Id), bytes.NewReader(studentBuff))
		requestGroup, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("localhost?table=group&id=%v", gr.Id), bytes.NewReader(groupBuff))

		tt := []struct {
			request      *http.Request
			resultExpect []byte
			codeExpect   int
		}{
			{request: requestGroup, resultExpect: []byte("Object has been successfully updated"), codeExpect: http.StatusOK},
			{request: requestStudent, resultExpect: []byte("Object has been successfully updated"), codeExpect: http.StatusOK},
		}

		for _, tc := range tt {
			result, code := s.server.Put(tc.request)
			assert.Equal(t, tc.codeExpect, code)
			assert.Equal(t, tc.resultExpect, result)
		}
	})
}

func TestServer_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var (
			st = student2.DefaultStudent().P()
			gr = group2.DefaultGroup().P()
		)
		s := NewTestHandler()
		s.setUp(st, gr)

		tt := []struct {
			request      *http.Request
			resultExpect []byte
			codeExpect   int
		}{
			{request: &http.Request{Method: http.MethodDelete,
				URL: &url.URL{RawQuery: fmt.Sprintf("id=%v&table=student", st.Id)}},
				resultExpect: []byte("Object has been successfully deleted"), codeExpect: http.StatusOK},
			{request: &http.Request{Method: http.MethodDelete,
				URL: &url.URL{RawQuery: fmt.Sprintf("id=%v&table=group", gr.Id)}},
				resultExpect: []byte("Object has been successfully deleted"), codeExpect: http.StatusOK},
		}
		for _, tc := range tt {
			result, code := s.server.Delete(tc.request)
			assert.Equal(t, tc.codeExpect, code)
			assert.Equal(t, tc.resultExpect, result)
		}
	})
}
