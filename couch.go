package couch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	URL = "http://localhost:5984"
)

type Response struct {
	OK        bool   `json:"ok"`
	ID        string `json:"id"`
	Rev       string `json:"rev"`
	ErrorType string `json:"error"`
	Reason    string `json:"reason"`
	Status    int
}

type DatabaseInfo struct {
	Exists             bool   `json:"exists"`
	Name               string `json:"db_name"`
	DocCount           int    `json:"doc_count"`
	DocDelCount        int    `json:"doc_del_count"`
	UpdateSeq          int    `json:"update_seq"`
	PurgeSeq           int    `json:"purge_seq"`
	CompactRunning     bool   `json:"compact_running"`
	DiskSize           int    `json:"disk_size"`
	DataSize           int    `json:"data_size"`
	InstanceStartTime  string `json:"instance_start_time"`
	DiskFormatVersion  int    `json:"disk_format_version"`
	CommittedUpdateSeq int    `json:"committed_update_seq"`
}

type database struct {
	url  string
	name string
}

func DB(name string) *database {
	return &database{URL, name}
}

func (db *database) URL() string {
	return fmt.Sprintf("%s/%s", db.url, db.name)
}

func (db *database) Create(v interface{}) (*Response, error) {
	return Post(db.URL()+"/", v)
}

func (db *database) Info() (*DatabaseInfo, error) {
	info := new(DatabaseInfo)

	status, err := sendRequest("GET", db.URL(), nil, info)

	if err != nil {
		return nil, err
	}

	info.Exists = status == 200

	return info, nil
}

func (db *database) Exists() (bool, error) {
	info, err := db.Info()

	if err != nil {
		return false, err
	}

	return info.Exists, nil
}

func buildURL(url string) string {
	return fmt.Sprintf("%s/%s", URL, url)
}

func Post(url string, v interface{}) (*Response, error) {
	cr := new(Response)
	s, err := sendRequest("POST", buildURL(url), v, cr)
	cr.Status = s
	return cr, err
}

func GetJSON(url string, v interface{}) (int, error) {
	return sendRequest("GET", buildURL(url), nil, v)
}

func Get(url string, v interface{}) (*Response, error) {
	cr := new(Response)
	s, err := sendRequest("GET", buildURL(url), nil, v)
	cr.Status = s
	return cr, err
}

func Put(url string, v interface{}) (*Response, error) {
	cr := new(Response)
	s, err := sendRequest("PUT", buildURL(url), v, cr)
	cr.Status = s
	return cr, err
}

func Delete(url string) (*Response, error) {
	dr := new(Response)
	status, err := sendRequest("DELETE", buildURL(url), nil, dr)
	if err != nil {
		return nil, err
	}
	dr.Status = status
	return dr, nil
}

func sendRequest(method string, url string, payload interface{}, r interface{}) (int, error) {

	//fmt.Printf("Send Request to: (%s) %s\n", method, url)

	var buf bytes.Buffer

	if payload != nil {
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(payload); err != nil {
			return 0, err
		}
	}

	req, err := http.NewRequest(method, url, &buf)

	if err != nil {
		return 0, err
	}

	req.Header["Content-Type"] = []string{"application/json"}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if r != nil {
		dec := json.NewDecoder(resp.Body)

		if err := dec.Decode(r); err != nil {
			return 0, err
		}
	}

	return resp.StatusCode, nil
}
