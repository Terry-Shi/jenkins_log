#! /bin/bash
JENKINS_URL=https://cas-cd.core.hpecorp.net:2543
JENKINS_USERNAME=1
JENKINS_PASSWORD=1
# ITBA URL is https://itba.itcs.hpe.com:8443/bsf/login.form
ITBA_HOST=16.202.70.65
ITBA_USER=hos
ITBA_DATA_DIR=/opt/itba/OV/InputSourceSH
BASE_DIR=/Users/terry/IdeaProjects/jenkins_log

$BASE_DIR/itba -URL=$JENKINS_URL -u=$JENKINS_USERNAME -p=$JENKINS_PASSWORD -config=$BASE_DIR/config.yaml 1> $BASE_DIR/itba.csv.new 2> $BASE_DIR/itba.log
mv $BASE_DIR/itba.csv.new $BASE_DIR/itba.csv
scp -Cp $BASE_DIR/itba.csv $ITBA_USER@$ITBA_HOST:$ITBA_DATA_DIR/JENKINS_JOB_CSTM.csv

