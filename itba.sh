#! /bin/bash
JENKINS_URL=https://cas-cd.core.hpecorp.net:2543
JENKINS_USERNAME=
JENKINS_PASSWORD=
ITBA_HOST=10.10.3.144
ITBA_USER=root
ITBA_DATA_DIR=/home/admin/InputData/Jenkins
BASE_DIR=/Users/terry/IdeaProjects/jenkins_log

$BASE_DIR/itba -URL=$JENKINS_URL -u=$JENKINS_USERNAME -p=$JENKINS_PASSWORD -config=$BASE_DIR/config.yaml 1> $BASE_DIR/itba.csv.new 2> $BASE_DIR/itba.log
mv $BASE_DIR/itba.csv.new $BASE_DIR/itba.csv
scp -Cp $BASE_DIR/itba.csv $ITBA_USER@$ITBA_HOST:$ITBA_DATA_DIR/JENKINS_JOB_CSTM.csv

