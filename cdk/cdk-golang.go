package cdk

import (
	configuration "github.com/tolubydesign/angular-story-backend/app/config"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	// "apigtw-lambda-ddb/config"
)

// resource - https://harshq.medium.com/building-apps-with-aws-sdk-for-golang-api-gateway-and-lambda-b254858b1d71
// 					-- https://gist.github.com/harshq/3c821984a96782d5aea329f132a384cb#file-test-api-go
// resource - https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-dynamo-db.html#http-api-dynamo-db-create-table
// resource - https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk/awsapigateway#section-readme
// resource - https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk/awsdynamodb
// (!) resource - https://github.com/aws-samples/serverless-patterns/tree/main/apigw-lambda-dynamodb-cdk-go
// resource - https://serverlessland.com/patterns/apigw-lambda-dynamodb-cdk-go
// (!) resource - https://github.com/cowcoa/aws-cdk-go-examples/blob/master/serverless/apigateway_lambda_ddb/cdk_main.go
// (!) resource - https://github.com/cowcoa/aws-cdk-go-examples/tree/master/serverless/apigateway_lambda_ddb
// (!) resource - https://github.com/abhirockzz/dynamodb-streams-lambda-golang
// resource - https://itnext.io/learn-how-to-use-dynamodb-streams-with-aws-lambda-and-go-f7abcee4d987
// resource - https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-dynamo-db.html#http-api-dynamo-db-create-table
// resource - https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdk/awsapigateway#section-readme



type CdkGolangStackProps struct {
	awscdk.StackProps
}
type ApiGtwLambdaDdbStackProps struct {
	awscdk.StackProps
}

