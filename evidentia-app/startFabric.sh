#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
starttime=$(date +%s)
CC_SRC_LANGUAGE=${1:-"go"}
CC_SRC_LANGUAGE=`echo "$CC_SRC_LANGUAGE" | tr [:upper:] [:lower:]`
if [ "$CC_SRC_LANGUAGE" != "go" -a "$CC_SRC_LANGUAGE" != "golang" -a "$CC_SRC_LANGUAGE" != "java" \
 -a  "$CC_SRC_LANGUAGE" != "javascript"  -a "$CC_SRC_LANGUAGE" != "typescript" ] ; then

	echo The chaincode language ${CC_SRC_LANGUAGE} is not supported by this script
 	echo Supported chaincode languages are: go, java, javascript, and typescript
 	exit 1

fi

# clean out any old identites in the wallets
rm -rf javascript/wallet/*
rm -rf java/wallet/*
rm -rf typescript/wallet/*

# launch network; create channel and join peer to channel
pushd ../test-network
./network.sh down
./network.sh up createChannel -ca -s couchdb
./network.sh deployCC -l ${CC_SRC_LANGUAGE}
popd

cat <<EOF

Total setup execution time : $(($(date +%s) - starttime)) secs ...

Next, setup the ETBNodes to connect to the network and interact with the Evidentia contract.
The javascript application under evidentia-app/javascript can also be used to interact with the network separately from the ETBNodes.


JavaScript Application:

  Start by changing into the "javascript" directory:
    cd javascript

  Next, install all required packages:
    npm install

  Then run the following applications to enroll the admin users, and register new users
  which will be used by the other applications to interact with the deployed
  Evidentia contract:
    To enroll and register one admin and one user for each organization:
    ./enrollAllAdmins.sh
    ./registerAllUsers.sh

    To enroll an admin for a specific organization(e.g. organization: Org1):
    node enrollAdmin Org1

    To register a user for a specific organization(e.g. organization: Org1, Username: user1):
    node registerUser Org1 user1

  You can run the invoke and query application by updating the application with the correct function calls
  and then executing them as followed:
    node invoke
    node query

EOF
