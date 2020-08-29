#!/bin/bash

# This script registers all users

# Exit on first error
set -e

ORG1=Org1
ORG2=Org2
ORG3=Org3
ORG4=Org4
ORG5=Org5
ORG6=Org6

ORG1_USER=user_org1
ORG2_USER=user_org2
ORG3_USER=user_org3
ORG4_USER=user_org4
ORG5_USER=user_org5
ORG6_USER=user_org6


node registerUser.js ${ORG1} ${ORG1_USER}
node registerUser.js ${ORG2} ${ORG2_USER}
node registerUser.js ${ORG3} ${ORG3_USER}
node registerUser.js ${ORG4} ${ORG4_USER}
node registerUser.js ${ORG5} ${ORG5_USER}
node registerUser.js ${ORG6} ${ORG6_USER}


echo '============= All users are registered successfully ============='
echo "Org1:" ${ORG1_USER}
echo "Org2:" ${ORG2_USER}
echo "Org3:" ${ORG3_USER}
echo "Org4:" ${ORG4_USER}
echo "Org5:" ${ORG5_USER}
echo "Org5:" ${ORG6_USER}
