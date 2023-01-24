package notification

import (
	repository "github.com/dkrizic/todo/server/backend/repository"
	"testing"
)

// test convert function
func TestConvert(t *testing.T) {
	// create dummy change object
	change := repository.Change{
		Before: &repository.Todo{
			Id:          "1",
			Title:       "title",
			Description: "description",
		},
		After: &repository.Todo{
			Id:          "2",
			Title:       "title",
			Description: "description",
		},
		ChangeType: "UPDATE",
	}

	// convert todo object to json
	data, err := convert(change)
	if err != nil {
		t.Errorf("Error converting change object to json: %v", err)
	}

	// convert data to string
	str := string(data)

	// compare data with expected value
	expected := "{\"Before\":{\"Id\":\"1\",\"Title\":\"title\",\"Description\":\"description\",\"Status\":\"\"},\"After\":{\"Id\":\"2\",\"Title\":\"title\",\"Description\":\"description\",\"Status\":\"\"},\"ChangeType\":\"UPDATE\"}"
	if str != expected {
		t.Errorf("Expected %v, got %v", expected, str)
	}
}
