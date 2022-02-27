# Remote Lambda Invoke using Go

`aws-remote-lambda-invoke-go` is a simple example of calling a Lambda function remotely without using an AWS API Gateway or the full SDK. The code creates an AWS SigV4 header, and then invokes the Lambda through the HTTP Endpoint. It relies on AWS credentials already being configured on the local machine

## Getting started

### Create IAM User

Following the security principle of least privilege, I recommend creating an API user in IAM that only has access to invoke the single required Lambda function. An example policy would look like this:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": "lambda:InvokeFunction",
            "Resource": "arn:aws:lambda:eu-west-2:123123123123:function:my-function-name"
        }
    ]
}
```

### Configure AWS Credentials

Using the API Key and Secret from the user that you created in IAM

```sh
aws configure
```

### Update code

There are a number of constants you'll want to update for your function

```go
const REGION string = "eu-west-2" // CHANGE THIS VALUE
const FUNCTION_NAME string = "## INSERT LAMBDA FUNCTION NAME ##" // CHANGE THIS VALUE
const PAYLOAD string = "{\"example-key\":\"example-value\"}" // CHANGE THIS VALUE
```

### Run

```sh
go run .
```

```sh
13:53:28.269864 Initialising
13:53:28.270507 Creating Request
13:53:28.270507 Generating Hash
13:53:28.270507 Hash: 2a92f9bb852a92f9bb852a92f9bb852a92f9bb852a92f9bb852a92f9bb852a92
13:53:28.270507 Signature: AWS4-HMAC-SHA256 Credential=ABCDEFABCDEFABCDEFABC/20220227/eu-west-2/lambda/aws4_request, SignedHeaders=content-length;host;x-amz-date, Signatur=2a92f9bb852a92f9bb852a92f9bb852a92f9bb852a92f9bb8a92f9bb852a2323
13:53:28.271015 Making web request
13:53:29.440913 Return Status: 200 OK
13:53:29.440913 Body: {"response-json-message":"example-value"}
```

## Useful Links

- <https://docs.aws.amazon.com/general/latest/gr/lambda-service.html>
- <https://docs.aws.amazon.com/lambda/latest/dg/API_Invoke.html>
- <https://github.com/aws/aws-sdk-go-v2>
- <https://github.com/aws/aws-sdk-go-v2/blob/main/aws/signer/v4/v4.go>
