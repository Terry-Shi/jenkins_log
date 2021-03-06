// Itba command exports the build data from jenkins to csv format to be consumed by ETL. The exported data are the following:
//
// JOB_ID", "REL_ID", "JOB_TYPE", "PIPELINE", "STAGE", "JOB_NAME", "START_TIME", "JOB_RESULT", "DURATION"
//
// where
//
// 	JOB_ID: Jenkins job number
// 	REL_ID: id of the release job (start job of the pipeline)
// 	JOB_TYPE: the type of the job as defined in the mapping config.yaml)
// 	PIPELINE: pipeline the job belongs to as defined in the mapping between the start job and the pipeline tag (see config.yaml)
// 	STAGE: stage in the pipeline as defined in the mapping between the parent job and the suffix (see config.yaml)
// 	JOB_NAME: job name
// 	START_TIME: date time of the job start
// 	JOB_RESULT: SUCCESS,FAILURE,UNSTABLE jenkins job statuses
// 	DURATION: job duration in ms
//
// Not all the job data are exported, it is possible to configure the jobs that are exported. The export works based on the following assumptions:
//
// 1. The jobs pipeline is defined with the paramterized job plugin (not tested with different plugins)
//
// 2. The stage of the job is defined in the parent job name suffix
//
// 3. The start job to pipeline mapping can be defined (multiple pipelines with the same start job are not supported)
//
// The configuration file format is yaml and expects to have defined the list of jobs to be exported (jobFilter), the mapping between the pipeline and the job name (jobPipeline) and the mapping between the suffix of the job name and the corresponding stage in the pipeline (jobStage).
// Here follows an example of the config file:
//	jobFilter:
//	 - jobName: <jobname>
//	 - jobType: <jobtype>
//	jobPipeline:
//	<jobname>:<pipeline>
//	jobStage:
//	 - suffix: <suffix>
//	 - stage: <stage>
//
// Usage:
//	itba -URL=<url> -u=<username> -p=<password> -config=<absolute path of config file>
//
// The log information is returned on stderr while the output (csv) is returned on stdout
package main

import (
	"flag"
	"fmt"
	//"log"
	"os"
	"time"
	//"strings"
	."./util"
)



func main() {

	var URL string
	var username string
	var password string
	var cfile string
	// windowns command line
	//go run itba.go -URL=https://cas-cd.core.hpecorp.net:2543 -u=1 -p=1 -config=C:\WORK\IDE\IdeaProjects\jenkins_log\config.yaml 1> JENKINS_JOB_CSTM.csv 2> itba.log
	flag.StringVar(&URL, "URL", "https://cas-cd.core.hpecorp.net:2543", "Jenkins server URL")
	flag.StringVar(&username, "u", "1", "username to connect to jenkins")
	flag.StringVar(&password, "p", "1", "password to connect to jenkins")
	flag.StringVar(&cfile, "config", "C:\\WORK\\IDE\\IdeaProjects\\jenkins_log\\config.yaml", "abs path of config file")

	flag.Parse()

	if len(os.Args) == 1 {
		flag.PrintDefaults()
		return
	}
	Logger.Printf("==== itba start ====\n")

	Conf = &Config{}
	LoadConfig(Conf, cfile)
	Logger.Printf("Loaded [%d] job filters\n", len(Conf.Filters))
	Logger.Printf("Loaded [%d] job pipeline mappings \n", len(Conf.JobPipelines))
	Logger.Printf("Loaded [%d] job stage mappings \n", len(Conf.JobStage))

	jk := NewJenkins(URL, username, password)
	fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,%s\n", "JOB_ID", "REL_ID", "JOB_TYPE", "PIPELINE", "STAGE", "JOB_NAME", "START_TIME", "JOB_RESULT", "DURATION")
	for _, f := range Conf.Filters { // 遍历所有Job
		Logger.Printf("Getting details for job: %s\n", f.JobName)
		jd := jk.JobDetails(f.JobName) // 读取一个Job的信息
		for _, b := range jd.Builds {
			Logger.Printf("Processing build %s-%d\n", f.JobName, b.Number)
			ts := time.Unix(0, b.Timestamp*int64(time.Millisecond))
			if b.Result != "" {
				Logger.Printf("Getting parent upstream job %s-%d\n", f.JobName, b.Number)
				rid, pipeline := GetUpstreamJob(jk, f.JobName, b.Number) // TODO should get "pipeline" when rid==0
				if rid != 0 {
					stage := b.Stage()
					//stage := getStage(f.JobName)
					fmt.Printf("%d,%d,%s,%s,%s,%s,%s,%s,%d\n", b.Number, rid, f.JobType, pipeline, stage, f.JobName, ts.Format(time.RFC3339), b.Result, b.Duration)
					Logger.Printf("%d,%d,%s,%s,%s,%s,%s,%s,%d\n", b.Number, rid, f.JobType, pipeline, stage, f.JobName, ts.Format(time.RFC3339), b.Result, b.Duration)
				} else {
					stage := b.Stage()
					//stage := getStage(f.JobName)
					fmt.Printf("%d,%d,%s,%s,%s,%s,%s,%s,%d\n", b.Number, rid, f.JobType, pipeline, stage, f.JobName, ts.Format(time.RFC3339), b.Result, b.Duration)
					Logger.Printf("Unable to find a parent upstream job - job started manually ?\n")
				}
			} else {
				Logger.Printf("Build still running, skipping\n")
			}
		}
	}
	Logger.Printf("==== itba stop ====\n")

}

//func getStage(jobName string) string {
//	stg := "UNK"
//	for _, js := range Conf.JobStage {
//		if strings.Contains(jobName, js.Suffix) {
//			stg = js.Stage
//			break
//		}
//	}
//	return stg
//}