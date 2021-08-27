package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
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

func newSiriusEvent(eventType string, sourceEntityClass string) *SiriusEvent {
	userID := rand.Intn(10)
	personID := rand.Intn(100)
	sourceEntityID := rand.Intn(200)

	eventClass := ""
	i := rand.Intn(4)
	switch i {
	case 1:
		eventClass = "Opg\\Core\\Model\\Event\\Common\\TaskCreated"
	case 2:
		eventClass = "Opg\\Core\\Model\\Event\\Task\\TaskEdited"
	case 3:
		eventClass = "Opg\\Core\\Model\\Event\\Task\\TaskReassigned"
	default:
		eventClass = "Opg\\Core\\Model\\Event\\Common\\TaskCompleted"
	}

	e := SiriusEvent{userID: userID, personID: personID, eventType: eventType, eventClass: eventClass, sourceEntityID: sourceEntityID, sourceEntityClass: sourceEntityClass}
	return &e
}

func runQuery(queryPtr *string, querySvc *timestreamquery.TimestreamQuery) {
	queryInput := &timestreamquery.QueryInput{
		QueryString: aws.String(*queryPtr),
	}
	// execute the query
	err := querySvc.QueryPages(queryInput,
		func(page *timestreamquery.QueryOutput, lastPage bool) bool {
			// process query response
			queryStatus := page.QueryStatus
			fmt.Println("Current query status:", queryStatus) //what happens when error?
			// query response metadata
			// includes column names and types
			metadata := page.ColumnInfo
			fmt.Println("Metadata:")
			fmt.Println(metadata)

			// query response data
			fmt.Println("All rows then Data:")
			//fmt.Println(page.Rows)
			// process rows
			rows := page.Rows
			for i := 0; i < len(rows); i++ {
				data := rows[i].Data
				value := processRowType(data, metadata)
				fmt.Println(value)
			}
			fmt.Println("Number of rows:", len(page.Rows))
			return true
		})
	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	}
}

func processRowType(data []*timestreamquery.Datum, metadata []*timestreamquery.ColumnInfo) string {
	value := ""
	for j := 0; j < len(data); j++ {
		if metadata[j].Type.ScalarType != nil {
			// process simple data types
			value += processScalarType(data[j])
		} else if metadata[j].Type.TimeSeriesMeasureValueColumnInfo != nil {
			// fmt.Println("Timeseries measure value column info")
			// fmt.Println(metadata[j].Type.TimeSeriesMeasureValueColumnInfo.Type)
			datapointList := data[j].TimeSeriesValue
			value += "["
			value += processTimeSeriesType(datapointList, metadata[j].Type.TimeSeriesMeasureValueColumnInfo)
			value += "]"
		} else if metadata[j].Type.RowColumnInfo != nil {
			columnInfo := metadata[j].Type.RowColumnInfo
			datumList := data[j].RowValue.Data
			value += "["
			value += processRowType(datumList, columnInfo)
			value += "]"
		} else {
			panic("Bad column type")
		}
		// comma seperated column values
		if j != len(data)-1 {
			value += ", "
		}
	}
	return value
}

func processScalarType(data *timestreamquery.Datum) string {
	return *data.ScalarValue
}

func processTimeSeriesType(data []*timestreamquery.TimeSeriesDataPoint, columnInfo *timestreamquery.ColumnInfo) string {
	value := ""
	for k := 0; k < len(data); k++ {
		time := data[k].Time
		value += *time + ":"
		if columnInfo.Type.ScalarType != nil {
			value += processScalarType(data[k].Value)
		} else {
			panic("Bad data type")
		}
		if k != len(data)-1 {
			value += ", "
		}
	}
	return value
}

func GetTimestream() {

	sess, err := session.NewSession()
	if err != nil {
		log.Fatalln(err)
	}

	RoleArn := "arn:aws:iam::288342028542:role/operator"

	creds := stscreds.NewCredentials(sess, RoleArn)
	awsConfig := aws.Config{Credentials: creds, Region: aws.String("eu-west-1")}

	querySvc := timestreamquery.New(sess, &awsConfig)

	var queryPtr *string
	if queryPtr == nil {
		query := "SELECT User_ID, Person_ID, Source_Entity_ID, measure_value::VARCHAR, time\nFROM \"audit-service-poc\".\"audit-service-poc\" \nWHERE Person_ID = '20'\nORDER BY time DESC\nLIMIT 10"
		queryPtr = &query
	}

	runQuery(queryPtr,querySvc)
}

func PostTimestream() {

	sess, err := session.NewSession()
	if err != nil {
		log.Fatalln(err)
	}

	RoleArn := "arn:aws:iam::288342028542:role/operator"

	creds := stscreds.NewCredentials(sess, RoleArn)
	awsConfig := aws.Config{Credentials: creds, Region: aws.String("eu-west-1")}

	se:= newSiriusEvent("INSERT", "Opg\\Core\\Model\\Entity\\Task\\Task")

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
						Name:  aws.String("User_ID"),
						Value: aws.String(strconv.Itoa(se.userID)),
					},
					&timestreamwrite.Dimension{
						Name:  aws.String("Person_ID"),
						Value: aws.String(strconv.Itoa(se.personID)),
					},
					&timestreamwrite.Dimension{
						Name:  aws.String("Event_Type"),
						Value: aws.String(se.eventType),
					},
					&timestreamwrite.Dimension{
						Name:  aws.String("Source_Entity_ID"),
						Value: aws.String(strconv.Itoa(se.sourceEntityID)),
					},
					&timestreamwrite.Dimension{
						Name:  aws.String("Source_Entity_Class"),
						Value: aws.String(se.sourceEntityClass),
					},
					&timestreamwrite.Dimension{
						Name:  aws.String("Event_Class"),
						Value: aws.String(se.eventClass),
					},
				},
				MeasureName:      aws.String("Event"),
				MeasureValue:     aws.String(se.eventClass),
				MeasureValueType: aws.String("VARCHAR"),
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
