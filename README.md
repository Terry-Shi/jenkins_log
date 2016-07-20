# jenkins_log
Jenkins build history extractor

#### How to install gopkg.in/yaml.v2; what's gopkg?
[what's gopkg ?](http://blog.csdn.net/siddontang/article/details/38083159)

1. setup GOPATH for mac OS in ~/.bashrc
    `export GOPATH=/Users/terry/weiyun/go`
    `export PATH=$GOPATH/bin:$PATH`
    make above change effective: `source ~/.bashrc`

2. setup GOPATH for Windows
   add GOPATH at "Environment Variables" -> "System variables"

    set proxy in windows if needed
    `set http_proxy=http://web-proxy.corp.hpecorp.net:8080`
    `set https_proxy=http://web-proxy.corp.hpecorp.net:8080`

3. install yaml `go get gopkg.in/yaml.v2`

#### how to setup GOPATH in IntelliJ IDEA
Ref:
https://rootpd.com/2016/02/04/setting-up-intellij-idea-for-your-first-golang-project/ 
http://stackoverflow.com/questions/17771091/i-use-intellij-idea-as-golang-ide-and-system-environment-have-already-set-gopat

#### How to run itba.go (when it depends on other source file), and they are all in main package
error msg:

    C:/Go\bin\go.exe run C:/WORK/IDE/IdeaProjects/jenkins_log/itba.go
    # command-line-arguments
    .\itba.go:51: undefined: Config
    .\itba.go:77: undefined: Config
    .\itba.go:78: undefined: loadConfig
    .\itba.go:83: undefined: NewJenkins
    .\itba.go:93: undefined: GetUpstreamJob
    
    Process finished with exit code 2
Ref: 
http://ami-gs.hatenablog.com/entry/2014/12/30/181445
http://stackoverflow.com/questions/26920969/why-does-the-go-run-command-fail-to-find-a-second-file-in-the-main-package
https://github.com/go-lang-plugin-org/go-lang-idea-plugin/issues/2012
https://golang.org/doc/code.html#Workspaces
https://groups.google.com/forum/#!topic/golang-nuts/urv3eP6ILaU

#### misc
- what's the meaning of all the GOXXXX env variable ? such as GOROOT, GOPATH
- check env list with `go env`

#### TODO:
1. 让程序跑起来!
   1.1 HTTPS URL的访问 http://stackoverflow.com/questions/12122159/golang-how-to-do-a-https-request-with-bad-certificate
2. 修改yaml文件,jobType and stage
3. command:
- Windows
   cd C:/WORK/IDE/IdeaProjects/jenkins_log/src
   go run itba.go -URL=https://cas-cd.core.hpecorp.net:2543 -u=1 -p=1 -config=C:\WORK\IDE\IdeaProjects\jenkins_log\config.yaml 1> JENKINS_JOB_CSTM.csv 2> itba.log
   
- Mac
   go run itba.go -URL=https://cas-cd.core.hpecorp.net:2543 -u=1 -p=1 -config=/Users/terry/IdeaProjects/jenkins_log/config.yaml 1> JENKINS_JOB_CSTM.csv 2> itba.log

4. import csv file
   ssh -i ~/Desktop/HPIT.ppk  hos@16.202.70.65
   https://itba.itcs.hpe.com:8443/bsf/login.form