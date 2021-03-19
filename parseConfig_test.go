package aerospikeConnect

import (
	"encoding/json"
	"fmt"
	"testing"
)

type configTest struct {
	Aerospike *AerospikeConfig `yaml:"aerospike"`
}

func TestParseConfigSuccess(t *testing.T) {
	ac := new(configTest)
	if err := ParseConfig("config.yml", ac); err != nil {
		t.FailNow()
	}
	data, err := json.MarshalIndent(ac, "", "    ")
	if err != nil {
		t.FailNow()
	}
	fmt.Println(string(data))
}
