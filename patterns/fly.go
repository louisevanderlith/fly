package patterns

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"sync"
)

type ProgramType = int

const (
	Cmd int = iota
	Pkg
	ConfName string = "./fly.json"
)

type (
	//Fly is the main configuration object
	Fly struct {
		Env      environment `json:"environment"`
		Programs []Program   `json:"programs"`
	}

	environment struct {
		Bin  string `json:"bin"`
		Mode string `json:"mode"`
	}

	Program struct {
		Type     ProgramType `json:"type"`
		Name     string      `json:"name"`
		Play     bool        `json:"play"`
		Priority int         `json:"priority"`
		Path     string      `json:"path"`
	}

	StructureInfo struct {
		Name       string
		Path       string
		HasMain    bool
		HasGoFiles bool
	}
)

func DetectConfig(path, mode string) (Fly, error) {
	_, err := os.Stat(ConfName)

	if err != nil {
		//you shouldn't be generating fly configs any other place than DEV.
		conf, err := generateConfig(path, mode)

		if err == nil && mode != "TEST" {
			writeConfig(conf)
		}

		return conf, nil
	}

	return loadConfig()
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

func writeConfig(conf Fly) {
	bits, err := json.Marshal(conf)

	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(ConfName, bits, 0644)
}

func (f *Fly) Build() {
	wg := &sync.WaitGroup{}

	for _, prog := range f.Programs {
		if _, err := os.Stat(prog.Path); err != nil {
			log.Printf("Directory Error: %+v\n", err)
			log.Println(err)
			continue
		}

		wg.Add(1)
		runBuildWg(wg, prog.Path)
	}

	wg.Wait()
}

func (f *Fly) Play(swagger bool) {
	//TODO:build only if application has changed...LATER~
	wg := &sync.WaitGroup{}

	for _, prog := range f.Programs {
		if !prog.Play {
			continue
		}

		//sanity check
		if _, err := os.Stat(prog.Path); err != nil {
			log.Printf("Directory Error: %+v\n", err)
			continue
		}

		buildRes := make(chan string)
		go runBuild(prog.Path, buildRes)

		if swagger && prog.Type == Cmd {
			swaggerDone := make(chan bool)
			go updateSwagger(prog.Path, swaggerDone)
			<-swaggerDone
		}

		log.Println(<-buildRes)
		wg.Add(1)
		go runPlayWg(wg, prog.Path, prog.Name, false)
	}

	wg.Wait()
}

func (f *Fly) Deploy() {
	log.Print("Not running yet. Needs to do what build.ps1 did, also refer to gbuild")
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

func updateSwagger(progDir string, done chan bool) {
	genCmd := exec.Command("bee", "generate", "docs")
	genCmd.Dir = progDir

	err := genCmd.Run()

	if err != nil {
		log.Println(err)
	}

	done <- true
}

func runBuildWg(wg *sync.WaitGroup, progDir string) {
	res := make(chan string)
	go runBuild(progDir, res)

	log.Printf("Build Result %v\n", <-res)
	wg.Done()
}

func runBuild(progDir string, buildRes chan string) {
	cmnd := exec.Command("go", "build")
	cmnd.Dir = progDir
	out, err := cmnd.CombinedOutput()

	if err != nil {
		fmt.Printf("Build Error: %s\n", out)
		buildRes <- "error"
		return
	}

	buildRes <- "complete"
}

func runPlayWg(wg *sync.WaitGroup, progDir, progName string, build bool) {
	if build {
		wg.Add(1)
		runBuildWg(wg, progDir)
	}

	ply := make(chan string, 1)
	go runPlay(progDir, progName, ply)

	log.Printf("%s build %v\n", progName, <-ply)
	wg.Done()
}

func runPlay(progDir, progName string, res chan string) {
	cmnd := exec.Command("./" + progName)
	cmnd.Dir = progDir
	cmnd.Stdout = os.Stdout
	cmnd.Stderr = os.Stderr

	err := cmnd.Start()
	if err != nil {
		fmt.Printf("Play Error: %s Stack:\n", err.Error())
		res <- "error"
		return
	}

	cmnd.Wait()

	res <- progName + " running"
}
