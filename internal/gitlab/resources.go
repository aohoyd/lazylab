package gitlab

import (
	"time"
)

type MergeRequestsService struct {
	client *Client
}

type MergeRequest struct {
	ID                          int64                 `json:"id"`
	IID                         int64                 `json:"iid"`
	ProjectID                   int64                 `json:"project_id"`
	Title                       string                `json:"title"`
	Description                 string                `json:"description"`
	State                       string                `json:"state"`
	CreatedAt                   string                `json:"created_at"`
	UpdatedAt                   *time.Time            `json:"updated_at"`
	MergedBy                    *User                 `json:"merged_by"`
	MergedAt                    *time.Time            `json:"merged_at"`
	ClosedBy                    *User                 `json:"closed_by"`
	ClosedAt                    *time.Time            `json:"closed_at"`
	TargetBranch                string                `json:"target_branch"`
	SourceBranch                string                `json:"source_branch"`
	UserNotesCount              int64                 `json:"user_notes_count"`
	Upvotes                     int64                 `json:"upvotes"`
	Downvotes                   int64                 `json:"downvotes"`
	Assignee                    *User                 `json:"assignee"`
	Author                      *User                 `json:"author"`
	Assignees                   []*User               `json:"assignees"`
	SourceProjectID             int64                 `json:"source_project_id"`
	TargetProjectID             int64                 `json:"target_project_id"`
	Labels                      []string              `json:"labels"`
	WorkInProgress              bool                  `json:"work_in_progress"`
	Milestone                   *Milestone            `json:"milestone"`
	MergeWhenPipelineSucceeds   bool                  `json:"merge_when_pipeline_succeeds"`
	MergeStatus                 string                `json:"merge_status"`
	SHA                         string                `json:"sha"`
	MergeCommitSHA              string                `json:"merge_commit_sha"`
	SquashCommitSHA             string                `json:"squash_commit_sha"`
	DiscussionLocked            bool                  `json:"discussion_locked"`
	ShouldRemoveSourceBranch    bool                  `json:"should_remove_source_branch"`
	ForceRemoveSourceBranch     bool                  `json:"force_remove_source_branch"`
	Reference                   string                `json:"reference"`
	References                  *References           `json:"references"`
	WebURL                      string                `json:"web_url"`
	TimeStats                   *TimeStats            `json:"time_stats"`
	Squash                      bool                  `json:"squash"`
	TaskCompletionStatus        *TaskCompletionStatus `json:"task_completion_status"`
	DiffRefs                    *DiffRefs             `json:"diff_refs"`
	HasConflicts                bool                  `json:"has_conflicts"`
	BlockingDiscussionsResolved bool                  `json:"blocking_discussions_resolved"`
	ApprovalsBeforeMerge        int64                 `json:"approvals_before_merge"`
}

type User struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Username  string     `json:"username"`
	State     string     `json:"state"`
	CreatedAt *time.Time `json:"created_at"`
	AvatarURL string     `json:"avatar_url"`
	WebURL    string     `json:"web_url"`
}

type Milestone struct {
	ID          int        `json:"id"`
	IID         int        `json:"iid"`
	ProjectID   int        `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	DueDate     *time.Time `json:"due_date"`
	State       string     `json:"state"`
	WebURL      string     `json:"web_url"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CreatedAt   *time.Time `json:"created_at"`
}

type References struct {
	Short    string `json:"short"`
	Relative string `json:"relative"`
	Full     string `json:"full"`
}

type TaskCompletionStatus struct {
	Count          int64 `json:"count"`
	CompletedCount int64 `json:"completed_count"`
}

type DiffRefs struct {
	BaseSha  string `json:"base_sha"`
	HeadSha  string `json:"head_sha"`
	StartSha string `json:"start_sha"`
}

type TimeStats struct {
	TimeEstimate        int64  `json:"time_estimate"`
	TotalTimeSpent      int64  `json:"total_time_spent"`
	HumanTimeEstimate   string `json:"human_time_estimate"`
	HumanTotalTimeSpent string `json:"human_total_time_spent"`
}

type ListMergeRequestsOptions struct {
	ApproverUsernames []string `url:"approver_usernames,brackets,omitempty" json:"approver_usernames,omitempty"`
	ApproverIDs       []int    `url:"approver_ids,brackets,omitempty" json:"approver_ids,omitempty"`
	AssigneeUsername  *string  `url:"assignee_username,omitempty" json:"assignee_username,omitempty"`
	AuthorUsername    *string  `url:"author_username,omitempty" json:"author_username,omitempty"`
	AuthorID          *int     `url:"author_id,omitempty" json:"author_id,omitempty"`
	State             *string  `url:"state,omitempty" json:"state,omitempty"`
	Scope             *string  `url:"scope,omitempty" json:"scope,omitempty"`
}

func (s *MergeRequestsService) ListMergeRequests(opt *ListMergeRequestsOptions) ([]*MergeRequest, *Response, error) {
	req, err := s.client.NewRequest("GET", "merge_requests", opt)
	if err != nil {
		return nil, nil, err
	}

	//fmt.Printf("%v %v %v\n", req.Method, req.URL, req.Host)

	var m []*MergeRequest
	resp, err := s.client.Do(req, &m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}
