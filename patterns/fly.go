package patterns

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
)

const (
	//ConfName is the configuration file's name
	ConfName string = "./fly.json"
)

type (
	//Fly is the main configuration object
	Fly struct {
		Env      environment `json:"environment"`
		Programs []Program   `json:"programs"`
	}

	environment struct {
		Bin  string   `json:"bin"`
		Mode string   `json:"mode"`
		Copy []string `json:"copy"`
	}

	Program struct {
		Name     string `json:"name"`
		Play     bool   `json:"play"`
		Priority int    `json:"priority"`
		Path     string `json:"path"`
	}
)

func DetectConfig(path, mode string) (Fly, error) {
	_, err := os.Stat(ConfName)

	if err != nil {
		//you shouldn't be generating fly configs any other place than DEV.
		conf := generateConfig(path, mode)

		if mode != "TEST" {
			writeConfig(conf)
		}

		return conf, nil
	}

	return loadConfig()
}

func writeConfig(conf Fly) {
	bits, err := json.Marshal(conf)

	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(ConfName, bits, 0644)
}

func loadConfig() (Fly, error) {
	result := Fly{}
	bits, err := ioutil.ReadFile(ConfName)

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(bits, &result)

	if err != nil {
		return result, err
	}

	//priority sorting
	sort.Sort(&result)

	return result, err
}

func (f *Fly) Len() int {
	return len(f.Programs)
}

func (f *Fly) Less(i, j int) bool {
	return f.Programs[i].Priority > f.Programs[j].Priority
}

func (f *Fly) Swap(i, j int) {
	f.Programs[i], f.Programs[j] = f.Programs[j], f.Programs[i]
}
