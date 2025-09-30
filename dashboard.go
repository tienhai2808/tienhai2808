package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	username = "tienhai2808"
)

var (
	excludedRepos = map[string]bool{
		"odoo19": true, "wscms": true, "vicem": true, "parfumerie": true, "mixi": true,
		"hotel_booking": true, "ems": true, "de_gk_test_1": true, "de_gk_test_2": true,
		"de_thi_gk_2023-2024": true, "de_gk_104_48K22.1_2024-2025": true,
		"de_gk_105_48K22.1_2024-2025": true, "bookr-django": true, "django_faceid": true,
		"django_first_project": true, "first_project-be": true, "decision-tree-ptdlpython": true,
		"film_recommendation_system": true, "simple_data_mining": true, "big4_stock": true, "ffrc": true,
	}
)

type gqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

type RepoActivity struct {
	Name    string
	Commits int
}

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN không được thiết lập")
	}

	client := &http.Client{}

	fmt.Println("📊 Đang thu thập dữ liệu từ GitHub...")

	repos := fetchAllRepos(client, token)
	fmt.Printf("✓ Tìm thấy %d repos\n", len(repos))

	submodules := fetchAllSubmodules(client, token, repos)
	fmt.Printf("✓ Tìm thấy %d submodules\n", len(submodules))

	langStats, _ := fetchLanguageStats(client, token, repos, submodules)
	fmt.Printf("✓ Thu thập language stats\n")

	repoActivity := fetchRepoActivity(client, token, repos, submodules)
	fmt.Printf("✓ Thu thập commit activity (30 ngày)\n")

	generateREADME(langStats, repoActivity)
	fmt.Println("\n✅ README.md đã được cập nhật!")
}

func doRequest(client *http.Client, token, query string, vars map[string]any) ([]byte, error) {
	payload, _ := json.Marshal(gqlRequest{query, vars})
	req, _ := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func doRESTRequest(client *http.Client, token, url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func fetchAllRepos(client *http.Client, token string) []string {
	query := `query($username:String!,$after:String){user(login:$username){repositories(first:100,after:$after){nodes{name}pageInfo{hasNextPage endCursor}}}}`
	var repos []string
	after := ""

	for {
		vars := map[string]any{"username": username}
		if after != "" {
			vars["after"] = after
		}

		body, err := doRequest(client, token, query, vars)
		if err != nil {
			log.Fatal(err)
		}

		var resp struct {
			Data struct {
				User struct {
					Repositories struct {
						Nodes    []struct{ Name string }
						PageInfo struct {
							HasNextPage bool
							EndCursor   string
						}
					}
				}
			}
		}
		json.Unmarshal(body, &resp)

		for _, n := range resp.Data.User.Repositories.Nodes {
			if !excludedRepos[n.Name] {
				repos = append(repos, n.Name)
			}
		}

		if !resp.Data.User.Repositories.PageInfo.HasNextPage {
			break
		}
		after = resp.Data.User.Repositories.PageInfo.EndCursor
	}
	return repos
}

func fetchAllSubmodules(client *http.Client, token string, repos []string) [][2]string {
	var allSubmodules [][2]string
	seenSubmodules := make(map[string]bool)

	for _, repo := range repos {
		branches := []string{"main", "master", "develop"}

		for _, branch := range branches {
			url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/.gitmodules?ref=%s", username, repo, branch)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Accept", "application/vnd.github.v3.raw")

			resp, err := client.Do(req)
			if err != nil || resp.StatusCode != 200 {
				if resp != nil {
					resp.Body.Close()
				}
				continue
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			text := string(body)

			if text != "" {
				re := regexp.MustCompile(`url\s*=\s*[^\n]*github\.com[/:]([^/]+)/([^/\.\s]+)`)
				matches := re.FindAllStringSubmatch(text, -1)

				for _, m := range matches {
					if len(m) >= 3 {
						key := m[1] + "/" + m[2]
						if !seenSubmodules[key] {
							allSubmodules = append(allSubmodules, [2]string{m[1], m[2]})
							seenSubmodules[key] = true
						}
					}
				}
				break
			}
		}
	}

	return allSubmodules
}

func fetchLanguageStats(client *http.Client, token string, repos []string, subs [][2]string) (map[string]int, map[string]int) {
	query := `query($owner:String!,$repo:String!){repository(owner:$owner,name:$repo){languages(first:10){edges{node{name}size}}}}`
	langStats := make(map[string]int)
	repoStats := make(map[string]int)

	process := func(owner, repo string) {
		if owner == username && excludedRepos[repo] {
			return
		}

		vars := map[string]any{"owner": owner, "repo": repo}
		body, err := doRequest(client, token, query, vars)
		if err != nil {
			return
		}

		var resp struct {
			Data struct {
				Repository struct {
					Languages struct {
						Edges []struct {
							Node struct{ Name string }
							Size int
						}
					}
				}
			}
		}
		json.Unmarshal(body, &resp)

		repoTotal := 0
		for _, e := range resp.Data.Repository.Languages.Edges {
			langStats[e.Node.Name] += e.Size
			repoTotal += e.Size
		}

		if repoTotal > 0 {
			repoStats[repo] = repoTotal
		}
	}

	for _, r := range repos {
		process(username, r)
	}
	for _, s := range subs {
		process(s[0], s[1])
	}

	return langStats, repoStats
}

func fetchRepoActivity(client *http.Client, token string, repos []string, subs [][2]string) []RepoActivity {
	since := time.Now().AddDate(0, 0, -30).Format(time.RFC3339)
	activities := make(map[string]int)

	process := func(owner, repo string) {
		if owner == username && excludedRepos[repo] {
			return
		}

		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?since=%s&author=%s&per_page=100",
			owner, repo, since, username)

		body, err := doRESTRequest(client, token, url)
		if err != nil {
			return
		}

		var commits []map[string]any
		json.Unmarshal(body, &commits)

		if len(commits) > 0 {
			activities[repo] = len(commits)
		}
	}

	for _, r := range repos {
		process(username, r)
	}
	for _, s := range subs {
		process(s[0], s[1])
	}

	var result []RepoActivity
	for name, commits := range activities {
		result = append(result, RepoActivity{name, commits})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Commits > result[j].Commits
	})

	return result
}

