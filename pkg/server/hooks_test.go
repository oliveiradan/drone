package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/drone/drone/pkg/config"
	"github.com/drone/drone/pkg/remote/builtin/github"
	"github.com/drone/drone/pkg/server/recorder"
	"github.com/drone/drone/pkg/server/session"
	"github.com/drone/drone/pkg/store/mock"

	. "github.com/drone/drone/Godeps/_workspace/src/github.com/franela/goblin"
	"github.com/drone/drone/Godeps/_workspace/src/github.com/gin-gonic/gin"
	//"github.com/drone/drone/Godeps/_workspace/src/github.com/stretchr/testify/mock"
	//"github.com/drone/drone/Godeps/_workspace/src/gopkg.in/yaml.v2"

	queue "github.com/drone/drone/pkg/queue/builtin"
	common "github.com/drone/drone/pkg/types"
)

func TestHooks(t *testing.T) {
	store := new(mocks.Store)
	//
	g := Goblin(t)
	g.Describe("Hooks", func() {

		g.It("Should post hooks", func() {
			//
			buildList := []*common.Build{
				&common.Build{
					CommitID: 1,
					State:    "success",
					ExitCode: 0,
					Sequence: 1,
				},
				&common.Build{
					CommitID: 3,
					State:    "error",
					ExitCode: 1,
					Sequence: 2,
				},
			}
			commit1 := &common.Commit{
				RepoID: 1,
				State:  common.StateSuccess,
				Ref:    "refs/heads/master",
				Sha:    "14710626f22791619d3b7e9ccf58b10374e5b76d",
				Builds: buildList,
			}
			user1 := &common.User{
				ID:    1,
				Login: "oliveiradan",
				Name:  "Daniel Oliveira",
				Email: "doliveirabrz@gmail.com",
				Token: "e1c372bc477d38972c54b1794bdf3932",
			}
			repo1 := &common.Repo{
				UserID:   1,
				Owner:    "oliveiradan",
				Name:     "drone-test1",
				FullName: "oliveiradan/drone-test1",
			}
			config1 := &config.Config{}
			config1.Auth.Client = "87e2bdf99eece72c73c1"
			config1.Auth.Secret = "6b4031674ace18723ac41f58d63bff69276e5d1b"
			remote1 := github.New(config1)
			//remote1.
			queue1 := queue.New()
			hook1 := &common.Hook{
				Repo:   repo1,
				Commit: commit1,
			}
			config1.Session.Secret = "Otto"
			session1 := session.New(config1)
			fakeYMLFile := fmt.Sprintf(`[{"type": "file",
"encoding": "base64",
"size": 5362,
"name": "README.md",
"path": "README.md",
"content": "encoded content ...",
"sha": "3d21ec53a331a6f037a91c368710b99387d012c1",
"url": "https://api.github.com/repos/octokit/octokit.rb/contents/README.md",
"git_url": "https://api.github.com/repos/octokit/octokit.rb/git/blobs/3d21ec53a331a6f037a91c368710b99387d012c1",
"html_url": "https://github.com/octokit/octokit.rb/blob/master/README.md",
"download_url": "https://raw.githubusercontent.com/octokit/octokit.rb/master/README.md",
"_links": {
"git": "https://api.github.com/repos/octokit/octokit.rb/git/blobs/3d21ec53a331a6f037a91c368710b99387d012c1",
"self": "https://api.github.com/repos/octokit/octokit.rb/contents/README.md",
"html": "https://github.com/octokit/octokit.rb/blob/master/README.md",
"owner": "oliveiradan",
"Name":  "drone-test1"
}
}]`)
			bufYMLFile, _ := json.Marshal(&fakeYMLFile)
			//
			//type Matrix map[string][]string
			//mtxData1 := struct {
			//	Matrix map[string][]string
			//}{}
			//err :=
			//yaml.Unmarshal([]byte(bufYMLFile), &mtxData1)
			//parseCond1 := struct {
			//	Condition *common.Condition `yaml:"when"`
			//}{}
			//err =
			//yaml.Unmarshal([]byte(bufYMLFile), &parseCond1)

			// GET /api/hook
			rw := recorder.New()
			ctx := &gin.Context{Engine: gin.Default(), Writer: rw}
			//ctx.Params = append(ctx.Params, gin.Param{Key: "number", Value: "1"})
			//
			urlBase := "/api/hook"
			//urlString := (repo1.Owner + "/" + repo1.Name + "/builds" + "/1")
			urlFull := urlBase //(urlBase + urlString)
			//
			buf, _ := json.Marshal(&hook1)
			ctx.Request, _ = http.NewRequest("GET", urlFull, bytes.NewBuffer(buf))
			ctx.Request.Header.Set("Content-Type", "application/json")
			//
			ctx.Set("datastore", store)
			ctx.Set("repo", repo1)
			ctx.Set("remote", remote1)
			ctx.Set("queue", queue1)
			ctx.Set("settings", config1)
			ctx.Set("session", session1)
			// Start mock
			store.On("RepoName", hook1.Repo.Owner, hook1.Repo.Name).Return(repo1, nil).Once()
			store.On("User", repo1.UserID).Return(user1, nil).Once()
			store.On("AddCommit", commit1).Return(nil).Once()
			PostHook(ctx)
			//
			//var readerOut []byte //bytes.Buffer //[]byte
			readerOut := &common.Hook{}
			json.Unmarshal(rw.Body.Bytes(), &readerOut)
			g.Assert(rw.Code).Equal(200)
			fmt.Println("YML File:", bufYMLFile)
		})
	})
}
