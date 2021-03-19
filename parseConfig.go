package aerospikeConnect

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func ParseConfig(filePath string, ac interface{}) error {
	if ac == nil {
		return fmt.Errorf("AerospikeConfig is nil")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open config file: %s", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(ac)
	if err != nil {
		return fmt.Errorf("decoding configuration file: %s", err)
	}
	return nil
}