func generateREADME(langStats map[string]int, repoActivity []RepoActivity) {
	var sb strings.Builder

	sb.WriteString("```bash\n")
	sb.WriteString("$ whoami\n")
	sb.WriteString("> Tien Hai Cao\n\n")

	sb.WriteString("$ field\n")
	sb.WriteString("> Backend Engineering, Distributed Systems, Cloud Infrastructure\n\n")

	sb.WriteString("$ languages --list\n")
	sb.WriteString("> Go, Java, Python, JavaScript, TypeScript\n\n")

	sb.WriteString("$ frameworks --list\n")
	sb.WriteString("> ExpressJS, NestJS, ElysiaJS, Spring Boot, Django, FastAPI, Gin, Fiber\n\n")

	sb.WriteString("$ databases --list\n")
	sb.WriteString("> MongoDB, PostgreSQL, MySQL, Microsoft SQL Server, Redis\n\n")

	sb.WriteString("$ tools --list\n")
	sb.WriteString("> RabbitMQ, Kafka, Docker, Jupyter Notebook, Ubuntu, AWS, Google Cloud\n\n")

	totalCommits := 0
	for _, r := range repoActivity {
		totalCommits += r.Commits
	}

	totalLang := 0
	for _, size := range langStats {
		totalLang += size
	}

	type lang struct {
		name string
		pct  float64
		size int
	}
	var langs []lang
	for n, s := range langStats {
		if s > 0 {
			langs = append(langs, lang{n, float64(s) * 100 / float64(totalLang), s})
		}
	}
	sort.Slice(langs, func(i, j int) bool { return langs[i].pct > langs[j].pct })

	sb.WriteString("$ languages --top5\n")
	maxLangs := min(len(langs), 5)

	for i := range maxLangs {
		l := langs[i]
		sizeKB := float64(l.size) / 1024
		var sizeStr string
		if sizeKB < 1024 {
			sizeStr = fmt.Sprintf("%.1f KB", sizeKB)
		} else {
			sizeStr = fmt.Sprintf("%.1f MB", sizeKB/1024)
		}

		bars := min(int(l.pct / 4), 25)
		barStr := strings.Repeat("█", bars) + strings.Repeat("░", 25-bars)

		prefix := " "
		if i == 0 {
			prefix = ">"
		}

		sb.WriteString(fmt.Sprintf("%s %-20s %-12s %s %.2f %%\n", prefix, l.name, sizeStr, barStr, l.pct))
	}
	sb.WriteString("\n")

	sb.WriteString("$ projects --top10\n")
	maxRepos := min(len(repoActivity), 10)

	if maxRepos == 0 {
		sb.WriteString("> No commits in last 30 days\n")
	} else {
		for i := range maxRepos {
			r := repoActivity[i]
			pct := float64(r.Commits) / float64(totalCommits) * 100

			bars := int(float64(r.Commits) / float64(totalCommits) * 25)
			if bars == 0 && r.Commits > 0 {
				bars = 1
			}
			barStr := strings.Repeat("█", bars) + strings.Repeat("░", 25-bars)

			commitStr := fmt.Sprintf("%d commits", r.Commits)

			displayName := r.Name
			if len(displayName) > 15 {
				displayName = displayName[:15] + "..."
			}

			prefix := " "
			if i == 0 {
				prefix = ">"
			}

			sb.WriteString(fmt.Sprintf("%s %-20s %-12s %s %.2f %%\n", prefix, displayName, commitStr, barStr, pct))
		}
	}

	sb.WriteString("```\n")

	readmeContent, err := os.ReadFile("README.md")
	if err != nil {
		log.Fatal("Lỗi khi đọc file README.md:", err)
	}

	startMarker := "<!--START_SECTION:dashboard-->"
	endMarker := "<!--END_SECTION:dashboard-->"
	newContent := sb.String()

	readmeStr := string(readmeContent)
	startIndex := strings.Index(readmeStr, startMarker)
	endIndex := strings.Index(readmeStr, endMarker)

	if startIndex == -1 || endIndex == -1 || startIndex >= endIndex {
		log.Fatal("Không tìm thấy hoặc marker không hợp lệ trong README.md")
	}

	updatedContent := readmeStr[:startIndex+len(startMarker)] + "\n" + newContent + readmeStr[endIndex:]

	err = os.WriteFile("README.md", []byte(updatedContent), 0644)
	if err != nil {
		log.Fatal("Lỗi khi ghi file README.md:", err)
	}
}
