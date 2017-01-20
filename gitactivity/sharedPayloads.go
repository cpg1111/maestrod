/*
Copyright 2016 Christian Grabowski All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
