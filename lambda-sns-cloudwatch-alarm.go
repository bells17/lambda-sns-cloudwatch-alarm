package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/bluele/slack"
)

var (
	SlackWebhookUrl string = os.Getenv("SlackWebhookUrl")
	SlackChannel    string = os.Getenv("SlackChannel")
	SlackUserName   string = os.Getenv("SlackUserName")
	SlackIconEmoji  string = os.Getenv("SlackIconEmoji")
	Debug           string = os.Getenv("Debug")
)

type CloudWatchAlermMessage struct {
	AWSAccountId     string `json:"AWSAccountId"`
	AlarmDescription string `json:"AlarmDescription"`
	AlarmName        string `json:"AlarmName"`
	NewStateReason   string `json:"NewStateReason"`
	NewStateValue    string `json:"NewStateValue"`
	OldStateValue    string `json:"OldStateValue"`
	Region           string `json:"Region"`
	StateChangeTime  string `json:"StateChangeTime"`
	// Trigger
}

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, snsEvent events.SNSEvent) {
	var err error
	if Debug != "" {
		eventJson, err := json.Marshal(snsEvent)
		if err != nil {
			panic(err)
		}
		log.Printf("[event start] %s\n", eventJson)
	}

	for _, record := range snsEvent.Records {
		message := &CloudWatchAlermMessage{}
		err = json.Unmarshal(([]byte)(record.SNS.Message), message)
		if err != nil {
			panic(err)
		}
		err = sendToSlack(message)
		if err != nil {
			panic(err)
		}
	}
}

func sendToSlack(message *CloudWatchAlermMessage) error {
	slackWebhookUrl, err := getSlackWebhookUrl()
	if err != nil {
		return err
	}
	webhook := slack.NewWebHook(slackWebhookUrl)
	return webhook.PostMessage(buildPayload(message))
}

func buildPayload(message *CloudWatchAlermMessage) *slack.WebHookPostPayload {
	var color string
	var statusEmojiIcon string

	payload := &slack.WebHookPostPayload{}
	payload.Channel = SlackChannel
	payload.Username = SlackUserName
	payload.IconEmoji = SlackIconEmoji

	if message.NewStateValue == "OK" {
		statusEmojiIcon = ":ok_hand:"
		color = "good"
	} else {
		statusEmojiIcon = ":exclamation:"
		color = "danger"
	}

	attachment := &slack.Attachment{}
	attachment.Color = color
	attachment.Pretext = fmt.Sprintf("%s %s",
		statusEmojiIcon,
		message.NewStateReason)
	attachment.Title = fmt.Sprintf("%s", message.AlarmName)

	var fields []*slack.AttachmentField
	fields = append(fields, &slack.AttachmentField{
		Title: "State",
		Value: message.NewStateValue,
		Short: true,
	}, &slack.AttachmentField{
		Title: "Description",
		Value: message.AlarmDescription,
		Short: true,
	}, &slack.AttachmentField{
		Title: "StateChangeTime",
		Value: message.StateChangeTime,
		Short: true,
	})

	attachment.Fields = fields
	payload.Attachments = []*slack.Attachment{attachment}
	return payload
}

func getSlackWebhookUrl() (string, error) {
	if strings.Index(SlackWebhookUrl, "https://hooks.slack.com/services/") == 0 {
		return SlackWebhookUrl, nil
	}

	kmsClient := kms.New(session.New())
	decodedBytes, err := base64.StdEncoding.DecodeString(SlackWebhookUrl)
	if err != nil {
		return "", err
	}
	input := &kms.DecryptInput{CiphertextBlob: []byte(decodedBytes)}
	response, err := kmsClient.Decrypt(input)
	if err != nil {
		return "", err
	}
	return string(response.Plaintext[:]), nil
}
