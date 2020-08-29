/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const FabricCAServices = require('fabric-ca-client');
const { Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');

async function main() {
    var args = process.argv.slice(2);
    if(args.length < 1) {
        console.log('Wrong number of arguments - Organisation name is missing')
        process.exit(1);
    }
    if(args.length >= 2) {
        console.log('Wrong number of arguments - Requires 1 argument')
        process.exit(1);
    }
    var Org = args[0];
    var org = args[0].toLowerCase();
    try {
        // load the network configuration
        const ccpPath = path.resolve(__dirname, '..', '..', 'test-network', 'organizations', 'peerOrganizations', `${org}.example.com`, `connection-${org}.json`);
        const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

        // Create a new CA client for interacting with the CA.
        const caInfo = ccp.certificateAuthorities[`ca.${org}.example.com`];
        const caTLSCACerts = caInfo.tlsCACerts.pem;
        const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the admin user.
        const identity = await wallet.get(`admin_${org}`);
        if (identity) {
            console.log(`An identity for the admin user "admin_${org}" already exists in the wallet`);
            return;
        }

        // Enroll the admin user, and import the new identity into the wallet.
        const enrollment = await ca.enroll({ enrollmentID: `admin`, enrollmentSecret: 'adminpw' });
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: `${Org}MSP`,
            type: 'X.509',
        };
        await wallet.put(`admin_${org}`, x509Identity);
        console.log(`Successfully enrolled admin user "admin_${org}" and imported it into the wallet`);

    } catch (error) {
        console.error(`Failed to enroll admin user "admin_${org}": ${error}`);
        process.exit(1);
    }
}

main();
