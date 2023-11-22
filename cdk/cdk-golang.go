package cdk

import (
	configuration "github.com/tolubydesign/angular-story-backend/app/config"
	
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	// "github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/awslambdago"
	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkGolangStackProps struct {
	awscdk.StackProps
}

/**
 * This file showcases how to split up a RestApi's Resources and Methods across nested stacks.
 *
 * The root stack 'RootStack' first defines a RestApi.
 * Two nested stacks BooksStack and PetsStack, create corresponding Resources '/books' and '/pets'.
 * They are then deployed to a 'prod' Stage via a third nested stack - DeployStack.
 *
 * To verify this worked, go to the APIGateway
 */
type rootStack struct {
	stack
}

func NewCdkGolangStack(scope constructs.Construct, id string, props *CdkGolangStackProps) awscdk.Stack {
	this := &rootStack{}
	var backend function
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here
	// Create a new api HTTP api on gateway v2.
	api := awsapigateway.NewRestApi(stack, jsii.String("cdk-lambda-api"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("cdk-lambda-api"),
		
	})
	// api = awsapigateway.NewHttpApi(stack, jsii.String("cdk-lambda-api"), &awsapigatewayv2.HttpApiProps{
	// 	CorsPreflight: &awsapigateway.CorsPreflightOptions{
	// 		AllowOrigins: &[]*string{jsii.String("*")}, //
	// 		AllowMethods: &[]awsapigateway.CorsHttpMethod{awsapigatewayv2.CorsHttpMethod_ANY},
	// 	},
	// })

	// The code that defines your stack goes here
	// example resource
	// queue := awssqs.NewQueue(stack, jsii.String("CdkGolangQueue"), &awssqs.QueueProps{
	// 	VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
	// })

	return stack
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
