package main

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

type Squeue struct {
	JobId     string `json:"job_id"`
	Partition string `json:"partition"`
	Name      string `json:"name"`
	User      string `json:"user"`
	Status    string `json:"status"`
	Time      string `json:"time"`
	Nodes     string `json:"nodes"`
	NodeList  string `json:"node_list"`
}

type Sbatch struct {
	JobId string `json:"job_id"`
}

type Scontrol struct {
	JobState string `json:"job_state"`
}

func runSqueue() ([]Squeue, error) {
	out, errMsg, err := runShell("squeue")
	if err != nil {
		return nil, errors.Wrap(err, string(errMsg))
	}
	lines := strings.Split(string(out), "\n")
	lines = lines[:len(lines)-1]

	if len(lines) == 1 {
		return nil, nil
	}
	for i := range lines {
		lines[i] = standardizeSpaces(lines[i])
	}
	res := make([]Squeue, len(lines)-1)
	for i := range res {
		x := strings.Split(lines[i+1], " ")
		res[i] = Squeue{x[0], x[1], x[2], x[3], x[4], x[5], x[6], x[7]}
	}

	return res, nil
}

func runSbatch(workDir string, arg string, file string) (string, error) {
	cmd := fmt.Sprintf("cd %s && sbatch %s %s", workDir, arg, file)
	out, errMsg, err := runShell(cmd)
	if err != nil {
		return "", errors.Wrap(err, string(errMsg))
	}

	rgx := regexp.MustCompile("([0-9]+)")
	sid := rgx.FindString(string(out))
	if sid == "" {
		return "", errors.Errorf("No job id was returned")
	}

	resJson, err := jsoniter.Marshal(Sbatch{sid})
	if err != nil {
		return "", errors.Wrap(err, "parse to json failed!")
	}

	return string(resJson), nil
}

func runScontrol(jobId string) (string, error) {
	cmd := fmt.Sprintf("scontrol show job %s", jobId)
	out, errMsg, err := runShell(cmd)
	if err != nil {
		return "", errors.Wrap(err, string(errMsg))
	}

	rgx := regexp.MustCompile("JobState=([A-Z]+)")
	s := rgx.FindSubmatch(out)
	if len(s) < 2 {
		return "", errors.Errorf("Failed to find job state")
	}
	resJson, err := jsoniter.Marshal(Scontrol{string(s[1])})
	if err != nil {
		return "", errors.Wrap(err, "parse to json failed!")
	}

	return string(resJson), nil
}

func runScancel(jobId string) error {
	cmd := fmt.Sprintf("scancel %s", jobId)
	_, errMsg, err := runShell(cmd)
	if err != nil {
		return errors.Wrap(err, string(errMsg))
	}

	return nil
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func runShell(cmd string) ([]byte, []byte, error) {
	command := exec.Command("sh", "-c", cmd)
	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	err = command.Start()
	if err != nil {
		return nil, nil, err
	}

	out, _ := ioutil.ReadAll(stdout)
	errMsg, _ := ioutil.ReadAll(stderr)

	err = command.Wait()

	return out, errMsg, err
}
