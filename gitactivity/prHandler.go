package gitactivity

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/cpg1111/maestrod/lifecycle"
)

// PRSubHandler is the subhandler for Pull Request events
type PRSubHandler struct {
	SubHandler
	Error ErrorHandler
	Queue *lifecycle.Queue
}

// NewPRHandler returns a pointer to a RouteHandler for Pull Request Events
func NewPRHandler(queue *lifecycle.Queue) *RouteHandler {
	sub := PRSubHandler{
		Error: ErrorHandler{},
		Queue: queue,
	}
	return &RouteHandler{
		Route: "/pullrequest",
		sub:   sub,
	}
}

type branchPayload struct {
	Label string      `json:"label"`
	Ref   string      `json:"ref"`
	SHA   string      `json:"sha"`
	User  UserPayload `json:"user"`
	Repo  RepoPayload `json:"repo"`
}

type prPayload struct {
	URL            string        `json:"url"`
	ID             uint          `json:"id"`
	HTMLURL        string        `json:"html_url"`
	DIFFURL        string        `json:"diff_url"`
	PATCHURL       string        `json:"patch_url"`
	ISSUEURL       string        `json:"issue_url"`
	Number         uint          `json:"number"`
	State          string        `json:"state"`
	Locked         bool          `json:"locked"`
	Title          string        `json:"title"`
	User           UserPayload   `json:"user"`
	Body           string        `json:"body"`
	CreatedAt      string        `json:"created_at"`
	UpdatedAt      string        `json:"updated_at"`
	ClosedAt       string        `json:"closed_at"`
	MergedAt       string        `json:"merged_at"`
	MergeCommitSha string        `json:"merge_commit_sha"`
	Assignee       UserPayload   `json:"assignee"`
	Head           branchPayload `json:"head"`
	Base           branchPayload `json:"base"`
}

// PRPayload is the Github api's Pull Request Event Object
type PRPayload struct {
	Action      string    `json:"action"`
	Number      int       `json:"number"`
	PullRequest prPayload `json:"pull_request"`
}

// Get is PRSubHandler's HTTP GET action
func (pr PRSubHandler) Get(res http.ResponseWriter, req *http.Request) {
	pr.Error.handle404(res, req)
}

// Post is PRSubHandler's HTTP POST action
func (pr PRSubHandler) Post(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	payload := &PRPayload{}
	decoder.Decode(payload)
	pullRequest := payload.PullRequest
	branchName := strings.Replace(pullRequest.Head.Ref, "refs/heads/", "", -1)
	log.Println(
		"Adding Job to Queue: ",
		pullRequest.Head.Repo.FullName,
		branchName,
		pullRequest.Head.SHA,
	)
	pr.Queue.Add(
		pullRequest.Head.Repo.FullName,
		branchName,
		pullRequest.Base.SHA,
		pullRequest.Head.SHA,
	)
	resp := PostResp{Status: 201, Message: "Created"}
	json.NewEncoder(res).Encode(resp)
}

// Put is PRSubHandler's HTTP PUT action
func (pr PRSubHandler) Put(res http.ResponseWriter, req *http.Request) {
	pr.Error.handle404(res, req)
}

// Patch is PRSubHandler's HTTP PATCH action
func (pr PRSubHandler) Patch(res http.ResponseWriter, req *http.Request) {
	pr.Error.handle404(res, req)
}

// Delete is PRSubHandler's HTTP DELETE action
func (pr PRSubHandler) Delete(res http.ResponseWriter, req *http.Request) {
	pr.Error.handle404(res, req)
}
