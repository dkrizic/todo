package backend

import (
	"context"
	"github.com/dkrizic/todo/server/backend/repository"
	"testing"
)

// test function convertTodoStructToJson
func Test_convertTodoStructToJson(t *testing.T) {
	todo := &repository.Todo{
		Id:          "1",
		Title:       "test",
		Description: "test",
	}
	todos := []*repository.Todo{todo}
	ctx := context.Background()
	data, err := convertTodoStructToJson(ctx, todos)
	if err != nil {
		t.Error("Error converting todo struct to json")
	}
	if string(data) != "[{\"id\":\"1\",\"title\":\"test\",\"description\":\"test\"}]" {
		t.Error("Error converting todo struct to json")
	}
}
