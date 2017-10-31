package benchmark

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/types"
)

const (
	transferRateLabel = "transfer rate [MB/s]"
)

type measurementSeries map[string]float64

type myReporter struct {
	Reporter

	results map[string]measurementSeries
}

var _ Reporter = &myReporter{}

func (r *myReporter) SpecSuiteWillBegin(config.GinkgoConfigType, *types.SuiteSummary) {}
func (r *myReporter) BeforeSuiteDidRun(*types.SetupSummary)                           {}
func (r *myReporter) SpecWillRun(*types.SpecSummary)                                  {}
func (r *myReporter) AfterSuiteDidRun(*types.SetupSummary)                            {}
func (r *myReporter) SpecSuiteDidEnd(*types.SuiteSummary)                             {}

func (r *myReporter) SpecDidComplete(specSummary *types.SpecSummary) {
	if !specSummary.IsMeasurement {
		return
	}
	test := specSummary.ComponentTexts[4]
	cond := specSummary.ComponentTexts[2]
	transferRate, ok := specSummary.Measurements[transferRateLabel]
	if !ok {
		return
	}
	r.addResult(cond, test, transferRate.Average)
}

func (r *myReporter) addResult(cond, ver string, transferRate float64) {
	if r.results == nil {
		r.results = make(map[string]measurementSeries)
	}
	if _, ok := r.results[cond]; !ok {
		r.results[cond] = make(measurementSeries)
	}
	r.results[cond][ver] = transferRate
}

func (r *myReporter) printResult() {
	table := tablewriter.NewWriter(os.Stdout)
	header := []string{""}
	for _, series := range r.results {
		for label := range series {
			header = append(header, label)
		}
		break
	}
	table.SetHeader(header)
	table.SetCaption(true, fmt.Sprintf("Based on %d samples of %d MB.\nAll values in MB/s.", samples, size))
	table.SetAutoFormatHeaders(false)

	for _, cond := range conditions {
		data := make([]string, len(header))
		data[0] = cond.Description

		for i := 1; i < len(header); i++ {
			val := r.results[cond.Description][header[i]]
			var out string
			if val == 0 {
				out = "-"
			} else {
				out = fmt.Sprintf("%.2f", val)
			}
			data[i] = out
		}
		table.Append(data)
	}
	table.Render()
}