func NewCdkGolangStack(scope constructs.Construct, id string, props *CdkGolangStackProps) awscdk.Stack {
	var sprops awscdk.StackProps

	c, err := configuration.GetConfiguration()
	storyDynamoDBTable := c.Configuration.Dynamodb.StoryTableName
	// userDynamoDBTable := c.Configuration.Dynamodb.UserTableName
	if err != nil {
		panic(err)
	}

	// DynamoDB related constants
	const (
		storyTableName          = "story"
		partitionKeyName        = "email"
		dynamoDBTableNameEnvVar = "DYNAMODB_TABLE_NAME"
	)

	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Create DynamoDB Base table.
	// Data Modeling
	// name(PK), time(SK),                  comment, chat_room
	// string    string(micro sec unixtime)	string   string
	// chatTable := awsdynamodb.NewTable(stack, jsii.String(storyDynamoDBTable), &awsdynamodb.TableProps{
	// 	TableName:     jsii.String(*stack.StackName() + "-Chat"),
	// 	BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
	// 	ReadCapacity:  jsii.Number(1),
	// 	WriteCapacity: jsii.Number(1),
	// 	RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	// 	PartitionKey: &awsdynamodb.Attribute{
	// 		Name: jsii.String("name"),
	// 		Type: awsdynamodb.AttributeType_STRING,
	// 	},
	// 	SortKey: &awsdynamodb.Attribute{
	// 		Name: jsii.String("time"),
	// 		Type: awsdynamodb.AttributeType_STRING,
	// 	},
	// 	PointInTimeRecovery: jsii.Bool(true),
	// })

	// DynamoDB table
	storyTable := awsdynamodb.NewTable(stack, jsii.String("dynamodb-table"), &awsdynamodb.TableProps{
		TableName:     jsii.String(storyDynamoDBTable),
		BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
		ReadCapacity:  jsii.Number(1),
		WriteCapacity: jsii.Number(1),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("time"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	// Create role for lambda function.
	lambdaRole := awsiam.NewRole(stack, jsii.String("LambdaRole"), &awsiam.RoleProps{
		RoleName:  jsii.String(*stack.StackName() + "-LambdaRole"),
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonDynamoDBFullAccess")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("CloudWatchFullAccess")),
		},
	})

	// Create put-chat-records function.
	putFunction := awslambda.NewFunction(stack, jsii.String("PutFunction"), &awslambda.FunctionProps{
		FunctionName: jsii.String(*stack.StackName() + "-PutChatRecords"),
		Runtime:      awslambda.Runtime_GO_1_X(),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(60)),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("functions/put-chat-records/."), nil),
		Handler:      jsii.String("put-chat-records"),
		Architecture: awslambda.Architecture_X86_64(),
		Role:         lambdaRole,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
		Environment: &map[string]*string{
			"DYNAMODB_TABLE": jsii.String(storyDynamoDBTable),
		},
	})

	// Create get-chat-records function.
	getFunction := awslambda.NewFunction(stack, jsii.String("GetChatRecords"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_GO_1_X(),
		FunctionName: jsii.String(*stack.StackName() + "-GetChatRecords"),
		MemorySize:   jsii.Number(128),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(60)),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("functions/get-chat-records/."), nil),
		Handler:      jsii.String("get-chat-records"),
		Architecture: awslambda.Architecture_X86_64(),
		Role:         lambdaRole,
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
		Environment: &map[string]*string{
			"DYNAMODB_TABLE": jsii.String(storyDynamoDBTable),
			"DYNAMODB_GSI":   jsii.String(storyDynamoDBTable + "-SecondaryIndex"),
		},
		// ReservedConcurrentExecutions: jsii.Number(1),
	})

	// Create API Gateway rest api.
	restApi := awsapigateway.NewRestApi(stack, jsii.String("LambdaRestApi"), &awsapigateway.RestApiProps{
		RestApiName:        jsii.String(*stack.StackName() + "-LambdaRestApi"),
		RetainDeployments:  jsii.Bool(false),
		EndpointExportName: jsii.String("RestApiUrl"),
		Deploy:             jsii.Bool(true),
		DeployOptions: &awsapigateway.StageOptions{
			StageName:           jsii.String("dev"),
			CacheClusterEnabled: jsii.Bool(true),
			CacheClusterSize:    jsii.String("0.5"),
			CacheTtl:            awscdk.Duration_Minutes(jsii.Number(1)),
			// https://www.petefreitag.com/item/853.cfm
			// This can help you better understand what burst and rate limite are.
			ThrottlingBurstLimit: jsii.Number(100),
			ThrottlingRateLimit:  jsii.Number(1000),
		},
	})

	// Add path resources to rest api.
	// You MUST associate ApiKey with the methods for the UsagePlane to work.
	putRecordsRes := restApi.Root().AddResource(jsii.String("put-chat-records"), nil)
	putRecordsRes.AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(putFunction, nil), &awsapigateway.MethodOptions{
		ApiKeyRequired: jsii.Bool(true),
	})
	getRecordsRes := restApi.Root().AddResource(jsii.String("get-chat-records"), nil)
	getMethod := getRecordsRes.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(getFunction, nil), &awsapigateway.MethodOptions{
		ApiKeyRequired: jsii.Bool(true),
	})

	// UsagePlane's throttle can override Stage's DefaultMethodThrottle,
	// while UsagePlanePerApiStage's throttle can override UsagePlane's throttle.
	usagePlane := restApi.AddUsagePlan(jsii.String("UsagePlane"), &awsapigateway.UsagePlanProps{
		Name: jsii.String(*stack.StackName() + "-UsagePlane"),
		Throttle: &awsapigateway.ThrottleSettings{
			BurstLimit: jsii.Number(10),
			RateLimit:  jsii.Number(100),
		},
		Quota: &awsapigateway.QuotaSettings{
			Limit:  jsii.Number(100),
			Offset: jsii.Number(0),
			Period: awsapigateway.Period_DAY,
		},
		ApiStages: &[]*awsapigateway.UsagePlanPerApiStage{
			{
				Api:   restApi,
				Stage: restApi.DeploymentStage(),
				Throttle: &[]*awsapigateway.ThrottlingPerMethod{
					{
						Method: getMethod,
						Throttle: &awsapigateway.ThrottleSettings{
							BurstLimit: jsii.Number(1),
							RateLimit:  jsii.Number(10),
						},
					},
				},
			},
		},
	})

	// Create ApiKey and associate it with UsagePlane.
	apiKey := restApi.AddApiKey(jsii.String("ApiKey"), &awsapigateway.ApiKeyOptions{})
	usagePlane.AddApiKey(apiKey, &awsapigateway.AddApiKeyOptions{})

	// Create DynamoDB GSI table.
	// Data Modeling

	// chat_room(PK), time(SK),                  comment, name
	// string         string(micro sec unixtime) string   string
	storyTable.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName: jsii.String(storyDynamoDBTable + "-SecondaryIndex"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("creator"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		ProjectionType: awsdynamodb.ProjectionType_ALL,
	})

	// Grant access to lambda functions.
	storyTable.GrantWriteData(putFunction)
	storyTable.GrantReadData(getFunction)

	return stack

	// var sprops awscdk.StackProps
	// if props != nil {
	// 	sprops = props.StackProps
	// }
	// stack := awscdk.NewStack(scope, &id, &sprops)
	// The code that defines your stack goes here
	// example resource
	// queue := awssqs.NewQueue(stack, jsii.String("CdkGolangQueue"), &awssqs.QueueProps{
	// 	VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
	// })
	// return stack
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	c, err := configuration.GetConfiguration()
	if err != nil {
		panic(err)
	}

	account := c.Configuration.AWS.AccountID
	region := c.Configuration.AWS.Region
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil
	return &awscdk.Environment{
		Account: jsii.String(account),
		Region:  jsii.String(region),
	}

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}

func RunCDK() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCdkGolangStack(app, "CdkGolangStack", &CdkGolangStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}
