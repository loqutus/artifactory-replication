package artifactory

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetArtifactoryFileSHA256(host string, fileName string, user string, pass string) (string, error) {
	url := "https://" + host + "/artifactory/api/storage/" + fileName
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(user, pass)
	var resp *http.Response
	var failed bool
	backOffTime := backOffStart
	for i := 1; i <= backOffSteps; i++ {
		resp, err = client.Do(req)
		if err != nil {
			failed = true
			log.Print("error HTTP GET", url, "retry", string(i))
			if i != backOffSteps {
				time.Sleep(time.Duration(backOffTime) * time.Millisecond)
			}
			backOffTime *= i
		} else {
			defer resp.Body.Close()
			failed = false
			break
		}
	}
	if failed == true {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	type storageInfo struct {
		Checksums map[string]string `json:"checksums"`
	}
	var result storageInfo
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	return result.Checksums["sha256"], nil
}
