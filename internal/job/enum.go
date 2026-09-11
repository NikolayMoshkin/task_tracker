package job

type JobType string

const (
	ADD              JobType = "add"
	UPDATE           JobType = "update"
	DELETE           JobType = "delete"
	LIST             JobType = "list"
	MARK_DONE        JobType = "mark-done"
	MARK_IN_PROGRESS JobType = "mark-in-progress"
)
