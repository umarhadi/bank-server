package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/umarhadi/bank-server/util"
)

func TestSendEmailWithGmail(t *testing.T) {

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="https://bank.api.umarhadi.dev">very bank service</a></p>
	`
	to := []string{"hi@umarhadi.dev"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}

func TestSendEmailWithInvalidAttachment(t *testing.T) {
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := "This is a test message"
	to := []string{"hi@umarhadi.dev"}
	attachFiles := []string{"non_existent_file.txt"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to attach file")
}
