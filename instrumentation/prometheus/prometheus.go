// Package prometheus implements Instrumentation against exported Prometheus
// metrics.
package prometheus

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/soundcloud/roshi/instrumentation"
)

// Satisfaction guaranteed.
var _ instrumentation.Instrumentation = PrometheusInstrumentation{}

// PrometheusInstrumentation holds metrics for all instrumented methods.
type PrometheusInstrumentation struct {
	insertCallCount             prometheus.Counter
	insertRecordCount           prometheus.Counter
	insertCallDuration          prometheus.Histogram
	insertRecordDuration        prometheus.Histogram
	insertQuorumFailureCount    prometheus.Counter
	selectCallCount             prometheus.Counter
	selectKeysCount             prometheus.Counter
	selectSendToCount           prometheus.Counter
	selectFirstResponseDuration prometheus.Histogram
	selectPartialErrorCount     prometheus.Counter
	selectBlockingDuration      prometheus.Histogram
	selectOverheadDuration      prometheus.Histogram
	selectDuration              prometheus.Histogram
	selectSendAllPromotionCount prometheus.Counter
	selectRetrievedCount        prometheus.Counter
	selectReturnedCount         prometheus.Counter
	selectRepairNeededCount     prometheus.Counter
	deleteCallCount             prometheus.Counter
	deleteRecordCount           prometheus.Counter
	deleteCallDuration          prometheus.Histogram
	deleteRecordDuration        prometheus.Histogram
	deleteQuorumFailureCount    prometheus.Counter
	repairCallCount             prometheus.Counter
	repairRequestCount          prometheus.Counter
	repairDiscardedCount        prometheus.Counter
	repairWriteSuccessCount     prometheus.Counter
	repairWriteFailureCount     prometheus.Counter
	walkKeysCount               prometheus.Counter
}

