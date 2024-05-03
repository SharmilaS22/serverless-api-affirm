package main

import (
	// "fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
)

type AppStackProps struct {
	awscdk.StackProps
}

func NewAppStack(scope constructs.Construct, id string, props *AppStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	api := awsapigatewayv2.NewHttpApi(stack, jsii.String("aff-API"), &awsapigatewayv2.HttpApiProps{
		CorsPreflight: &awsapigatewayv2.CorsPreflightOptions{
			AllowOrigins: &[]*string{jsii.String("*")},
			AllowMethods: &[]awsapigatewayv2.CorsHttpMethod{awsapigatewayv2.CorsHttpMethod_ANY},
		},
		ApiName: jsii.String("affirmations-api"),
	})

	lambdaRole := awsiam.NewRole(stack, jsii.String("dynamodbRoleForLambda"), &awsiam.RoleProps{
		RoleName: jsii.String("dynamoAccessToLambda"),
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), &awsiam.ServicePrincipalOpts{
			// filter out only lambda from this project
		}),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSLambdaDynamoDBExecutionRole")),
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonDynamoDBFullAccess")),
		},
	})

	getAffFunctionIntegration := awslambda.NewFunction(stack, jsii.String("get-aff-Lambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PYTHON_3_9(),
		Handler: jsii.String("index.get_aff_handler"),
		Code: awslambda.Code_FromAsset(jsii.String("lambda"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("getAffirmation-Function"),
		Role: lambdaRole,
	})

	getlambdaintegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("Lambda"),
		getAffFunctionIntegration,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)

	postAffFunctionIntegration := awslambda.NewFunction(stack, jsii.String("post-aff-Lambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PYTHON_3_9(),
		Handler: jsii.String("index.post_aff_handler"),
		Code: awslambda.Code_FromAsset(jsii.String("lambda"), &awss3assets.AssetOptions{}),
		FunctionName: jsii.String("addAffirmation-Function"),
		Role: lambdaRole,
	})

	postlambdaintegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("Lambda"),
		postAffFunctionIntegration,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)

	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: getlambdaintegration,
		Path:        jsii.String("/affirm"),
		Methods:    &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_GET,
		},
	})

	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: postlambdaintegration,
		Path:        jsii.String("/affirm"),
		Methods:    &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_POST,
		},
	})

	awsdynamodb.NewTable(stack, jsii.String("aff-Table"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_NUMBER,
		},
		TableName: jsii.String("affirmations"),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewAppStack(app, "AppStack", &AppStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

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
