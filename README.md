# jenkins_log
Jenkins build history extractor

#### How to install gopkg.in/yaml.v2; what's gopkg?
[what's gopkg ?](http://blog.csdn.net/siddontang/article/details/38083159)

1. setup GOPATH in ~/.bashrc
    export GOPATH=/Users/terry/weiyun/go
    export PATH=$GOPATH/bin:$PATH
2. make above change effective: `source ~/.bashrc`
3. install yaml `go get gopkg.in/yaml.v2`

#### how to setup GOPATH in IntelliJ IDEA
Ref:
https://rootpd.com/2016/02/04/setting-up-intellij-idea-for-your-first-golang-project/ 
http://stackoverflow.com/questions/17771091/i-use-intellij-idea-as-golang-ide-and-system-environment-have-already-set-gopat


#### misc
- what's the meaning of all the GOXXXX env variable ? such as GOROOT, GOPATH
- check env list with `go env`


#### TODO:
1. 让程序跑起来!
   1.1 HTTPS URL的访问 http://stackoverflow.com/questions/12122159/golang-how-to-do-a-https-request-with-bad-certificate
2. 修改yaml文件,jobType and stage
3. command:
   go run itba.go config.go jenkins.go -URL=https://cas-cd.core.hpecorp.net:2543 -u= -p= -config=/Users/terry/IdeaProjects/jenkins_log/config.yaml 1> itba.csv.new 2> itba.log