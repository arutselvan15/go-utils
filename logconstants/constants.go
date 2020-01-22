package logconstants

const (
	// events

	// OnAdd OnAdd
	OnAdd = "OnAdd"
	// OnUpdate OnUpdate
	OnUpdate = "OnUpdate"
	// OnDelete OnDelete
	OnDelete = "OnDelete"

	// actions

	// Create Create
	Create = "Create"
	// Update Update
	Update = "Update"
	// Read Read
	Read = "Read"
	// Delete Delete
	Delete = "Delete"

	// subcomponents

	// Audit Audit
	Audit = "Audit"
	// Mutate Mutate
	Mutate = "Mutate"
	// Validate Validate
	Validate = "Validate"

	// general state of requests

	// Pending Pending
	Pending = "Pending"
	// Start Start
	Start = "Start"
	// InProgress InProgress
	InProgress = "InProgress"
	// Retry Retry
	Retry = "Retry"
	// End End
	End = "End"

	// disposition state

	// Success Success
	Success = "Success"
	// Failure Failure
	Failure = "Failure"
)
