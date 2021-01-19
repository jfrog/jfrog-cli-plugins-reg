package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ReleaseRes struct {
	Name string
}

type ReleaseReq struct {
	Owner string
	Repo  string
}

type IssuesReq struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

func GetLatestRelease(owner, repo string) (version string, err error) {
	uri := "https://api.github.com/repos/" + owner + "/" + repo + "/releases?per_page=1"
	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var githubResponse []ReleaseRes
	if err := json.Unmarshal(content, &githubResponse); err != nil {
		return "", errors.New("failed to unmarshal GitHub response:" + string(content))
	}
	return strings.TrimSpace(githubResponse[0].Name), nil
}

func OpenIssue(owner, repo, token string, issueReq IssuesReq) error {
	if token == "" {
		return errors.New("missing Token to open an issue")
	}
	client := &http.Client{}
	client.Timeout = time.Second * 15
	uri := "https://api.github.com/repos/" + owner + "/" + repo + "/issues"
	data, err := json.Marshal(&issueReq)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return errors.New("Github response: " + resp.Status + " while opening an issue")
	}
	return nil
}
