package test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	helper "github.com/cloudposse/test-helpers/pkg/atmos/aws-component-helper"
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

func TestComponent(t *testing.T) {
	awsRegion := "us-east-2"

	fixture := helper.NewFixture(t, "../", awsRegion, "test/fixtures")

	defer fixture.TearDown()
	fixture.SetUp(&atmos.Options{})

	fixture.Suite("default", func(t *testing.T, suite *helper.Suite) {
		suite.Test(t, "basic", func(t *testing.T, atm *helper.Atmos) {
			inputs := map[string]interface{}{}

			defer atm.GetAndDestroy("sns-topic/basic", "default-test", inputs)
			component := atm.GetAndDeploy("sns-topic/basic", "default-test", inputs)
			assert.NotNil(t, component)

			var snsTopic SnsTopic

			atm.OutputStruct(component, "sns_topic_name", &snsTopic)

			snsTopicId := atm.Output(component, "sns_topic_id")
			assert.Equal(t, snsTopic.Name, snsTopicId)

			snsTopicOwner := atm.Output(component, "sns_topic_owner")
			assert.Equal(t, snsTopic.Owner, snsTopicOwner)

			snsTopicArn := atm.Output(component, "sns_topic_arn")
			assert.Equal(t, fmt.Sprintf("arn:aws:sns:%s:%s:%s", awsRegion, snsTopicOwner, snsTopicId), snsTopicArn)
			assert.Equal(t, snsTopic.Arn, snsTopicArn)

			snsTopicSubscriptions := atm.OutputMapOfObjects(component, "sns_topic_subscriptions")
			assert.NotNil(t, snsTopicSubscriptions)

			client := aws.NewSnsClient(t, awsRegion)

			topicAttributes, err := client.GetTopicAttributes(&sns.GetTopicAttributesInput{
				TopicArn: &snsTopicArn,
			})
			assert.NoError(t, err)
			assert.NotNil(t, topicAttributes)

			// You can add assertions for specific attributes if needed
			assert.Equal(t, snsTopic.Name, *topicAttributes.Attributes["DisplayName"])
			assert.Equal(t, snsTopic.Owner, *topicAttributes.Attributes["Owner"])
			assert.Equal(t, snsTopic.Arn, *topicAttributes.Attributes["TopicArn"])
			assert.Equal(t, snsTopic.DisplayName, *topicAttributes.Attributes["DisplayName"])

			applicationSuccessFeedbackSampleRate, err := strconv.Atoi(*topicAttributes.Attributes["ApplicationSuccessFeedbackSampleRate"])
			assert.NoError(t, err)
			assert.Equal(t, snsTopic.ApplicationSuccessFeedbackSampleRate, applicationSuccessFeedbackSampleRate)

			firehoseSuccessFeedbackSampleRate, err := strconv.Atoi(*topicAttributes.Attributes["FirehoseSuccessFeedbackSampleRate"])
			assert.NoError(t, err)
			assert.Equal(t, snsTopic.FirehoseSuccessFeedbackSampleRate, firehoseSuccessFeedbackSampleRate)

			lambdaSuccessFeedbackSampleRate, err := strconv.Atoi(*topicAttributes.Attributes["LambdaSuccessFeedbackSampleRate"])
			assert.NoError(t, err)
			assert.Equal(t, snsTopic.LambdaSuccessFeedbackSampleRate, lambdaSuccessFeedbackSampleRate)

			// deadLetterQueueUrl := atm.Output(component, "dead_letter_queue_url")
			// assert.Equal(t, deadLetterQueueUrl, "")

			// deadLetterQueueId := atm.Output(component, "dead_letter_queue_id")
			// assert.Equal(t, deadLetterQueueId, "")

			// deadLetterQueueName := atm.Output(component, "dead_letter_queue_name")
			// assert.Equal(t, deadLetterQueueName, "")

			// deadLetterQueueArn := atm.Output(component, "dead_letter_queue_arn")
			// assert.Equal(t, deadLetterQueueArn, "")
		})

	})
}
