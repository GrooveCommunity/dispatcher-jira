module github.com/GrooveCommunity/dispatcher-jira

go 1.16

require (
	github.com/GrooveCommunity/glib-cloud-storage v0.0.0
	github.com/GrooveCommunity/glib-noc-event-structs v0.0.0
	github.com/andygrunwald/go-jira v1.13.0
	github.com/gorilla/mux v1.8.0
	google.golang.org/api v0.48.0
)

replace (
	github.com/GrooveCommunity/glib-cloud-storage v0.0.0 => /home/dejair/go/src/glib-cloud-storage
	github.com/GrooveCommunity/glib-noc-event-structs v0.0.0 => /home/dejair/go/src/glib-noc-event-structs
)
