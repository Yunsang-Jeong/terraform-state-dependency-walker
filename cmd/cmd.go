package cmd

type stringFlag struct {
	shorten      string
	defaultValue string
	description  string
	requirement  bool
}

const (
	AWSClientRegion          = "aws-region"
	TerraformStateBucketName = "bucket-name"
	TerraformStateObjectName = "state-name"
	DynamodbTableName        = "ddb-name"
	DynamodbTablePrimaryKey  = "ddb-pk"
)