// New returns a new Instrumentation that prints metrics to the passed
// io.Writer. All metrics are prefixed with an appropriate bucket name, and
// take the form e.g. "insert.record.count 10".
func New(prefix string) PrometheusInstrumentation {
	i := PrometheusInstrumentation{
		insertCallCount:             prometheus.NewCounter(),
		insertRecordCount:           prometheus.NewCounter(),
		insertCallDuration:          prometheus.NewDefaultHistogram(),
		insertRecordDuration:        prometheus.NewDefaultHistogram(),
		insertQuorumFailureCount:    prometheus.NewCounter(),
		selectCallCount:             prometheus.NewCounter(),
		selectKeysCount:             prometheus.NewCounter(),
		selectSendToCount:           prometheus.NewCounter(),
		selectFirstResponseDuration: prometheus.NewDefaultHistogram(),
		selectPartialErrorCount:     prometheus.NewCounter(),
		selectBlockingDuration:      prometheus.NewDefaultHistogram(),
		selectOverheadDuration:      prometheus.NewDefaultHistogram(),
		selectDuration:              prometheus.NewDefaultHistogram(),
		selectSendAllPromotionCount: prometheus.NewCounter(),
		selectRetrievedCount:        prometheus.NewCounter(),
		selectReturnedCount:         prometheus.NewCounter(),
		selectRepairNeededCount:     prometheus.NewCounter(),
		deleteCallCount:             prometheus.NewCounter(),
		deleteRecordCount:           prometheus.NewCounter(),
		deleteCallDuration:          prometheus.NewDefaultHistogram(),
		deleteRecordDuration:        prometheus.NewDefaultHistogram(),
		deleteQuorumFailureCount:    prometheus.NewCounter(),
		repairCallCount:             prometheus.NewCounter(),
		repairRequestCount:          prometheus.NewCounter(),
		repairDiscardedCount:        prometheus.NewCounter(),
		repairWriteSuccessCount:     prometheus.NewCounter(),
		repairWriteFailureCount:     prometheus.NewCounter(),
		walkKeysCount:               prometheus.NewCounter(),
	}

	prometheus.Register(
		prefix+"insert_call_count",
		"How many insert calls have been made.",
		prometheus.NilLabels,
		i.insertCallCount,
	)
	prometheus.Register(
		prefix+"insert_record_count",
		"How many records have been inserted.",
		prometheus.NilLabels,
		i.insertRecordCount,
	)
	prometheus.Register(
		prefix+"insert_call_duration_nanoseconds",
		"Insert duration per-call.",
		prometheus.NilLabels,
		i.insertCallDuration,
	)
	prometheus.Register(
		prefix+"insert_record_duration_nanoseconds",
		"Insert duration per-record.",
		prometheus.NilLabels,
		i.insertRecordDuration,
	)
	prometheus.Register(
		prefix+"insert_quorum_failure_count",
		"Insert quorum failure count.",
		prometheus.NilLabels,
		i.insertQuorumFailureCount,
	)
	prometheus.Register(
		prefix+"select_call_count",
		"How many select calls have been made.",
		prometheus.NilLabels,
		i.selectCallCount,
	)
	prometheus.Register(
		prefix+"select_keys_count",
		"How many keys have been selected.",
		prometheus.NilLabels,
		i.selectKeysCount,
	)
	prometheus.Register(
		prefix+"select_send_to_count",
		"How many clusters have received select calls.",
		prometheus.NilLabels,
		i.selectSendToCount,
	)
	prometheus.Register(
		prefix+"select_first_response_duration_nanoseconds",
		"Select first response duration.",
		prometheus.NilLabels,
		i.selectFirstResponseDuration,
	)
	prometheus.Register(
		prefix+"select_partial_error_count",
		"How many partial errors have occurred in selects.",
		prometheus.NilLabels,
		i.selectPartialErrorCount,
	)
	prometheus.Register(
		prefix+"select_blocking_duration_nanoseconds",
		"Select blocking duration.",
		prometheus.NilLabels,
		i.selectBlockingDuration,
	)
	prometheus.Register(
		prefix+"select_overhead_duration_nanoseconds",
		"Select overhead duration.",
		prometheus.NilLabels,
		i.selectOverheadDuration,
	)
	prometheus.Register(
		prefix+"select_duration_nanoseconds",
		"Overall select duration.",
		prometheus.NilLabels,
		i.selectDuration,
	)
	prometheus.Register(
		prefix+"select_send_all_promotion_count",
		"How many select requests were promoted to a send-all, in appropriate read strategies.",
		prometheus.NilLabels,
		i.selectSendAllPromotionCount,
	)
	prometheus.Register(
		prefix+"select_retrieved_count",
		"How many key-score-member tuples have been retrieved from clusters by select calls.",
		prometheus.NilLabels,
		i.selectRetrievedCount,
	)
	prometheus.Register(
		prefix+"select_returned_count",
		"How many key-score-member tuples have been returned to clients by select calls.",
		prometheus.NilLabels,
		i.selectReturnedCount,
	)
	prometheus.Register(
		prefix+"select_repair_needed_count",
		"How many repairs have been detected and requested by select calls.",
		prometheus.NilLabels,
		i.selectRepairNeededCount,
	)
	prometheus.Register(
		prefix+"delete_call_count",
		"How many delete calls have been made.",
		prometheus.NilLabels,
		i.deleteCallCount,
	)
	prometheus.Register(
		prefix+"delete_record_count",
		"How many records have been deleted in delete calls.",
		prometheus.NilLabels,
		i.deleteRecordCount,
	)
	prometheus.Register(
		prefix+"delete_call_duration_nanoseconds",
		"Delete duration, per-call.",
		prometheus.NilLabels,
		i.deleteCallDuration,
	)
	prometheus.Register(
		prefix+"delete_record_duration_nanoseconds",
		"Delete duration, per-record.",
		prometheus.NilLabels,
		i.deleteRecordDuration,
	)
	prometheus.Register(
		prefix+"delete_quorum_failure_count",
		"Delete quorum failure count.",
		prometheus.NilLabels,
		i.deleteQuorumFailureCount,
	)
	prometheus.Register(
		prefix+"repair_call_count",
		"How many repair calls have been made.",
		prometheus.NilLabels,
		i.repairCallCount,
	)
	prometheus.Register(
		prefix+"repair_request_count",
		"How many key-member tuples have been repaired.",
		prometheus.NilLabels,
		i.repairRequestCount,
	)
	prometheus.Register(
		prefix+"repair_discarded_count",
		"How many repair calls have been discarded due to rate or buffer limits.",
		prometheus.NilLabels,
		i.repairDiscardedCount,
	)
	prometheus.Register(
		prefix+"repair_write_success_count",
		"Repair write success count.",
		prometheus.NilLabels,
		i.repairWriteSuccessCount,
	)
	prometheus.Register(
		prefix+"repair_write_failure_count",
		"Repair write failure count.",
		prometheus.NilLabels,
		i.repairWriteFailureCount,
	)
	prometheus.Register(
		prefix+"walk_keys_count",
		"How many keys have been walked by the walker process.",
		prometheus.NilLabels,
		i.walkKeysCount,
	)

	return i
}

// Install installs the Prometheus handlers, so the metrics are available.
func (i PrometheusInstrumentation) Install(pattern string, mux *http.ServeMux) {
	mux.Handle(pattern, prometheus.DefaultHandler)
}

