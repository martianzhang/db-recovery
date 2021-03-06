package redo

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var p *Parse
var devPath string
var fixturePath string

// go test -mysql-version 5.6 -v -timeout 30s -run TestFunc
var mysqlRelease = flag.String("mysql-release", "mysql", "mysql docker image release vendor, eg. mysql, percona, mariadb")
var mysqlVersion = flag.String("mysql-version", "5.7", "mysql docker image versions, eg. 5.7, 8.0")

func TestMain(m *testing.M) {
	p = NewParse("", "")
	_, filename, _, _ := runtime.Caller(0)

	flag.Parse()
	devPath = filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	fixturePath = devPath + "/cmd/test/fixture/" + *mysqlRelease + "_" + *mysqlVersion

	// set logs to stderr, and log-level = trace
	flag.Set("logtostderr", "true")
	flag.Set("v", "5")

	m.Run()
}

func TestParseRedoLogs(t *testing.T) {
	flag.Set("v", "3")
	err := p.ParseDictPage(fixturePath + "/ibdata1")
	if err != nil {
		panic(err)
	}

	err = p.ParseRedoLogs([]string{fixturePath + "/ib_logfile0", fixturePath + "/ib_logfile1"})
	if err != nil {
		panic(err)
	}
	flag.Set("v", "5")
}

func TestLogHeader(t *testing.T) {
	fd, err := os.Open(fixturePath + "/ib_logfile0")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	err = p.readRedoLogFileHeader(fd)
	if err != nil {
		panic(err)
	}

	err = p.readRedoLogFileCheckpoint(fd)
	if err != nil {
		panic(err)
	}
}
