package janus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (j *Janus) addAPI(body []byte) error {
	url := fmt.Sprintf("%s/apis", j.URL)
	_, err := do("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "failed to add api")
	}
	return nil
}

func (j *Janus) modifyAPI(name string, body []byte) error {
	url := fmt.Sprintf("%s/apis/%s", j.URL, name)
	_, err := do("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "failed to modify api")
	}
	return nil
}

func (j *Janus) getAPI(name string) (*Service, error) {
	url := fmt.Sprintf("%s/apis/%s", j.URL, name)
	resp, err := do("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get api")
	}

	if resp == nil {
		return nil, nil
	}

	var janus Service
	err = json.Unmarshal(resp, &janus)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal to JanusService struct")
	}
	return &janus, nil
}

func (j *Janus) deleteAPI(name string) error {
	url := fmt.Sprintf("%s/apis/%s", j.URL, name)
	_, err := do("DELETE", url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to delete api")
	}
	return nil
}

func do(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}
	defer response.Body.Close()
	bodyBB, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated:
		return bodyBB, nil
	case http.StatusNotFound, http.StatusNoContent:
		return nil, nil
	case http.StatusBadRequest:
		return nil, E(
			"request was invalid in some way",
			Response(bodyBB),
			Status(response.StatusCode),
		)
	case http.StatusUnauthorized:
		return nil, E(
			"not allowed",
			Response(bodyBB),
			Status(response.StatusCode),
		)
	default:
		return nil, E(
			"something bad happened",
			Response(bodyBB),
			Status(response.StatusCode),
		)
	}
}
