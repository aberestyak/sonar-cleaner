package sonarclient

import (
	"io/ioutil"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func DoHttpRequest(request *http.Request) (*[]byte, error) {
	client := &http.Client{}
	log.Tracef("Request URL: %s", request.URL.String())
	resp, err := client.Do(request)

	if err != nil {
		return nil, errors.Wrap(err, "Couldn't make the request")
	}
	defer resp.Body.Close()
	if 200 > resp.StatusCode || resp.StatusCode > 300 {
		return nil, errors.Wrapf(err, "Not successful request, response code %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.Wrap(err, "Couldn't get response body")
	}
	log.Tracef("Response body: %s", string(body))
	return &body, nil
}
