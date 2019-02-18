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

type slackPostBody struct {
	Attachments []struct {
		Text string `json:"text"`
	} `json:"attachments"`
}

type slackPostSimpleBody struct {
	Text     string `json:"text"`
	Username string `json:"username"`
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
		"*%s*\nBranch: %s\nAuthor: %s\nDate: %s\nHash: %s\nRepository: %s",
		args.Subject, args.Branch, args.Author, args.CommitDate, args.Hash, args.RootDirName)
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
	println(args)
	gitArgs := makeGitArgs(args[0])
	err := slackPost(gitArgsToSlackPostSimpleBody(gitArgs))
	if err != nil {
		panic(err)
	}
}

func gitArgsToSlackPostSimpleBody(gitArgs gitArgs) slackPostSimpleBody {
	text := gitArgs.format()
	return slackPostSimpleBody{
		Text:     text,
		Username: "コミット",
	}
}

func slackPost(body slackPostSimpleBody) error {
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
