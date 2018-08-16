package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/buger/jsonparser"
)

// gitlab-pipeline-trigger https://gitlab.com XXXXX projectid pipelineid jobname

func main() {
	uri := os.Args[1]
	token := os.Args[2]
	projectID := os.Args[3]
	pipelineID := os.Args[4]
	jobName := os.Args[5]

	resultData := getPipelineJobs(uri, token, projectID, pipelineID)
	jobID := findJobIDByName(resultData, jobName)
	triggerJob(uri, token, projectID, jobID)
}

func getPipelineJobs(uri string, token string, projectID string, pipelineID string) []byte {
	url := uri + "/api/v4/projects/" + projectID + "/pipelines/" + pipelineID + "/jobs"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("PRIVATE-TOKEN", token)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	resultData, _ := ioutil.ReadAll(res.Body)
	return resultData
}

func findJobIDByName(resultData []byte, jobName string) int {
	count := 0
	jobID := 0
	jsonparser.ArrayEach(
		resultData,
		func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			name, _ := jsonparser.GetString(value, "name")
			if jobName == name {
				id, _ := jsonparser.GetInt(resultData, "["+strconv.Itoa(count)+"]", "id")
				jobID = int(id)
			}
			count++
		},
	)

	return jobID
}

func triggerJob(uri string, token string, projectID string, jobID int) {

	url := uri + "/api/v4/projects/" + projectID + "/jobs/" + strconv.Itoa(jobID) + "/play"
	req, _ := http.NewRequest("POST", url, nil)

	req.Header.Add("PRIVATE-TOKEN", token)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}
