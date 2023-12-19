package gerrit_commit

// topic = "commit"
type commitMerged struct {
	ChangeId               string
	Project                string
	Branch                 string
	Commiter               string
	Insertions             uint
	Deletions              uint
	TotalCommentCount      uint
	UnresolvedCommentCount uint
	ScanCommentCount       uint
	JiraNumber             string
	ReviewBlockCount       uint
	VerifyBlockCount       uint
}