// InsertCall satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) InsertCall() {
	i.insertCallCount.Increment(prometheus.NilLabels)
}

// InsertRecordCount satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) InsertRecordCount(n int) {
	i.insertRecordCount.IncrementBy(prometheus.NilLabels, float64(n))
}

// InsertCallDuration satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) InsertCallDuration(d time.Duration) {
	i.insertCallDuration.Add(prometheus.NilLabels, float64(d.Nanoseconds()))
}

// InsertRecordDuration satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) InsertRecordDuration(d time.Duration) {
	i.insertRecordDuration.Add(prometheus.NilLabels, float64(d.Nanoseconds()))
}

// InsertQuorumFailure satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) InsertQuorumFailure() {
	i.insertQuorumFailureCount.Increment(prometheus.NilLabels)
}

// SelectCall satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) SelectCall() {
	i.selectCallCount.Increment(prometheus.NilLabels)
}

// SelectKeys satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) SelectKeys(n int) {
	i.selectKeysCount.IncrementBy(prometheus.NilLabels, float64(n))
}

// SelectSendTo satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) SelectSendTo(n int) {
	i.selectSendToCount.IncrementBy(prometheus.NilLabels, float64(n))
}

// SelectFirstResponseDuration satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) SelectFirstResponseDuration(d time.Duration) {
	i.selectFirstResponseDuration.Add(prometheus.NilLabels, float64(d.Nanoseconds()))
}

// SelectPartialError satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) SelectPartialError() {
	i.selectPartialErrorCount.Increment(prometheus.NilLabels)
}

// SelectBlockingDuration satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) SelectBlockingDuration(d time.Duration) {
	i.selectBlockingDuration.Add(prometheus.NilLabels, float64(d.Nanoseconds()))
}

// SelectOverheadDuration satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) SelectOverheadDuration(d time.Duration) {
	i.selectOverheadDuration.Add(prometheus.NilLabels, float64(d.Nanoseconds()))
}

// SelectDuration satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) SelectDuration(d time.Duration) {
	i.selectDuration.Add(prometheus.NilLabels, float64(d.Nanoseconds()))
}

// SelectSendAllPromotion satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) SelectSendAllPromotion() {
	i.selectSendAllPromotionCount.Increment(prometheus.NilLabels)
}

// SelectRetrieved satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) SelectRetrieved(n int) {
	i.selectRetrievedCount.IncrementBy(prometheus.NilLabels, float64(n))
}

// SelectReturned satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) SelectReturned(n int) {
	i.selectReturnedCount.IncrementBy(prometheus.NilLabels, float64(n))
}

// SelectRepairNeeded satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) SelectRepairNeeded(n int) {
	i.selectRepairNeededCount.IncrementBy(prometheus.NilLabels, float64(n))
}

// DeleteCall satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) DeleteCall() {
	i.deleteCallCount.Increment(prometheus.NilLabels)
}

// DeleteRecordCount satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) DeleteRecordCount(n int) {
	i.deleteRecordCount.IncrementBy(prometheus.NilLabels, float64(n))
}

// DeleteCallDuration satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) DeleteCallDuration(d time.Duration) {
	i.deleteCallDuration.Add(prometheus.NilLabels, float64(d.Nanoseconds()))
}

// DeleteRecordDuration satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) DeleteRecordDuration(d time.Duration) {
	i.deleteRecordDuration.Add(prometheus.NilLabels, float64(d.Nanoseconds()))
}

// DeleteQuorumFailure satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) DeleteQuorumFailure() {
	i.deleteQuorumFailureCount.Increment(prometheus.NilLabels)
}

// RepairCall satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) RepairCall() {
	i.repairCallCount.Increment(prometheus.NilLabels)
}

// RepairRequest satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) RepairRequest(n int) {
	i.repairRequestCount.IncrementBy(prometheus.NilLabels, float64(n))
}

// RepairDiscarded satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) RepairDiscarded(n int) {
	i.repairDiscardedCount.IncrementBy(prometheus.NilLabels, float64(n))
}

// RepairWriteSuccess satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) RepairWriteSuccess(n int) {
	i.repairWriteSuccessCount.IncrementBy(prometheus.NilLabels, float64(n))
}

// RepairWriteFailure satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) RepairWriteFailure(n int) {
	i.repairWriteFailureCount.IncrementBy(prometheus.NilLabels, float64(n))
}

// WalkKeys satisfies the Instrumentation interface.
func (i PrometheusInstrumentation) WalkKeys(n int) {
	i.walkKeysCount.IncrementBy(prometheus.NilLabels, float64(n))
}
