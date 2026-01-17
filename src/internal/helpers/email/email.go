package email

import "github.com/Harmeet10000/Fortress_API/src/internal/helpers/email/templates"

func (c *Client) SendWelcomeEmail(to, firstName string) error {
	data := map[string]string{
		"UserFirstName": firstName,
	}

	return c.SendEmail(
		to,
		"Welcome to Boilerplate!",
		templates.TemplateWelcome,
		data,
	)
}
