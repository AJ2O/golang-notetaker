package notes

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Note represents a single note that can be interacted with by the user.
type Note struct {
	NoteID           string
	UserID           string
	LastModifiedDate string
	CreationDate     string
	Content          string
	Views            int
}

// NotePage represents a page of notes.
type NotePage struct {
	Notes []Note
}

var awsSession *session.Session
var ddb *dynamodb.DynamoDB

// DDBTable is the DynamoDB table storing the notes data.
var DDBTable string

func init() {
	awsSession = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	ddb = dynamodb.New(awsSession)
	log.Println("DynamoDB session started")
}

func getCurrentTime() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// CreateNote creates a new note to be stored in the user's account.
func CreateNote(userID string, content string) error {
	currentTime := getCurrentTime()

	// construct input struct
	inputNote := struct {
		NoteID           string
		UserID           string
		LastModifiedDate string
		CreationDate     string
		Content          string
		Views            int
	}{
		ksuid.New().String(),
		userID,
		currentTime,
		currentTime,
		content,
		0,
	}
	noteAVMap, err := dynamodbattribute.MarshalMap(inputNote)
	if err != nil {
		panic("Cannot marshal note into AttributeValue map")
	}

	// construct the input params
	params := &dynamodb.PutItemInput{
		Item:      noteAVMap,
		TableName: aws.String(DDBTable),
	}

	// insert the item into DynamoDB
	_, err = ddb.PutItem(params)
	if err != nil {
		return errors.New("Could not add new note")
	}

	return nil
}

// ReadNote returns the stored note with the given ID.
func ReadNote(noteID string) (Note, error) {
	// construct the query parameters
	query := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"NoteID": {
				S: aws.String(noteID),
			},
		},
		TableName: aws.String(DDBTable),
	}

	// query note from DynamoDB, checking for errors
	result, err := ddb.GetItem(query)
	if err != nil {
		return Note{}, err
	}
	if result.Item == nil {
		msg := "A note with ID " + noteID + " does not exist"
		return Note{}, errors.New(msg)
	}

	// unmarshall return date into a Note struct
	note := Note{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &note)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return note, nil
}

// ReadAllNotes returns all the notes stored in the user's account.
func ReadAllNotes(userID string) ([]Note, error) {
	// construct the query parameters
	query := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":u": {
				S: aws.String(userID),
			},
		},
		KeyConditionExpression: aws.String("UserID = :u"),
		ScanIndexForward:       aws.Bool(false),
		TableName:              aws.String(DDBTable),
		IndexName:              aws.String("UserID-LastModifiedDate-index"),
	}

	// get notes from DynamoDB, checking for errors
	result, err := ddb.Query(query)
	if err != nil {
		return []Note{}, err
	}

	// unmarshall return data into a slice of Note structs
	var allNotes []Note
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &allNotes)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Records, %v", err))
	}

	return allNotes, nil
}

// UpdateNote modifies the stored note with newly submitted contents.
func UpdateNote(noteID string, contentUpdate string) error {
	currentTime := getCurrentTime()

	// construct update parameters
	updateParams := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {
				S: aws.String(contentUpdate),
			},
			":l": {
				S: aws.String(currentTime),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"NoteID": {
				S: aws.String(noteID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		TableName:        aws.String(DDBTable),
		UpdateExpression: aws.String("SET Content = :c, LastModifiedDate = :l"),
	}

	// update the database
	_, err := ddb.UpdateItem(updateParams)
	if err != nil {
		return err
	}
	return nil
}

// UpdateNoteView adds +1 to the given note's view count.
func UpdateNoteView(noteID string) error {
	// construct update parameters
	updateParams := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v": {
				N: aws.String("1"),
			},
		},
		// Views is a reserved keyword, so it has it be substituted
		ExpressionAttributeNames: map[string]*string{
			"#views": aws.String("Views"),
		},
		Key: map[string]*dynamodb.AttributeValue{
			"NoteID": {
				S: aws.String(noteID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		TableName:        aws.String(DDBTable),
		UpdateExpression: aws.String("ADD #views :v"),
	}

	// update the database
	_, err := ddb.UpdateItem(updateParams)
	if err != nil {
		return err
	}
	return nil
}

// DeleteNote removes the stored note from the user's account.
func DeleteNote(noteID string) error {

	// construct update parameters
	deleteParams := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"NoteID": {
				S: aws.String(noteID),
			},
		},
		TableName: aws.String(DDBTable),
	}

	// delete from the database
	_, err := ddb.DeleteItem(deleteParams)
	if err != nil {
		return err
	}
	return nil
}
