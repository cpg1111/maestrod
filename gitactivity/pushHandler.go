package gitactivity

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/cpg1111/maestrod/lifecycle"
)

type PushSubHandler struct {
	SubHandler
	Error ErrorHandler
	Queue *lifecycle.Queue
}

func NewPushHandler(queue *lifecycle.Queue) *RouteHandler {
	sub := PushSubHandler{
		Error: ErrorHandler{},
		Queue: queue,
	}
	return &RouteHandler{
		Route: "/push",
		sub:   sub,
	}
}

type authorPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserName string `json:"username"`
}

type commitPayload struct {
	ID        string        `json:"id"`
	TreeID    string        `json:"tree_id"`
	Distinct  bool          `json:"distinct"`
	Message   string        `json:"message"`
	Timestamp string        `json:"timestamp"`
	URL       string        `json:"url"`
	Author    authorPayload `json:"author"`
	Commiter  authorPayload `json:"committer"`
	Added     []string      `json:"added"`
	Removed   []string      `json:"removed"`
	Modified  []string      `json:"modified"`
}

type repositoryPayload struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Name  string `json"name`
		Email string `json:"email"`
	} `json:"owner"`
	Private     bool   `json:"private"`
	HtmlURL     string `json:"html_url"`
	Description string `json:"description"`
	Fork        bool   `json:"fork"`
	URL         string `json:"url"`
	CreatedAt   uint   `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	PushedAt    uint   `json:"pushed_at"`
}

type pushPayload struct {
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	After      string `json:"after"`
	Created    bool   `json:"created"`
	Deleted    bool   `json:"deleted"`
	Forced     bool   `json:"forced"`
	BaseRef    string `json:"base_ref"`
	Compare    string
	Commits    []commitPayload   `json:"commits"`
	Repository repositoryPayload `json:"repository"`
	Pusher     authorPayload     `json:"pusher"`
}

func (p PushSubHandler) Get(res http.ResponseWriter, req *http.Request) {
	p.Error.handle404(res, req)
}

type postResp struct {
	Status  uint   `json:"status"`
	Message string `json:"message"`
}

func (p PushSubHandler) Post(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	payload := &pushPayload{}
	decoder.Decode(payload)
	branchName := strings.Replace(payload.Ref, "refs/heads/", "", -1)
	log.Println("Adding Job to Queue: ", payload.Repository.FullName, branchName, payload.Before)
	p.Queue.Add(payload.Repository.FullName, branchName, payload.Before, payload.After)
	resp := postResp{Status: 201, Message: "Created"}
	json.NewEncoder(res).Encode(resp)
}

func (p PushSubHandler) Put(res http.ResponseWriter, req *http.Request) {
	p.Error.handle404(res, req)
}

func (p PushSubHandler) Patch(res http.ResponseWriter, req *http.Request) {
	p.Error.handle404(res, req)
}

func (p PushSubHandler) Delete(res http.ResponseWriter, req *http.Request) {
	p.Error.handle404(res, req)
}
