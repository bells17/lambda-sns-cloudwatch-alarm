# lambda-sns-cloudwatch-alarm

This is an AWS Lambda function, for Send AWS CloudWatch Alarm to Slack.

## Build & Upload a function

```
make dep
make pack
make publish FUNCTION_NAME=lambda-function-name
```

## Lambda setting

- set handle name to `lambda-sns-cloudwatch-alarm`
- set below environment variables

key | value
--- | ---
SlackWebhookUrl | the Slack incoming webhook URL (support a KMS encryption)
SlackChannel | Channel name of send an alarm
SlackUserName | Slack username
SlackIconEmoji | Slack user icon emoji
