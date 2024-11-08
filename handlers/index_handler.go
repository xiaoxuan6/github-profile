package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olekukonko/tablewriter"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/xiaoxuan6/github-profile/services"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var IndexHandler = new(indexHandler)

type indexHandler struct {
}

func (indexHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "homepage.html", nil)
}

func (indexHandler) Generate(c *gin.Context) {
	if len(os.Getenv("GITHUB_TOKEN")) < 1 {
		log.Fatalln("GITHUB_TOKEN environment variable not set")
	}

	userName, _ := c.GetPostForm("r")
	userName = strings.TrimSpace(userName)
	profiles, prRepository := services.GenerateProfile(userName)

	var profilesData [][]string
	for k, profile := range profiles {
		re := regexp.MustCompile(`\s+`)
		profile.Description = re.ReplaceAllString(profile.Description, "")
		profile.Description = strings.ReplaceAll(profile.Description, "|", "")
		profilesData = append(profilesData, []string{strconv.Itoa(k + 1), profile.FullName(), profile.Description, profile.Create, profile.Update, profile.Language, strconv.Itoa(profile.Stars)})
	}

	profilesHtml := fmt.Sprintf("## Created repos\n")
	profilesHtml += renderTable([]string{"ID", "Repo", "Description", "Crate", "Update", "Language", "Stars"}, profilesData)

	var prRepositoriesData [][]string
	for k, p := range prRepository {
		prRepositoriesData = append(prRepositoriesData, []string{strconv.Itoa(k + 1), p.FullName(), p.State, p.Created, p.CountUrl(userName)})
	}

	prRepositoriesHtml := fmt.Sprintf("## Pr repos \n")
	prRepositoriesHtml += renderTable([]string{"ID", "Repo", "State", "Pr Date", "prCount"}, prRepositoriesData)

	html := github_flavored_markdown.Markdown([]byte(profilesHtml + prRepositoriesHtml))
	c.HTML(http.StatusOK, "homepage.html", gin.H{
		"Title": userName,
		"Body":  template.HTML(string(html)),
	})
}

func renderTable(header []string, data [][]string) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader(header)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data)
	table.Render()
	return tableString.String()
}
