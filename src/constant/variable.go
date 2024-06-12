package constant

type ApplicationStatus string
type JobStatus string
type Role string

const (
	Pending   ApplicationStatus = "PENDING"
	Interview ApplicationStatus = "INTERVIEW"
	Accepted  ApplicationStatus = "ACCEPTED"
	Rejected  ApplicationStatus = "REJECTED"

	Available JobStatus = "AVAILABLE"
	Closed    JobStatus = "CLOSED"

	Recruiter Role = "recruiter"
	Talent    Role = "talent"
)
