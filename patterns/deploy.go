package patterns

import (
	"log"
	"path/filepath"
	"sync"
)

func (f *Fly) Deploy() {
	wg := &sync.WaitGroup{}

	for _, prog := range f.Programs {
		wg.Add(1)
		go runDeployWg(wg, f.Env.Bin, f.Env.Mode, prog)
	}

	wg.Wait()
}

func runDeployWg(wg *sync.WaitGroup, binPath, mode string, p Program) {
	res := make(chan bool)
	go runDeploy(binPath, mode, p, res)

	log.Printf("Deployed %s - %v\n", mode, <-res)
	wg.Done()
}

func runDeploy(binPath, mode string, p Program, res chan bool) {
	appBin := filepath.Join(binPath, mode, p.Name)
	err := tarFolder(appBin, binPath, mode)

	if err != nil {
		log.Println("Error:", err)
		res <- false
		return
	}

	res <- true
}
