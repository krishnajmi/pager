package communicator

import (
	"context"
	"errors"
	"testing"

	"github.com/kp/pager/communicator/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testNotification struct {
	NotificationType
	logCreator interface {
		CreateLog(ctx context.Context, db interface{}, to string, templateID int64, requestID string) (*models.CommunicationLogs, error)
	}
}

func (n *testNotification) Save(ctx context.Context) error {
	entry, err := n.logCreator.CreateLog(ctx, nil, n.To, n.TemplateID, n.RequestId)
	if err != nil {
		return err
	}
	n.LogID = entry.ID
	return nil
}

type mockLogCreator struct {
	mock.Mock
}

func (m *mockLogCreator) CreateLog(ctx context.Context, db interface{}, to string, templateID int64, requestID string) (*models.CommunicationLogs, error) {
	args := m.Called(ctx, db, to, templateID, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CommunicationLogs), args.Error(1)
}

func TestNotificationType_Save(t *testing.T) {
	tests := []struct {
		name         string
		notification *testNotification
		setupMock    func(*mockLogCreator)
		wantErr      bool
		errMsg       string
		wantLogID    int64
	}{
		{
			name: "successful save",
			notification: &testNotification{
				NotificationType: NotificationType{
					To:         "test@example.com",
					TemplateID: 123,
					RequestId:  "req123",
				},
			},
			setupMock: func(m *mockLogCreator) {
				m.On("CreateLog", mock.Anything, nil, "test@example.com", int64(123), "req123").
					Return(&models.CommunicationLogs{ID: 123}, nil)
			},
			wantLogID: 123,
		},
		{
			name: "database error",
			notification: &testNotification{
				NotificationType: NotificationType{
					To:         "test@example.com",
					TemplateID: 123,
					RequestId:  "req123",
				},
			},
			setupMock: func(m *mockLogCreator) {
				m.On("CreateLog", mock.Anything, nil, "test@example.com", int64(123), "req123").
					Return((*models.CommunicationLogs)(nil), errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "failed to save communication log",
		},
		{
			name: "empty request ID",
			notification: &testNotification{
				NotificationType: NotificationType{
					To:         "test@example.com",
					TemplateID: 123,
					RequestId:  "",
				},
			},
			setupMock: func(m *mockLogCreator) {
				m.On("CreateLog", mock.Anything, nil, "test@example.com", int64(123), "").
					Return(&models.CommunicationLogs{ID: 456}, nil)
			},
			wantLogID: 456,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockCreator := &mockLogCreator{}
			tt.setupMock(mockCreator)
			tt.notification.logCreator = mockCreator

			err := tt.notification.Save(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLogID, tt.notification.LogID)
			}
			mockCreator.AssertExpectations(t)
		})
	}
}

func TestNotificationType_Validate(t *testing.T) {
	tests := []struct {
		name         string
		notification *NotificationType
		wantErr      bool
		errMsg       string
	}{
		{
			name: "valid notification",
			notification: &NotificationType{
				To:         "test@example.com",
				TemplateID: 123,
				RequestId:  "req123",
			},
			wantErr: false,
		},
		{
			name: "empty To field",
			notification: &NotificationType{
				To:         "",
				TemplateID: 123,
				RequestId:  "req123",
			},
			wantErr: true,
			errMsg:  "To field is required",
		},
		{
			name: "invalid email",
			notification: &NotificationType{
				To:         "invalid-email",
				TemplateID: 123,
				RequestId:  "req123",
			},
			wantErr: true,
			errMsg:  "invalid email address",
		},
		{
			name: "zero TemplateID",
			notification: &NotificationType{
				To:         "test@example.com",
				TemplateID: 0,
				RequestId:  "req123",
			},
			wantErr: true,
			errMsg:  "TemplateID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := tt.notification.Validate(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNotificationType_Prepare(t *testing.T) {
	tests := []struct {
		name         string
		notification *NotificationType
		wantPayload  interface{}
		wantErr      bool
	}{
		{
			name: "with template data",
			notification: &NotificationType{
				To:         "test@example.com",
				TemplateID: 123,
				RequestId:  "req123",
				Context: map[string]string{
					"name": "John",
					"age":  "30",
				},
			},
			wantPayload: map[string]string{
				"name": "John",
				"age":  "30",
			},
		},
		{
			name: "without template data",
			notification: &NotificationType{
				To:         "test@example.com",
				TemplateID: 123,
				RequestId:  "req123",
			},
			wantPayload: nil,
		},
		{
			name: "empty template data",
			notification: &NotificationType{
				To:         "test@example.com",
				TemplateID: 123,
				RequestId:  "req123",
				Context:    map[string]string{},
			},
			wantPayload: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			payload, err := tt.notification.Prepare(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPayload, payload)
			}
		})
	}
}

func TestNotificationType_Send(t *testing.T) {
	tests := []struct {
		name         string
		notification *NotificationType
		sender       interface {
			Send(ctx context.Context, payload interface{}) error
		}
		setupMock func(*mockSender)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful send",
			notification: &NotificationType{
				To:         "test@example.com",
				TemplateID: 123,
				RequestId:  "req123",
			},
			sender: &mockSender{},
			setupMock: func(m *mockSender) {
				m.On("Send", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name: "send error",
			notification: &NotificationType{
				To:         "test@example.com",
				TemplateID: 123,
				RequestId:  "req123",
			},
			sender: &mockSender{},
			setupMock: func(m *mockSender) {
				m.On("Send", mock.Anything, mock.Anything).Return(errors.New("send failed"))
			},
			wantErr: true,
			errMsg:  "send failed",
		},
		{
			name: "nil sender",
			notification: &NotificationType{
				To:         "test@example.com",
				TemplateID: 123,
				RequestId:  "req123",
			},
			sender:  nil,
			wantErr: true,
			errMsg:  "sender is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if mockSender, ok := tt.sender.(*mockSender); ok {
				tt.setupMock(mockSender)
			}

			err := tt.notification.Send(ctx, tt.sender)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}

			if mockSender, ok := tt.sender.(*mockSender); ok {
				mockSender.AssertExpectations(t)
			}
		})
	}
}

type mockSender struct {
	mock.Mock
}

func (m *mockSender) Send(ctx context.Context, payload interface{}) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}
