package app

import "github.com/prometheus/client_golang/prometheus"

var (
	GroupOpProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "group_service_operations_received",
		Help: "Total number of request to group_service",
	})

	StudentOpProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "student_service_operations_received",
		Help: "Total number of request to student_service",
	})
)

func Init() {
	prometheus.MustRegister(GroupOpProcessed)
	prometheus.MustRegister(StudentOpProcessed)
}
