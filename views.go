package cloudant

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/google/go-querystring/query"
)

// View struct which will allow us to get results from a view
type View struct {
	ddoc   string
	name   string
	Params *ViewParams
}

// ViewParams is the struct which will strucutre the query string for the view find
type ViewParams struct {
	IncludeDocs    bool     `url:"include_docs"`
	Limit          int      `url:"limit"`
	Reduce         bool     `url:"reduce"`
	GroupLevel     int      `url:"group_level"`
	StartKeyString []string `url:"start_key"`
	StartKeyInt    []int64  `url:"start_key"`
	EndKeyString   []string `url:"end_key"`
	EndKeyInt      []int64  `url:"end_key"`
}

// NewView returns a new View object
func NewView(ddoc, name string) *View {
	return &View{
		ddoc:   ddoc,
		name:   name,
		Params: &ViewParams{},
	}
}

// GetViewResults gets a views results
func (db *Database) GetViewResults(v *View, results interface{}) error {
	qp, _ := query.Values(v.Params)
	endpoint := fmt.Sprintf("/%s/_design/%s/_view/%s?%s&group=true", db.Name(), v.ddoc, v.name, qp.Encode())
	resp, err := db.client.doRequest("GET", endpoint, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return newCloudantError(resp)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)

	//fmt.Printf(buf.String()) // NOTE: uncomment this line when debugging results

	// for test coverage purposes, performing nil check using non-idiomatic pattern
	var objmap map[string]*json.RawMessage
	if err == nil {
		err = json.Unmarshal(buf.Bytes(), &objmap)
	}

	if err == nil {
		err = json.Unmarshal(*objmap["rows"], results)
	}

	return err
}
