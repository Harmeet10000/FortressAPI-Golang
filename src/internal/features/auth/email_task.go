package auth

import (
	"context"
	"encoding/json"
	"fmt"

	// "github.com/Harmeet10000/Fortress_API/src/internal/helpers/email"
	"github.com/hibiken/asynq"
)

const TypeWelcomeEmail = "email:welcome"

type WelcomeEmailPayload struct {
	UserID int
	Email  string
}

// Task Producer: Use this in your API handlers
func NewWelcomeEmailTask(userID int, email string) (*asynq.Task, error) {
	payload, err := json.Marshal(WelcomeEmailPayload{UserID: userID, Email: email})
	if err != nil {
		return nil, err
	}
	// We set the queue to 'critical' as per your config requirement
	return asynq.NewTask(TypeWelcomeEmail, payload, asynq.Queue("critical")), nil
}

// Task Handler: The worker will execute this
func HandleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var p WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	fmt.Printf("Sending welcome email to User %d at %s\n", p.UserID, p.Email)
	// err := email.SendWelcomeEmail(p.Email, p.UserID)
	// if err != nil {
	// 	return fmt.Errorf("send welcome email failed: %w", err)
	// }
	return nil
}
