[//]: # (SPDX-License-Identifier: CC-BY-4.0)

# Distributed Evidence Network

## Network Model

![alt text](/images/Prototype.png)

## Getting Started

### Prerequisites

Check that all the [prerequisites](https://hyperledger-fabric.readthedocs.io/en/release-2.0/prereqs.html) are installed
on your machine in order to run a Hyperledger Fabric Network.

### Install Binaries and Docker Images

Enter the following line to install all the binaries and Docker Images.

```console
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.0.1 1.4.6 0.4.18
```
More information can be found [here](https://hyperledger-fabric.readthedocs.io/en/release-2.0/install.html)

### Installation

Clone the repository

```console
git clone https://github.com/ernesto-valentiner/evidentia-prototype.git
```

## Usage

### Bring the network up

After installing all the prerequisites and binaries, start the network by changing directory to evidentia-app 
and executing ./startFabric.sh

```console
cd evidentia-app
./startFabric.sh
```
The network might take some time to build.
When the following message is shown in the console, the network has finished building.

```console 
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
``` 

### Bring the network down

To bring the network down execute the following command in the evidentia-app directory.
```console
./networkDown.sh
```

### Next steps

To connect the ETBNodes with the network, follow the description in the [ETB repository](https://github.com/ernesto-valentiner/ETB).

## License <a name="license"></a>

Hyperledger Project source code files are made available under the Apache
License, Version 2.0 (Apache-2.0), located in the [LICENSE](LICENSE) file.
Hyperledger Project documentation files are made available under the Creative
Commons Attribution 4.0 International License (CC-BY-4.0), available at http://creativecommons.org/licenses/by/4.0/.
