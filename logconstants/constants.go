// Package logconstants lib
package logconstants

const (
	// operation

	// Create Create
	Create = "create"
	// Update Update
	Update = "update"
	// Read Read
	Read = "read"
	// Delete Delete
	Delete = "delete"

	// resource state

	// Received Received
	Received = "received"
	// Pending Pending
	Pending = "pending"
	//Processing Processing
	Processing = "processing"
	//Deleting Deleting
	Deleting = "deleting"
	// Successful Success
	Successful = "successful"
	// Failed Failure
	Failed = "failed"
	// Retry Retry
	Retry = "retry"
	// Unknown Unknown
	Unknown = "unknown"
	//Ignored Ignored
	Ignored = "ignored"

	// step state

	// Start Start
	Start = "start"
	// Skip Skip
	Skip = "skip"
	// InProgress InProgress
	InProgress = "inprogress"
	// Complete Complete
	Complete = "complete"
	// Error error
	Error = "error"

	// Audit Audit
	Audit = "audit"
	// Mutate Mutate
	Mutate = "mutate"
	// Validate Validate
	Validate = "validate"
)
