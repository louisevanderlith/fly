package patterns

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func (f *Fly) Build() {
	wg := &sync.WaitGroup{}

	for _, prog := range f.Programs {
		if _, err := os.Stat(prog.Path); err != nil {
			log.Printf("Directory Error: %+v\n", err)
			log.Println(err)
			continue
		}

		wg.Add(1)
		go runBuildWg(wg, f.Env, prog)
	}

	wg.Wait()
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

//func runBuildWg(wg *sync.WaitGroup, progDir, buildDir, progName, mode string) {
func runBuildWg(wg *sync.WaitGroup, env environment, p Program) {
	res := make(chan string)
	go runBuild(env, p, res)

	log.Printf("Build %s Result %v\n", p.Name, <-res)
	wg.Done()
}

//func runBuild(progDir, buildDir, progName, mode string, buildRes chan string) {
func runBuild(env environment, p Program, buildRes chan string) {
	outDir := getOutDir(env, p.Name)
	outFile := filepath.Join(outDir, p.Name)
	cmnd := exec.Command("go", "build", "-o", outFile, "-i")
	cmnd.Dir = p.Path
	out, err := cmnd.CombinedOutput()

	if err != nil {
		fmt.Printf("Build %s Error: %s\n", p.Name, out)
		buildRes <- "error"
		return
	}

	copyAdditionalItems(env, p, outDir)

	buildRes <- "complete"
}

func getOutDir(env environment, progName string) string {
	wd, _ := os.Getwd()

	return fmt.Sprintf("%s/%s/%s/%s", wd, env.Bin, env.Mode, progName)
}

func copyAdditionalItems(env environment, p Program, outDir string) {
	for _, v := range env.Copy {
		wd, _ := os.Getwd()
		src := filepath.Join(wd, p.Path, v)
		dst := filepath.Join(outDir, v)

		err := CopyFile(src, dst)

		if err != nil {
			fmt.Printf("Unable to copy %s; to %s %+v\n", src, dst, err)
		}
	}
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	dir, fle := filepath.Split(dst)

	if len(fle) > 0 {
		os.MkdirAll(dir, os.ModePerm)
	}

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
