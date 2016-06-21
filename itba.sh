#! /bin/bash
JENKINS_URL=http://10.10.2.46:8080/jenkins
JENKINS_USERNAME=devops
JENKINS_PASSWORD=devops
ITBA_HOST=10.10.3.144
ITBA_USER=root
ITBA_DATA_DIR=/home/admin/InputData/Jenkins
BASE_DIR=/root/itba

$BASE_DIR/itba -URL=$JENKINS_URL -u=$JENKINS_USERNAME -p=$JENKINS_PASSWORD -config=$BASE_DIR/config.yaml 1> $BASE_DIR/itba.csv.new 2> $BASE_DIR/itba.log
mv $BASE_DIR/itba.csv.new $BASE_DIR/itba.csv
scp -Cp $BASE_DIR/itba.csv $ITBA_USER@$ITBA_HOST:$ITBA_DATA_DIR/JENKINS_JOB_CSTM.csv

