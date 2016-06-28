package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	//"strings"

	"crypto/tls"
	"strings"
	"log"
	"os"
)

var Logger log.Logger
var Conf *Config

func init() {
	Logger = *log.New(os.Stderr, "itba: ", log.LstdFlags)
}

type Jenkins struct {
	Url      string
	Username string
	Password string
	client   *http.Client
	cache    map[string]Job
}

func NewJenkins(Url, Username, Password string) *Jenkins {
	jk := &Jenkins{
		Url:      Url,
		Username: Username,
		Password: Password,
		client:   &http.Client{},
		cache:    make(map[string]Job),
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// TODO: should find a better way to handle "x509: certificate signed by unknown authority panic:"
	jk.client = &http.Client{Transport: tr}

	return jk
}

// JobDetails executes the API call to collect the detailed information about the jenkins job. The results are cached.
func (jk *Jenkins) JobDetails(name string) Job {

	Logger.Printf("Getting job details: %s\n", name)
	if item, ok := jk.cache[name]; ok {
		Logger.Printf("Details found in cache [%d] %s\n", len(jk.cache), name)
		return item
	}
	req, _ := http.NewRequest("GET", jk.Url+"/job/"+name+"/api/json", nil) // 某个Job的信息
	//req.SetBasicAuth(jk.Username, jk.Password)
	req.Header.Set("Accept", "application/json")

	res, _ := jk.client.Do(req)
	//if (err != nil) {
	//	return
	//}
	b1, _ := ioutil.ReadAll(res.Body)
	//if (err != nil) {
	//	return
	//}

	defer res.Body.Close()
	aj := &Job{}
	json.Unmarshal(b1, aj)
	for i, ab := range aj.Builds {
		url := fmt.Sprintf(jk.Url+"/job/%s/%d/api/json", name, ab.Number) // 某次Build的输出信息

		req2, _ := http.NewRequest("GET", url, nil)
		//req2.SetBasicAuth(jk.Username, jk.Password)
		req2.Header.Set("Accept", "application/json")
		res2, _ := jk.client.Do(req2)
		b2, _ := ioutil.ReadAll(res2.Body)
		defer res2.Body.Close()
		bld := &Build{}
		err := json.Unmarshal(b2, bld)
		aj.Builds[i] = *bld
		if err != nil {
			fmt.Printf("Error parsing build info: %s\n", err)
		}

	}
	jk.cache[name] = *aj
	Logger.Printf("Added to cache [%d] %s\n", len(jk.cache), name)
	return *aj
}

type Jobs struct {
	Desc string `json:"nodeDescription"`
	Jobs []Job  `json:"jobs"`
}

type Job struct {
	Name   string  `json:"name"`
	Url    string  `json:"url"`
	Builds []Build `json:"builds"`
}

type Build struct {
	Number    int      `json:"number"`
	Duration  int      `json:"duration,omitempty"`
	Result    string   `json:"result,omitempty"`
	Timestamp int64    `json:"timestamp,omitempty"`
	Actions   []Action `json:"actions,omitempty"`
}

// Stage returns the stage of a given job by looking up on the build chain to find the first job matching the suffix search
func (bld Build) Stage() string {
	stg := "UNK"
	Logger.Printf("stage processing build %d - %s\n", bld.Number, stg)
	for _, a := range bld.Actions {
		//for i, c := range a.Causes {
		//	for _, js := range conf.JobStage {
		//		if strings.Contains(c.Up, js.Suffix) {
		//			stg = js.Stage
		//			break
		//		}
		//	}
		//	Logger.Printf("stage causes [%d] %s - stage: %s\n", i, c.Up, stg)
		//	if stg != "UNK" {
		//		break
		//	}
		//}
		for _, p := range a.Parameters {
			if (p.Name == "ENVIRONMENT") {
				stg = strings.ToUpper(p.Value)
				break
			}
		}

	}
	Logger.Printf("stage return: %s\n", stg)
	return stg
}

type Action struct {
	Causes     []Upstream  `json:"causes,omitempty`
	Parameters []Parameter `json:"parameters,omitempty`
}
type Upstream struct {
	Desc    string `json:"shortDescription"`
	Up      string `json:"upstreamProject,omitempty"`
	UpBuild int    `json:"upstreamBuild,omitempty"`
}
type Parameter struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// GetUpstreamJob recursively walks the job call tree and loops through the jobs Actions and Causes to identify the pipeline the job belongs to.
// The pipeline is identified by applying the job-pipeline mapping defined in the configuration file
func GetUpstreamJob(jk *Jenkins, name string, build int) (int, string) {
	Logger.Printf("Getting upstream job: %s-%d\n", name, build)
	jd := jk.JobDetails(name)
	for _, b := range jd.Builds {
		if b.Number == build {
			Logger.Printf("Found build : %d-%d\n", b.Number, build)
			for i, a := range b.Actions {
				Logger.Printf("Processing action (%d/%d)%d\n", i, len(b.Actions), len(a.Causes))
				for j, c := range a.Causes {
					Logger.Printf("Processing causes (%d/%d) job: %s-%d\n", j, len(a.Causes), c.Up, c.UpBuild)
					if _, ok := Conf.JobPipelines[c.Up]; ok {
						Logger.Printf("Found pipeline %s\n", Conf.JobPipelines[c.Up])
						return c.UpBuild, Conf.JobPipelines[c.Up]
					} else {
						if c.Up != "" {
							i, j := GetUpstreamJob(jk, c.Up, c.UpBuild)
							return i, j
						}
					}
				}
			}
		}
	}
	Logger.Printf("Exit criteria not found returning 0,NotFound\n")
	return 0, "NotFound"
}
