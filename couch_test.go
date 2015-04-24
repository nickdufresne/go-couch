package couch

import (
	"testing"
)

func should(t *testing.T) func(*Response, error) {
	return func(r *Response, err error) {
		if err != nil {
			t.Error(err)
		}
	}
}

func TestPutDB(t *testing.T) {
	check := should(t)
	check(Delete("database_test"))

	create, err := Put("database_test", nil)
	if err != nil {
		t.Error(err)
	}

	if create.Status != 201 {
		t.Errorf("State should be 201.  Got: %d", create.Status)
	}
}

func TestPostToDB(t *testing.T) {
	check := should(t)
	check(Delete("database_test"))
	check(Put("database_test", nil))

	json := map[string]string{"Message": "Hello"}

	resp, err := Post("database_test/", json)

	if err != nil {
		t.Error(err)
	}

	if resp.Status != 201 {
		t.Errorf("State should be 201.  Got: %d", resp.Status)
	}
}

func TestDBExists(t *testing.T) {
	check := should(t)
	check(Delete("database_test"))

	db := DB("database_test")
	ok, err := db.Exists()

	if err != nil {
		t.Error(err)
	}

	if ok {
		t.Error("Expected DB to not exist")
	}

	check(Put("database_test", nil))

	ok, err = db.Exists()
	if err != nil {
		t.Error(err)
	}

	if !ok {
		t.Error("Expected DB to exist")
	}

}
