package utility

import (
	"encoding/json"
	"io/ioutil"
	"github.com/logpost/poc-suggestion-algorithm/models"
)

// JobExpected use for destruct job from object
type JobExpected struct {
	Job	models.Job	`json:"job"`
}

// GetterExpected use for parse only field required
type GetterExpected struct { 
	Getter	[]JobExpected	`json:"getter"`
}

// LoadJSON is method for loading JSON file 
func LoadJSON() []JobExpected {
	var data GetterExpected
	
	readFile, _ := ioutil.ReadFile("./google-maps-response-raw.json")
	_ = json.Unmarshal([]byte(readFile), &data)

	saveFile, _ := json.MarshalIndent(data.Getter, "", " ")
	_ = ioutil.WriteFile("google-maps-response-parsed.json", saveFile, 0644)

	return data.Getter
}