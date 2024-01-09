package report

import (
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

type ReportData struct {
	name       string
	data       [][]string
	reportFile string
}

func InitGenesisReport(name string) *ReportData {
	return &ReportData{
		name:       name,
		data:       [][]string{{"Variable Name", "Variable Type", "Storage Layout", "Key", "Value"}},
		reportFile: fmt.Sprintf("./out/%s-genesis-report.csv", name),
	}
}

func (d *ReportData) AddGenesisKV(name, typeString, layout, key, value string) {
	d.data = append(d.data, []string{name, typeString, layout, key, value})
}

func InitVerifyReport(name string) *ReportData {
	return &ReportData{
		name:       name,
		data:       [][]string{{"Test Case", "Status", "From Network", "To Network"}},
		reportFile: fmt.Sprintf("./out/%s-report.csv", name),
	}
}

func (d *ReportData) AddStorageResult(key, value string, from, to common.Hash) {
	success := true
	if from.Cmp(to) != 0 {
		success = false
	}
	testcase := fmt.Sprintf("\"%s\":\"%s\"", key, value)
	d.data = append(d.data, []string{testcase, strconv.FormatBool(success), from.Hex(), to.Hex()})
}

func (d *ReportData) AddFunctionResult(method string, args interface{}, from, to []byte) {
	success := true
	if !bytes.Equal(from, to) {
		success = false
	}
	var testcase string
	if args != nil {
		testcase = fmt.Sprintf("method: %s, args: %s", method, fmt.Sprintf("%v", args))
	} else {
		testcase = fmt.Sprintf("method: %s", method)
	}
	d.data = append(d.data, []string{testcase, strconv.FormatBool(success), "0x" + hex.EncodeToString(from), "0x" + hex.EncodeToString(to)})
}

func (d *ReportData) ReportVerifyResult() {
	f, err := os.Create(d.reportFile)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	w := csv.NewWriter(f)
	err = w.WriteAll(d.data)
	if err != nil {
		panic(err)
	}
}
