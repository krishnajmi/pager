package batchprocessor

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/kp/pager/common"
	"github.com/kp/pager/communicator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) Publish(ctx context.Context, topic string, message []byte) error {
	args := m.Called(ctx, topic, message)
	return args.Error(0)
}

type MockNotificationType struct {
	communicator.NotificationType
}

func (m MockNotificationType) Save(ctx context.Context) error {
	return nil
}

func (m MockNotificationType) Validate(ctx context.Context) error {
	return nil
}

func (m MockNotificationType) Prepare(ctx context.Context) (interface{}, error) {
	return nil, nil
}

func (m MockNotificationType) Send(ctx context.Context, payload interface{}) error {
	return nil
}

func TestNewBatchProcessor(t *testing.T) {
	ctx := context.Background()
	audiences := []common.AudienceType{
		{Email: "test1@example.com", Context: map[string]string{}},
		{Email: "test2@example.com", Context: map[string]string{}},
	}
	model := communicator.NotificationType{
		To:         "test@example.com",
		TemplateID: 123,
		RequestId:  "req123",
	}
	topic := "test-topic"
	producer := &MockKafkaProducer{}

	processor := NewBatchProcessor(ctx, audiences, model, topic, producer)

	assert.NotNil(t, processor)
	assert.IsType(t, &BatchChannelBased{}, processor)
}

func TestProcess_Success(t *testing.T) {
	tests := []struct {
		name          string
		batchSize     int
		expectedCalls int
	}{
		{"Small batch", 4, 1},
		{"Exact batch multiple", 10, 2},
		{"Larger than batch", 13, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			audiences := make([]common.AudienceType, tt.batchSize)
			for i := 0; i < tt.batchSize; i++ {
				audiences[i] = common.AudienceType{
					Email:   string(rune(i+1)) + "@example.com",
					Context: map[string]string{},
				}
			}

			model := communicator.NotificationType{
				To:         "test@example.com",
				TemplateID: 123,
				RequestId:  "req123",
			}
			topic := "test-topic"
			producer := &MockKafkaProducer{}
			producer.On("Publish", mock.Anything, topic, mock.Anything).Return(nil)

			processor := NewBatchProcessor(ctx, audiences, model, topic, producer)
			err := processor.Process(ctx)

			assert.NoError(t, err)
			producer.AssertNumberOfCalls(t, "Publish", tt.expectedCalls)
		})
	}
}

func TestProcess_Concurrency(t *testing.T) {
	ctx := context.Background()
	audiences := make([]common.AudienceType, 100) // Large enough to test concurrency
	for i := 0; i < 100; i++ {
		audiences[i] = common.AudienceType{
			Email:   string(rune(i+1)) + "@example.com",
			Context: map[string]string{},
		}
	}

	model := communicator.NotificationType{
		To:         "test@example.com",
		TemplateID: 123,
		RequestId:  "req123",
	}
	topic := "test-topic"
	producer := &MockKafkaProducer{}
	producer.On("Publish", mock.Anything, topic, mock.Anything).Return(nil)

	processor := NewBatchProcessor(ctx, audiences, model, topic, producer)
	err := processor.Process(ctx)

	assert.NoError(t, err)
	minExpectedCalls := 100 / 5 // At least 20 batches (100 items / 5 per batch)
	assert.GreaterOrEqual(t, len(producer.Calls), minExpectedCalls)
}

func TestProcess_Error(t *testing.T) {
	ctx := context.Background()
	audiences := make([]common.AudienceType, 10)
	for i := 0; i < 10; i++ {
		audiences[i] = common.AudienceType{
			Email:   string(rune(i+1)) + "@example.com",
			Context: map[string]string{},
		}
	}

	model := communicator.NotificationType{
		To:         "test@example.com",
		TemplateID: 123,
		RequestId:  "req123",
	}
	topic := "test-topic"
	producer := &MockKafkaProducer{}
	producer.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("kafka error"))

	processor := NewBatchProcessor(ctx, audiences, model, topic, producer)
	err := processor.Process(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "kafka error")
}

func TestSendBatchToQueue_Success(t *testing.T) {
	ctx := context.Background()
	audiences := []common.AudienceType{
		{Email: "test1@example.com", Context: map[string]string{}},
		{Email: "test2@example.com", Context: map[string]string{}},
	}
	model := communicator.NotificationType{
		To:         "test@example.com",
		TemplateID: 123,
		RequestId:  "req123",
	}
	topic := "test-topic"
	producer := &MockKafkaProducer{}
	producer.On("Publish", ctx, topic, mock.Anything).Return(nil)

	processor := &BatchChannelBased{
		Model:         model,
		TopicName:     topic,
		Audiences:     nil, // Not used in sendBatchToQueue
		KafkaProducer: producer,
	}

	err := processor.sendBatchToQueue(ctx, audiences)
	assert.NoError(t, err)

	// Verify the published message contains expected data
	producer.AssertCalled(t, "Publish", ctx, topic, mock.Anything)
	capturedMsg := producer.Calls[0].Arguments[2].([]byte)

	var msg communicator.QMessage
	err = json.Unmarshal(capturedMsg, &msg)
	assert.NoError(t, err)
	assert.Equal(t, model, msg.GenericModel)
	assert.Equal(t, audiences, msg.Audiences)
	assert.NotEmpty(t, msg.BatchID)
}

func TestSendBatchToQueue_MarshalError(t *testing.T) {
	ctx := context.Background()
	audiences := []common.AudienceType{
		{Email: "test1@example.com", Context: map[string]string{}},
		{Email: "test2@example.com", Context: map[string]string{}},
	}
	topic := "test-topic"
	producer := &MockKafkaProducer{}

	// Create a test processor with a mock marshal function
	processor := &BatchChannelBased{
		Model: communicator.NotificationType{
			To:         "test@example.com",
			TemplateID: 123,
			RequestId:  "req123",
		},
		TopicName:     topic,
		Audiences:     audiences,
		KafkaProducer: producer,
	}

	// Save original json.Marshal
	originalMarshal := jsonMarshal
	defer func() { jsonMarshal = originalMarshal }()

	// Mock marshal to return error
	jsonMarshal = func(v interface{}) ([]byte, error) {
		return nil, errors.New("forced marshal error")
	}

	err := processor.sendBatchToQueue(ctx, audiences)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal message")
	producer.AssertNotCalled(t, "Publish")
}

// jsonMarshal is a variable that can be replaced in tests
var jsonMarshal = json.Marshal
