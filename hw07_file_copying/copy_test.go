package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

type testcase struct {
	from, to                    string
	limit, offset, expectedsize int64
	expectedbyte                []byte
	error                       error
}

var testcases []testcase

func TestCopy(t *testing.T) {
	tmpdir := path.Join(os.TempDir(),"testcopytmp")
	defer os.RemoveAll(tmpdir)
	from = path.Join("testdata", "input.txt")
	frominfo, err := os.Lstat(from)
	if err != nil {
		log.Fatal("unable to stat from test file")
	}
	fromsize := frominfo.Size()
	offsetregx, err := regexp.Compile(`offset(\d+)[_\.]`)
	if err != nil {
		log.Fatal("unable to compile regex fo offset")
	}
	limitregx, err := regexp.Compile(`limit(\d+)[_\.]`)
	if err != nil {
		log.Fatal("unable to compile regex fo limit")
	}

	testcases = append(testcases, testcase{
		from:         "bad_from.txt",
		to:           newTempFile(tmpdir, ""),
		limit:        0,
		offset:       0,
		expectedsize: 0,
		expectedbyte: []byte{},
		error:        ErrUnsupportedFile})
	testcases = append(testcases, testcase{
		from:         from,
		to:           newTempFile(tmpdir, ""),
		limit:        0,
		offset:       fromsize + 1,
		expectedsize: 0,
		expectedbyte: []byte{},
		error:        ErrOffsetExceedsFileSize})
	filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.Index(info.Name(), "out_") == 0 {
			setOffset(offsetregx, info)
			setLimit(limitregx, info)
			expectedbyte, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalf("unable to read test out file in tests. Error: %v", err)
			}
			testcases = append(testcases, testcase{
				from:         from,
				to:           newTempFile(tmpdir, info.Name()),
				expectedbyte: expectedbyte,
				limit:        limit,
				offset:       offset,
				expectedsize: info.Size(),
				error:        nil,
			})

		}
		return nil
	})
	for _, testcase := range testcases {
		err := Copy(testcase.from, testcase.to, testcase.offset, testcase.limit)
		if testcase.error != nil {
			require.Equal(t, testcase.error, err, "expected and actual errors are not equal")
			continue
		}
		require.NoErrorf(t, err, "function Copy returned error")
		f, err := os.Open(testcase.to)
		defer f.Close()
		require.NoError(t, err)
		s, _ := f.Stat()
		outf, err := ioutil.ReadAll(f)
		require.NoError(t, err)
		require.Equalf(t, outf, testcase.expectedbyte, "content for file %v with offset %v limit %v", testcase.from, testcase.offset, testcase.limit)
		require.Equalf(t, testcase.expectedsize, s.Size(), "size for file %v with offset %v limit %v", testcase.from, testcase.offset, testcase.limit)
	}

}
func newTempFile(dir, templ string) string {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatalf("unable to create temp dir fo test. Error: %v", err)
	}
	tempFile, err := ioutil.TempFile(dir, "*_"+templ)
	if err != nil {
		log.Fatalf("unable to create temp file fo test. Error: %v", err)
	}
	defer tempFile.Close()
	return tempFile.Name()
}

func setLimit(limitr *regexp.Regexp, info os.FileInfo) {
	lm := limitr.FindStringSubmatch(info.Name())
	if len(lm) > 1 {
		l, err := strconv.Atoi(lm[1])
		if err != nil {
			log.Fatalf("unable to convert limit to int. Error: %v", err)
		}
		limit = int64(l)
	}
}

func setOffset(offsetr *regexp.Regexp, info os.FileInfo) {
	om := offsetr.FindStringSubmatch(info.Name())
	if len(om) > 1 {
		o, err := strconv.Atoi(om[1])
		if err != nil {
			log.Fatalf("unable to convert offset to int. Error: %v", err)
		}
		offset = int64(o)
	}
}
