package notification

import (
	"github.com/dkrizic/todo/api/todo"
	"testing"
)

// test convert function
func TestConvert(t *testing.T) {
	// create dummy change object
	change := todo.Change{
		Before: &todo.ToDo{
			Id:          "1",
			Title:       "title",
			Description: "description",
		},
		After: &todo.ToDo{
			Id:          "2",
			Title:       "title",
			Description: "description",
		},
	}

	// convert todo object to json
	data, err := convert(change)
	if err != nil {
		t.Errorf("Error converting change object to json: %v", err)
	}

	// convert data to string
	str := string(data)

	// compare data with expected value
	expected := "{\"before\":{\"id\":\"1\",\"title\":\"title\",\"description\":\"description\"},\"after\":{\"id\":\"2\",\"title\":\"title\",\"description\":\"description\"}}"
	if str != expected {
		t.Errorf("Expected %v, got %v", expected, str)
	}
}
