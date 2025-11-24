#!/bin/bash

export GIT_REPOSITORY_URL="${GIT_REPOSITORY_URL}"

# clone the repo
git clone "$GIT_REPOSITORY_URL" /home/app/output

# run the js script
exec node script.js

