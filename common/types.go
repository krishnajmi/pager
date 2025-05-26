package common

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Context is used to store personalized data for the audience, such as their name,
// which can be used to personalize the notification template.
// Note: The "Email" field is used to store the audience's email address because
// we are planning to attach an email service for testing in phase 2.
type AudienceType struct {
	Email   string            `json:"email"`
	Context map[string]string `json:"context"`
}
