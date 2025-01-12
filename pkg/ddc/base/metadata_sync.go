package base

import (
	"errors"
	"strconv"
	"time"

	"github.com/fluid-cloudnative/fluid/pkg/metrics"
	"github.com/fluid-cloudnative/fluid/pkg/utils"
	"github.com/go-logr/logr"
)

// MetadataSyncResult describes result for asynchronous metadata sync
type MetadataSyncResult struct {
	Done      bool
	StartTime time.Time
	UfsTotal  string
	FileNum   string
	Err       error
}

// SafeClose closes the metadataSyncResultChannel but ignores panic when the channel is already closed.
// Returns true if the channel is already closed.
func SafeClose(ch chan MetadataSyncResult) (closed bool) {
	if ch == nil {
		return
	}
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()

	close(ch)
	return false
}

// SafeSend sends result to the metadataSyncResultChannel but ignores panic when the channel is already closed
// Returns true if the channel is already closed.
func SafeSend(ch chan MetadataSyncResult, result MetadataSyncResult) (closed bool) {
	if ch == nil {
		return
	}
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()

	ch <- result
	return false
}

// RecordDatasetMetrics records dataset-related metrics from the given MetadataSyncResult
func RecordDatasetMetrics(result MetadataSyncResult, datasetNamespace, datasetName string, log logr.Logger) {
	if len(datasetNamespace) == 0 {
		argErr := errors.New("invalid argument: datasetNamespace should not be empty")
		log.Error(argErr, "fail to validate RecordDatasetMetrics arguments")
		return
	}

	if len(datasetName) == 0 {
		argErr := errors.New("invalid argument: datasetName should not be empty")
		log.Error(argErr, "fail to validate RecordDatasetMetrics arguments")
		return
	}

	if len(result.UfsTotal) != 0 {
		if ufsTotal, parseErr := utils.FromHumanSize(result.UfsTotal); parseErr == nil {
			metrics.GetDatasetMetrics(datasetNamespace, datasetName).SetUFSTotalSize(float64(ufsTotal))
		} else {
			log.Error(parseErr, "fail to parse result.UfsTotal", "result.UfsTotal", result.UfsTotal)
		}
	}

	if len(result.FileNum) != 0 {
		if fileNum, parseErr := strconv.Atoi(result.FileNum); parseErr == nil {
			metrics.GetDatasetMetrics(datasetNamespace, datasetName).SetUFSFileNum(float64(fileNum))
		} else {
			log.Error(parseErr, "fail to atoi result.FileNum", "result.FileNum", result.FileNum)
		}
	}
}
