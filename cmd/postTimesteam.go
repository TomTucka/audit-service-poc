package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
)

type Session struct {
	AwsSession *session.Session
}

type SiriusEvent struct {
	userID 				int
	personID 			int				//Client ID or Deputy ID
	eventType			string			//Insert Update
	eventClass			string			//Opg\Core\Model\Event\Common\TaskCreated
	sourceEntityID 		int				//144
	sourceEntityClass 	string			//Opg\Core\Model\Entity\Task\Task
}

func newSiriusEvent(userID int, personID int, eventType string, eventClass string, sourceEntityID int, sourceEntityClass string) *SiriusEvent {
	e := SiriusEvent{userID: userID, personID: personID, eventType: eventType, eventClass: eventClass, sourceEntityID: sourceEntityID, sourceEntityClass: sourceEntityClass}
	return &e
}

func PostTimestream() {

	sess, err := session.NewSession()
	if err != nil {
		log.Fatalln(err)
	}

	RoleArn := "arn:aws:iam::288342028542:role/operator"

	creds := stscreds.NewCredentials(sess, RoleArn)
	awsConfig := aws.Config{Credentials: creds, Region: aws.String("eu-west-1")}

	se:= newSiriusEvent(100, 1,"INSERT", "Opg\\Core\\Model\\Event\\Common\\TaskCreated", 144, "Opg\\Core\\Model\\Entity\\Task\\Task")

	svc := timestreamwrite.New(sess, &awsConfig)

	now := time.Now()
	currentTimeInSeconds := now.Unix()

	records := &timestreamwrite.WriteRecordsInput{
		DatabaseName: aws.String("audit-service-poc"),
		TableName: aws.String("audit-service-poc"),
		Records: []*timestreamwrite.Record{
			&timestreamwrite.Record{
				Dimensions: []*timestreamwrite.Dimension{
					&timestreamwrite.Dimension{
						Name:  aws.String("User ID"),
						Value: aws.String(string(rune(se.userID))),
					},
					&timestreamwrite.Dimension{
						Name:  aws.String("Person ID"),
						Value: aws.String(string(rune(se.personID))),
					},
					&timestreamwrite.Dimension{
						Name:  aws.String("Event Type"),
						Value: aws.String(se.eventType),
					},
					&timestreamwrite.Dimension{
						Name:  aws.String("Event Class"),
						Value: aws.String(se.eventClass),
					},
					&timestreamwrite.Dimension{
						Name:  aws.String("Source Entity ID"),
						Value: aws.String(string(rune(se.sourceEntityID))),
					},
					&timestreamwrite.Dimension{
						Name:  aws.String("Source Entity Class"),
						Value: aws.String(se.sourceEntityClass),
					},
				},
				Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
				TimeUnit:         aws.String("SECONDS"),
			},
		},
	}

	_, err = svc.WriteRecords(records)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		fmt.Println("Write records is successful")
	}
}
