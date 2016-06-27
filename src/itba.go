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
	"log"
	"os"
	"time"
)

var logger log.Logger
var conf *Config

func init() {
	logger = *log.New(os.Stderr, "itba: ", log.LstdFlags)
}

func main() {
	fmt.Print("first line of code !")
	var URL string
	var username string
	var password string
	var cfile string

	flag.StringVar(&URL, "URL", "", "Jenkins server URL")
	flag.StringVar(&username, "u", "", "username to connect to jenkins")
	flag.StringVar(&password, "p", "", "password to connect to jenkins")
	flag.StringVar(&cfile, "config", "config.yaml", "abs path of config file")

	flag.Parse()

	if len(os.Args) == 1 {
		flag.PrintDefaults()
		return
	}
	logger.Printf("==== itba start ====\n")

	conf = &Config{}
	loadConfig(conf, cfile)
	logger.Printf("Loaded [%d] job filters\n", len(conf.Filters))
	logger.Printf("Loaded [%d] job pipeline mappings \n", len(conf.JobPipelines))
	logger.Printf("Loaded [%d] job stage mappings \n", len(conf.JobStage))

	jk := NewJenkins(URL, username, password)
	fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,%s\n", "JOB_ID", "REL_ID", "JOB_TYPE", "PIPELINE", "STAGE", "JOB_NAME", "START_TIME", "JOB_RESULT", "DURATION")
	for _, f := range conf.Filters {
		logger.Printf("Getting details for job: %s\n", f.JobName)
		jd := jk.JobDetails(f.JobName)
		for _, b := range jd.Builds {
			logger.Printf("Processing build %s-%d\n", f.JobName, b.Number)
			ts := time.Unix(0, b.Timestamp*int64(time.Millisecond))
			if b.Result != "" {
				logger.Printf("Getting parent upstream job %s-%d\n", f.JobName, b.Number)
				rid, pipeline := GetUpstreamJob(jk, f.JobName, b.Number) // TODO should get "pipeline" when rid==0
				if rid != 0 {
					stage := b.Stage()
					fmt.Printf("%d,%d,%s,%s,%s,%s,%s,%s,%d\n", b.Number, rid, f.JobType, pipeline, stage, f.JobName, ts.Format(time.RFC3339), b.Result, b.Duration)
					logger.Printf("%d,%d,%s,%s,%s,%s,%s,%s,%d\n", b.Number, rid, f.JobType, pipeline, stage, f.JobName, ts.Format(time.RFC3339), b.Result, b.Duration)
				} else {
					stage := b.Stage()
					fmt.Printf("%d,%d,%s,%s,%s,%s,%s,%s,%d\n", b.Number, rid, f.JobType, pipeline, stage, f.JobName, ts.Format(time.RFC3339), b.Result, b.Duration)

					logger.Printf("Unable to find a parent upstream job - job started manually ?\n")
				}
			} else {
				logger.Printf("Build still running, skipping\n")
			}
		}
	}
	logger.Printf("==== itba stop ====\n")

}
