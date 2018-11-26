package patterns

import (
	"io/ioutil"
	"testing"
)

func TestPackage_Tar_GeneratesFile(t *testing.T) {
	files := []*FileBody{
		&FileBody{
			Name: "ABC.txt",
			Body: []byte("ABC"),
		},
		&FileBody{
			Name: "OtherFile.data",
			Body: []byte("123 - ABC"),
		},
	}

	buff, err := tarFiles(files)

	if err != nil {
		t.Error(err)
	}

	err = ioutil.WriteFile("testing.tar.gz", buff.Bytes(), 0600)

	if err != nil {
		t.Error(err)
	}
}
