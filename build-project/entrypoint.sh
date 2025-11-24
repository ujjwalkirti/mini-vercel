#!/bin/bash
set -e

# Load NVM
export NVM_DIR=/root/.nvm
source $NVM_DIR/nvm.sh

# Execute your main script
exec /home/app/main.sh
