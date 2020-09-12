#!/bin/bash

# This script enrolls all the admin users

# Exit on first error
set -e

ORG1=Org1
ORG2=Org2
ORG3=Org3
ORG4=Org4
ORG5=Org5
ORG6=Org6


node enrollAdmin.js ${ORG1}
node enrollAdmin.js ${ORG2}
node enrollAdmin.js ${ORG3}
node enrollAdmin.js ${ORG4}
node enrollAdmin.js ${ORG5}
node enrollAdmin.js ${ORG6}

echo '============= All admin users were registered successfully ============='
echo "Org1:" ${ORG1}
echo "Org2:" ${ORG2}
echo "Org3:" ${ORG3}
echo "Org4:" ${ORG4}
echo "Org5:" ${ORG5}
echo "Org6:" ${ORG6}
