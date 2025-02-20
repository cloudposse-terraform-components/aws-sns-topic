package test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	helper "github.com/cloudposse/test-helpers/pkg/atmos/component-helper"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/assert"
)

type SnsTopic struct {
	ApplicationFailureFeedbackRoleArn    string                 `json:"application_failure_feedback_role_arn"`
	ApplicationSuccessFeedbackRoleArn    string                 `json:"application_success_feedback_role_arn"`
	ApplicationSuccessFeedbackSampleRate int                    `json:"application_success_feedback_sample_rate"`
	ArchivePolicy                        string                 `json:"archive_policy"`
	Arn                                  string                 `json:"arn"`
	BeginningArchiveTime                 string                 `json:"beginning_archive_time"`
	ContentBasedDeduplication            bool                   `json:"content_based_deduplication"`
	DeliveryPolicy                       string                 `json:"delivery_policy"`
	DisplayName                          string                 `json:"display_name"`
	FifoTopic                            bool                   `json:"fifo_topic"`
	FirehoseFailureFeedbackRoleArn       string                 `json:"firehose_failure_feedback_role_arn"`
	FirehoseSuccessFeedbackRoleArn       string                 `json:"firehose_success_feedback_role_arn"`
	FirehoseSuccessFeedbackSampleRate    int                    `json:"firehose_success_feedback_sample_rate"`
	HttpFailureFeedbackRoleArn           string                 `json:"http_failure_feedback_role_arn"`
	HttpSuccessFeedbackRoleArn           string                 `json:"http_success_feedback_role_arn"`
	HttpSuccessFeedbackSampleRate        int                    `json:"http_success_feedback_sample_rate"`
	Id                                   string                 `json:"id"`
	KmsMasterKeyId                       string                 `json:"kms_master_key_id"`
	LambdaFailureFeedbackRoleArn         string                 `json:"lambda_failure_feedback_role_arn"`
	LambdaSuccessFeedbackRoleArn         string                 `json:"lambda_success_feedback_role_arn"`
	LambdaSuccessFeedbackSampleRate      int                    `json:"lambda_success_feedback_sample_rate"`
	Name                                 string                 `json:"name"`
	NamePrefix                           string                 `json:"name_prefix"`
	Owner                                string                 `json:"owner"`
	Policy                               string                 `json:"policy"`
	SignatureVersion                     int                    `json:"signature_version"`
	SqsFailureFeedbackRoleArn            string                 `json:"sqs_failure_feedback_role_arn"`
	SqsSuccessFeedbackRoleArn            string                 `json:"sqs_success_feedback_role_arn"`
	SqsSuccessFeedbackSampleRate         int                    `json:"sqs_success_feedback_sample_rate"`
	Tags                                 map[string]interface{} `json:"tags"`
	TagsAll                              map[string]interface{} `json:"tags_all"`
	TracingConfig                        string                 `json:"tracing_config"`
}

type ComponentSuite struct {
	helper.TestSuite
}

func (s *ComponentSuite) TestBasic() {
	const component = "sns-topic/basic"
	const stack = "default-test"
	const awsRegion = "us-east-2"

	defer s.DestroyAtmosComponent(s.T(), component, stack, nil)
	options, _ := s.DeployAtmosComponent(s.T(), component, stack, nil)
	assert.NotNil(s.T(), options)

	var snsTopic SnsTopic
	atmos.OutputStruct(s.T(), options, "sns_topic_name", &snsTopic)

	snsTopicId := atmos.Output(s.T(), options, "sns_topic_id")
	assert.Equal(s.T(), snsTopic.Name, snsTopicId)

	snsTopicOwner := atmos.Output(s.T(), options, "sns_topic_owner")
	assert.Equal(s.T(), snsTopic.Owner, snsTopicOwner)

	snsTopicArn := atmos.Output(s.T(), options, "sns_topic_arn")
	assert.Equal(s.T(), fmt.Sprintf("arn:aws:sns:%s:%s:%s", awsRegion, snsTopicOwner, snsTopicId), snsTopicArn)
	assert.Equal(s.T(), snsTopic.Arn, snsTopicArn)

	snsTopicSubscriptions := atmos.OutputMapOfObjects(s.T(), options, "sns_topic_subscriptions")
	assert.NotNil(s.T(), snsTopicSubscriptions)

	client := aws.NewSnsClient(s.T(), awsRegion)

	topicAttributes, err := client.GetTopicAttributes(context.Background(), &sns.GetTopicAttributesInput{
		TopicArn: &snsTopicArn,
	})
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), topicAttributes)

	assert.Equal(s.T(), snsTopic.Name, topicAttributes.Attributes["DisplayName"])
	assert.Equal(s.T(), snsTopic.Owner, topicAttributes.Attributes["Owner"])
	assert.Equal(s.T(), snsTopic.Arn, topicAttributes.Attributes["TopicArn"])
	assert.Equal(s.T(), snsTopic.DisplayName, topicAttributes.Attributes["DisplayName"])

	applicationSuccessFeedbackSampleRate, err := strconv.Atoi(topicAttributes.Attributes["ApplicationSuccessFeedbackSampleRate"])
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), snsTopic.ApplicationSuccessFeedbackSampleRate, applicationSuccessFeedbackSampleRate)

	firehoseSuccessFeedbackSampleRate, err := strconv.Atoi(topicAttributes.Attributes["FirehoseSuccessFeedbackSampleRate"])
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), snsTopic.FirehoseSuccessFeedbackSampleRate, firehoseSuccessFeedbackSampleRate)

	lambdaSuccessFeedbackSampleRate, err := strconv.Atoi(topicAttributes.Attributes["LambdaSuccessFeedbackSampleRate"])
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), snsTopic.LambdaSuccessFeedbackSampleRate, lambdaSuccessFeedbackSampleRate)

	s.DriftTest(component, stack, nil)
}

func (s *ComponentSuite) TestEnabledFlag() {
	const component = "sns-topic/disabled"
	const stack = "default-test"
	const awsRegion = "us-east-2"

	s.VerifyEnabledFlag(component, stack, nil)
}

func TestRunSuite(t *testing.T) {
	suite := new(ComponentSuite)
	helper.Run(t, suite)
}
