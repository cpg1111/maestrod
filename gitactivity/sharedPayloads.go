package gitactivity

// UserPayload is the struct for Github's api user object
type UserPayload struct {
	Login     string `json:"login"`
	ID        uint   `json:"id"`
	AvatarURL string `json:"avatar_url"`
	URL       string `json:"url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
}

// RepoPayload is the struct for Github's api repo object
type RepoPayload struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	FullName    string      `json:"full_name"`
	Owner       UserPayload `json:"owner"`
	Description string      `json:"description"`
	Private     bool        `json:"private"`
	Fork        bool        `json:"fork"`
	HTMLURL     string      `json:"html_url"`
	URL         string      `json:"url"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	PushedAt    string      `json:"pushed_at"`
}

// PostResp is the standard response on a HTTP POST
type PostResp struct {
	Status  uint   `json:"status"`
	Message string `json:"message"`
}
