package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var channelURL string

func init() {
	channelURL = os.Getenv("GIT_COMMIT_CHANNEL")
	if len(channelURL) == 0 {
		panic("GIT_COMMIT_CHANNEL is empty")
	}
}

type attachments struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Color string `json:"color"`
}

type slackPostBody struct {
	Text        string        `json:"text"`
	Username    string        `json:"username"`
	IconURL     string        `json:"icon_url"`
	Attachments []attachments `json:"attachments"`
}

type gitArgs struct {
	Subject     string
	Hash        string
	Author      string
	CommitDate  string
	Branch      string
	RootDirName string
}

func (args gitArgs) format() string {
	return fmt.Sprintf(
		"Branch: %s\nAuthor: %s\nDate: %s\nHash: %s\nRepository: %s",
		args.Branch, args.Author, args.CommitDate, args.Hash, args.RootDirName)
}

func makeGitArgs(str string) gitArgs {
	splitedString := strings.Split(str, ",")
	return gitArgs{
		Subject:     splitedString[0],
		Hash:        splitedString[1],
		Author:      splitedString[2],
		CommitDate:  splitedString[3],
		Branch:      splitedString[4],
		RootDirName: splitedString[5],
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	gitArgs := makeGitArgs(args[0])
	err := slackPost(gitArgsToSlackPostSimpleBody(gitArgs))
	if err != nil {
		panic(err)
	}
}

var defaultCommitBotIconURL = "https://i1.wp.com/pbs.twimg.com/profile_images/602729491916435458/hSu0UjMC_400x400.jpg?resize=300%2C300&ssl=1"

func gitArgsToSlackPostSimpleBody(gitArgs gitArgs) slackPostBody {
	botIcon := os.Getenv("GIT_COMMIT_BOT_ICON")

	if len(botIcon) == 0 {
		botIcon = defaultCommitBotIconURL
	}

	titleAttachment := attachments{Title: gitArgs.Subject, Text: gitArgs.format(), Color: "good"}
	attachments := []attachments{titleAttachment}
	return slackPostBody{
		Text:        "",
		Username:    "コミットログ",
		IconURL:     botIcon,
		Attachments: attachments,
	}
}

func slackPost(body slackPostBody) error {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return err
	}
	req, err := http.NewRequest(
		"POST",
		channelURL,
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	println(resp)
	defer resp.Body.Close()

	return err
}
