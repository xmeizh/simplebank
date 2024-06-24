package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmeizh/simplebank/util"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Simple Bank Test Email"
	content := `
	<h1>Welcome to Simple Bank Service</h1>
	<p>This is a test message from Simple Bank</p>
	`

	to := []string{"dummy@gmail.com"}
	attachments := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachments)
	require.NoError(t, err)
}
