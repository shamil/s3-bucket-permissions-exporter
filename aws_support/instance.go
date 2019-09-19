package aws_support

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/support"
	"github.com/prometheus/common/log"
)

// AwsSupport is an instance of AWS support client
type AwsSupport struct {
	Client *support.Support
}

// New creates new Trusted Advisor API client
func New() *AwsSupport {
	config := aws.NewConfig().WithCredentialsChainVerboseErrors(true).WithRegion("us-east-1") // Trusted Advisor API needs 'us-east-1' region
	sess := session.Must(session.NewSession(config))

	if _, err := (sess.Config.Credentials.Get()); err != nil {
		log.Fatal(err)
	}

	return &AwsSupport{
		Client: support.New(sess),
	}
}
