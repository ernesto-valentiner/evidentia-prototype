/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Gateway, Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');

async function main() {
    var args = process.argv.slice(2);
    if(args.length < 1) {
        console.log('Wrong number of arguments - Username of coordinator is missing')
        process.exit(1);
    }
    if(args.length > 1) {
        console.log('Wrong number of arguments - Requires 1 argument')
        process.exit(1);
    }

    var username = args[0]
    try {
        // load the network configuration
        const ccpPath = path.resolve(__dirname, '..', '..', 'test-network', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
        let ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get(username);
        if (!identity) {
            console.log(`An identity for the user "${username}" does not exist in the wallet`);
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: username, discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('fgs-channel');

        // Get the contract from the network.
        const contract = network.getContract('evidentia');

        //Subscribe to serviceExecution event
        await contract.addContractListener('my-contract-listener', 'serviceExecution', async(err, event, blockNumber, transactionId, status)  => {
            if (err) {
                console.error(err);
                return;
            }

            //convert event to something we can parse
            const service = JSON.parse(event.payload.toString());

            console.log('************************ Service Event *******************************************************');
            console.log(`Service name: ${service.serviceName}`);
            console.log(`Params: ${service.params}`);
            console.log(`Target User: ${service.target}`);
            if (services.includes(service.serviceName)) {
                console.log(`I will execute the service`);
                console.log(`Updating Target for ${service.serviceName}${service.params}`);
                await contract.submitTransaction('updateServiceExecutionTarget', service.serviceName, service.params);
                console.log(`Updated Target for ${service.serviceName}${service.params}`);
                console.log(`Updating Response for ${service.serviceName}${service.params}`);
                await contract.submitTransaction('updateServiceExecutionResponse', service.serviceName, service.params, "EVIDENCE", "RESPONSE");
                console.log(`Updated Response for ${service.serviceName}${service.params}`);
            } else {
                console.log(`I can't execute the service`);
            }
        });

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
