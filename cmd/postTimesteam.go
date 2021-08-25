package cmd

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
)

type Session struct {
	AwsSession *session.Session
}

func PostTimestream() {

	sess, err := session.NewSession()
	if err != nil {
		log.Fatalln(err)
	}

	RoleArn := "arn:aws:iam::288342028542:role/operator"

	creds := stscreds.NewCredentials(sess, RoleArn)
	awsConfig := aws.Config{Credentials: creds, Region: aws.String("eu-west-1")}

	svc := timestreamwrite.New(sess, &awsConfig)

	records := &timestreamwrite.WriteRecordsInput{
		DatabaseName: aws.String("audit-service-poc"),
		TableName: aws.String("audit-service-poc"),
		Records: &timestreamwrite.Record{

		}
	}

	svc.WriteRecords(records)
}
