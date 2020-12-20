https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-read-table-item.html

# Create DynamoDB Table
- Name: Project-NoteTaker-Test
- Partition Key: NoteID (String)
- Secondary Index (GSI):
  - UserID (String)
  - LastModifiedDate (String)
- Primary:
  - RCU=5
  - WCU=5
- Secondary
  - RCU=5
  - WCU=1
- Other Attributes
  - CreationDate (String)
  - NoteID (String)
  - Content (String)
  - Views (Number)