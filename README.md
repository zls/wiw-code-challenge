# Running locally
> this project uses ports tcp/8001 and tcp/8181

> interacting with dynamodb via cli requires `--endpoint-url` to be set as well as an aws profile or `AWS_*` keys to be set. 

1. docker-compose up
1. Run dynamodb table creation scripts [create ddb table](#create-ddb-table) and [update gsi](#update-gsi)


### Create ddb table
aws dynamodb create-table --endpoint-url http://localhost:8001 \
  --table-name shifts \
  --attribute-definitions \
    AttributeName=UserID,AttributeType=N \
    AttributeName=StartTime,AttributeType=S \
  --key-schema \
    AttributeName=UserID,KeyType=HASH \
    AttributeName=StartTime,KeyType=RANGE \
  --provisioned-throughput=ReadCapacityUnits=1,WriteCapacityUnits=1

#### Update GSI
aws dynamodb update-table --endpoint-url http://localhost:8001 \
  --table-name shifts \
  --attribute-definitions \
    AttributeName=ID,AttributeType=B \
  --global-secondary-index-updates '[ { "Create": {"IndexName": "ByID", "KeySchema": [{"AttributeName": "ID", "KeyType": "HASH"} ], "Projection": { "ProjectionType": "ALL" }, "ProvisionedThroughput": { "ReadCapacityUnits": 1, "WriteCapacityUnits": 1 } } } ]'

aws dynamodb update-table --endpoint-url http://localhost:8001 \
  --table-name shifts \
  --attribute-definitions \
    AttributeName=AccountID,AttributeType=N \
  --global-secondary-index-updates '[ { "Create": {"IndexName": "ByAccount", "KeySchema": [{"AttributeName": "AccountID", "KeyType": "HASH"}, {"AttributeName": "StartTime", "KeyType": "RANGE"} ], "Projection": { "ProjectionType": "ALL" }, "ProvisionedThroughput": { "ReadCapacityUnits": 1, "WriteCapacityUnits": 1 } } } ]'

#### List ddb tables
aws dynamodb list-tables --endpoint-url http://localhost:8001

### Describe ddb table
aws dynamodb describe-table --table shifts --endpoint-url http://localhost:8001

### TODOS

* Figure out why delete isn't working
* Add an update shift endpoint
* better logging. a few different log levels (debug, info and error) and formatting are a necessity. Logging session info too.
* better handling of aws ddb errors, currently don't distinguish between client or server errors.
* ddb framework? dynamodb is a bit cumbersome.
* rework dynamodb schema. `ByID` index could be a local secondary index. Though, it might make more sense to structure primary partition key off `AccountID` as this is probably the more common access pattern (`ByID` can be update accordingly)
* return client errors. Currently the backend just returns a 404 and empty json.
* timezones are not normalized, this interferes with natural sort and comparison. As long as a timezone is consistent across shifts for a user this isn't a problem. DynamoDB, not having a time type might not have been the best choice unless you can make the restriction to only allow shifts to a specific timezone, possibly through a location. Per `rfc3339`
>>> 5.1. Ordering

   If date and time components are ordered from least precise to most
   precise, then a useful property is achieved.  Assuming that the time
   zones of the dates and times are the same (e.g., all in UTC),
   expressed using the same string (e.g., all "Z" or all "+00:00"), and
   all times have the same number of fractional second digits, then the
   date and time strings may be sorted as strings (e.g., using the
   strcmp() function in C) and a time-ordered sequence will result.  The
   presence of optional punctuation would violate this characteristic.