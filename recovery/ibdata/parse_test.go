package ibdata

import (
	"flag"
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
	p = NewParse()
	_, filename, _, _ := runtime.Caller(0)

	flag.Parse()
	devPath = filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	fixturePath = devPath + "/cmd/test/fixture/" + *mysqlRelease + "_" + *mysqlVersion

	// set logs to stderr, and log-level = trace
	flag.Set("logtostderr", "true")
	flag.Set("v", "5")

	m.Run()
}

// go test -v -timeout 30s -run ^TestParse$ github.com/zbdba/db-recovery/recovery/ibdata > TestParse.log 2>&1
func TestParse(t *testing.T) {
	flag.Set("v", "3")
	err := p.ParseDictPage(fixturePath + "/ibdata1")
	if err != nil {
		panic(err)
	}

	err = p.ParseDataPage(fixturePath+"/ibdata1", "test", "test_int", true)
	if err != nil {
		panic(err)
	}
	flag.Set("v", "5")
}

func TestParseFile(t *testing.T) {
	// test parse ibdata
	pages, err := p.parseFile(fixturePath + "/ibdata1")
	if err != nil {
		panic(err)
	}

	if len(pages) == 0 {
		panic("no pages found, parse error!")
	}

	// test parse ibd
	pages, err = p.parseFile(fixturePath + "/test/test_int.ibd")
	if err != nil {
		panic(err)
	}

	if len(pages) == 0 {
		panic("no pages found, parse error!")
	}
}

func TestAddColumns(t *testing.T) {
	var columns []Columns
	// add first
	columns = addColumns(columns, 0, Columns{FieldName: "first"})

	// append last
	columns = addColumns(columns, 1, Columns{FieldName: "last"})

	// add into middle
	columns = addColumns(columns, 1, Columns{FieldName: "middle"})

	if len(columns) != 3 {
		panic("addColumns count error")
	}

	if columns[0].FieldName != "first" {
		panic("addColumns value error")
	}

	if columns[1].FieldName != "middle" {
		panic("addColumns value error")
	}

	if columns[2].FieldName != "last" {
		panic("addColumns value error")
	}
}
