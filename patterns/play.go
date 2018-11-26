package patterns

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

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
		go runBuild(f.Env, prog, buildRes)

		if swagger {
			swaggerDone := make(chan bool)
			go updateSwagger(prog.Path, swaggerDone)
			<-swaggerDone
		}

		log.Println(<-buildRes)
		wg.Add(1)
		go runPlayWg(wg, f.Env, prog, false)
	}

	wg.Wait()
}

func runPlayWg(wg *sync.WaitGroup, env environment, p Program, build bool) {
	if build {
		wg.Add(1)
		runBuildWg(wg, env, p)
	}

	ply := make(chan string, 1)
	go runPlay(p.Path, p.Name, ply)

	log.Printf("%s build %v\n", p.Name, <-ply)
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
