package aws_support

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/support"
)

const checkID = "Pfx0RwqBli" // S3 Bucket Permissions check

// RefreshtS3BucketPermissionsCheck is used to refresh the S3 Bucket Permissions check
func (s *AwsSupport) RefreshtS3BucketPermissionsCheck() (*support.TrustedAdvisorCheckRefreshStatus, error) {
	input := &support.RefreshTrustedAdvisorCheckInput{
		CheckId: aws.String(checkID),
	}

	output, err := s.Client.RefreshTrustedAdvisorCheck(input)
	if err != nil {
		return nil, err
	}

	return output.Status, nil
}

// DescribeS3BucketPermissionsCheck is used to retrieve results from  S3 Bucket Permissions check
func (s *AwsSupport) DescribeS3BucketPermissionsCheck() (*support.TrustedAdvisorCheckResult, error) {
	input := &support.DescribeTrustedAdvisorCheckResultInput{
		CheckId: aws.String(checkID),
	}

	output, err := s.Client.DescribeTrustedAdvisorCheckResult(input)
	if err != nil {
		return nil, err
	}

	return output.Result, nil
}
