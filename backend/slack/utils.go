package slack

import (
	"bytes"
	"fmt"
	"github.com/slack-go/slack"
	"io/ioutil"
	"net/http"
)

func verifySigningSecret(r *http.Request, signingSecret string) error {
	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	// Need to use r.Body again when unmarshalling SlashCommand and InteractionCallback
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	verifier.Write(body)
	if err = verifier.Ensure(); err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
