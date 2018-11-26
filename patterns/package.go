package patterns

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type BuildBall struct {
	Name  string //application bundled by the Ball
	Files []*FileBody
}

type FileBody struct {
	Name string
	Body []byte
}

/*func scanAppBin(path string, p Program) (BuildBall, error) {
	result := BuildBall{
		Name: p.Name,
	}

	_, err := os.Stat(path)

	if err != nil {
		return result, err
	}

	log.Printf("Reading for %s at %s\n", p.Name, path)
	appFiles, err := ioutil.ReadDir(path)

	if err != nil {
		return result, err
	}

	for _, appFile := range appFiles {
		content, err := ioutil.ReadFile(appFile.Name())

		if err != nil {
			log.Println(err)
			continue
		}

		fbody := &FileBody{
			Name: appFile.Name(),
			Body: content,
		}

		result.Files = append(result.Files, fbody)
	}

	return result, nil
}*/

func tarFolder(source, target, mode string) error {
	filename := filepath.Base(source)
	stamp := time.Now().Format("20060102150405")
	target = filepath.Join(target, fmt.Sprintf("%s.%s.%s.tar", filename, mode, stamp))
	tarfile, err := os.Create(target)

	if err != nil {
		return err
	}

	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(source)

	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if err := tarball.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()
		_, err = io.Copy(tarball, file)

		return err
	})
}

/*
func tarFiles(files []*FileBody) (bytes.Buffer, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return buf, err
		}

		if _, err := tw.Write([]byte(file.Body)); err != nil {
			return buf, err
		}
	}

	if err := tw.Close(); err != nil {
		return buf, err
	}

	return buf, nil
}
*/
