package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

func GetTimestream(p string) *timestreamquery.QueryOutput {

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
		fmt.Println(p)
		query := "SELECT User_ID, Person_ID, Source_Entity_ID, measure_value::VARCHAR, time\nFROM \"audit-service-poc\".\"audit-service-poc\" \nWHERE Person_ID = '" + p + "'\nORDER BY time DESC\nLIMIT 10"
		queryPtr = &query
	}

	response := runQuery(queryPtr, querySvc)
	return response
}

func processRowType(data []*timestreamquery.Datum, metadata []*timestreamquery.ColumnInfo) string {
	value := ""

	for j := 0; j < len(data); j++ {
		if metadata[j].Type.ScalarType != nil {
			// process simple data types
			value += processScalarType(data[j])
		} else if metadata[j].Type.TimeSeriesMeasureValueColumnInfo != nil {
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
	fmt.Println("time series type")


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


func runQuery(queryPtr *string, querySvc *timestreamquery.TimestreamQuery) *timestreamquery.QueryOutput {
	queryInput := &timestreamquery.QueryInput{
		QueryString: aws.String(*queryPtr),
	}
	// execute the query
	response, _ := querySvc.Query(queryInput)
	fmt.Print(response)
	metadata := response.ColumnInfo
	for i := 0; i < len(metadata); i++ {
		name := metadata[i].Name
		fmt.Print(*name + " | ")
	}
	rows := response.Rows
	for i := 0; i < len(rows); i++ {
		data := rows[i].Data
		value := processRowType(data, metadata)
		fmt.Println(value)
	}
	fmt.Println("Number of rows:", len(response.Rows))
	return response

	//err := querySvc.QueryPages(queryInput,
	//	func(page *timestreamquery.QueryOutput, lastPage bool) bool {
	//		// process query response
	//		queryStatus := page.QueryStatus
	//		fmt.Println("Current query status:", queryStatus) //what happens when error?
	//		// query response metadata
	//		// includes column names and types
	//
	//
	//		metadata := page.ColumnInfo
	//		for i := 0; i < len(metadata); i++ {
	//			name := metadata[i].Name
	//			//scalar := metadata[i].Type.ScalarType
	//			//fmt.Print(*name + "," + *scalar + " | ")
	//			fmt.Print(*name + " | ")
	//		}
	//		fmt.Println("")
	//
	//		// process rows
	//		rows := page.Rows
	//		for i := 0; i < len(rows); i++ {
	//			data := rows[i].Data
	//			value := processRowType(data, metadata)
	//			fmt.Println(value)
	//		}
	//		fmt.Println("Number of rows:", len(page.Rows))
	//		return true
	//	})
	//if err != nil {
	//	fmt.Println("Error:")
	//	fmt.Println(err)
	//}
}

