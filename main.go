package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
)

type (
	fly struct {
		Env      environment `json:"environment"`
		Programs []program   `json:"programs"`
		Dir      string      `json;"_"`
	}

	environment struct {
		Bin  string `json:"bin"`
		Mode string `json:"mode"`
	}

	program struct {
		Type string `json"type"`
		Name string `json:"name"`
		Play bool   `json"play"`
	}
)

func main() {
	modePtr := flag.String("mode", "B", "[P]lay || [B]uild (default) || [D]eploy (lab)")
	swaggerPtr := flag.Bool("swagger", false, "Updates Swagger docs and routers.")
	flag.Parse()

	log.Printf("/*FLYing Mode:%s Swagger:%t*/\n", *modePtr, *swaggerPtr)
	conf, err := loadConfig()

	if err != nil {
		panic(err)
	}

	switch *modePtr {
	case "B":
		build(conf)
	case "P":
		play(conf, *swaggerPtr)
	case "D":
		deploy(conf)
	}

	log.Print("Thank you, FLY again soon.")
}

func loadConfig() (fly, error) {
	result := fly{}

	bits, err := ioutil.ReadFile("./fly.json")

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(bits, &result)

	if err != nil {
		return result, err
	}

	wd, err := os.Getwd()

	if err != nil {
		return result, err
	}

	result.Dir = wd

	return result, err
}

func build(conf fly) {
	wg := &sync.WaitGroup{}

	for _, prog := range conf.Programs {
		if !prog.Play {
			continue
		}

		progDir := fmt.Sprintf("%s/%s/%s", conf.Dir, prog.Type, prog.Name)

		if _, err := os.Stat(progDir); err != nil {
			log.Println(err)
			continue
		}

		wg.Add(1)
		runBuildWg(wg, progDir)
	}

	wg.Wait()
}

func play(conf fly, swagger bool) {
	//build only if application has changed...LATER~
	wg := &sync.WaitGroup{}

	for _, prog := range conf.Programs {
		if !prog.Play {
			continue
		}

		progDir := fmt.Sprintf("%s/%s/%s", conf.Dir, prog.Type, prog.Name)

		if _, err := os.Stat(progDir); err != nil {
			log.Println(err)
			continue
		}

		buildRes := make(chan string)
		go runBuild(progDir, buildRes)

		if swagger && prog.Type == "api" {
			swaggerDone := make(chan bool)
			go updateSwagger(progDir, swaggerDone)
			<-swaggerDone
		}

		wg.Add(1)
		go runPlay(wg, progDir, prog.Name, buildRes)
	}

	wg.Wait()
}

func deploy(conf fly) {
	log.Print("Not running yet. Needs to do what build.ps1 did")
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

	log.Printf("Result %v\n", <-res)
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

func runPlay(wg *sync.WaitGroup, progDir, progName string, buildRes chan string) {
	//cmnd := os.Exec
	res := <-buildRes

	if res == "error" {
		return
	}

	cmnd := exec.Command(progName)
	cmnd.Dir = progDir

	//run in new window...
	err := cmnd.Run()

	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Running %s PID: %v\n", progName, cmnd.Process.Pid)
	}

	wg.Done()
}
