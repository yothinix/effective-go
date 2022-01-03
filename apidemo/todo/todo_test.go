package todo

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestCreateTodoNotAllowSleepTask(t *testing.T) {
	handler := NewTodoHandler(&gorm.DB{})

	w := httptest.NewRecorder()
	payload := bytes.NewBufferString(`{"text": "sleep"}`)
	req, _ := http.NewRequest("POST", "http://0.0.0.0:8080/todos", payload)
	req.Header.Add("TransactionID", "testIDxxx")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.NewTask(c)

	want := `{"error":"not allowed"}`

	if want != w.Body.String() {
		t.Errorf("want %s but get %s\n", want, w.Body.String())
	}
}
