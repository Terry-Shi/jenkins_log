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


####
what's the meaning of all the GOXXXX env variable ? such as GOROOT, GOPATH
check env list with `go env`