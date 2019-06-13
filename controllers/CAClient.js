/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const FabricCAServices = require('fabric-ca-client');
const { FileSystemWallet, X509WalletMixin, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve(__dirname, '..', 'config', 'connection-org1.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

// Create a new file system based wallet for managing identities.
const walletPath = path.join(process.cwd(), 'wallet');
const wallet = new FileSystemWallet(walletPath);

class CAClient {

    /**
     * Enroll admin user
     */
    async enrollAdmin() {
        try {
            // Create a new CA client for interacting with the CA.
            const caInfo = ccp.certificateAuthorities['ca.example.com'];
            const caTLSCACertsPath = path.resolve(__dirname, '..', caInfo.tlsCACerts.path);
            const caTLSCACerts = fs.readFileSync(caTLSCACertsPath);
            const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

            console.log(`Wallet path: ${walletPath}`);

            // Check to see if we've already enrolled the admin user.
            const adminExists = await wallet.exists('admin');
            if (adminExists) {
                return {
                    status: 422,
                    data: 'An identity for the admin user "admin" already exists in the wallet'
                };
            }

            // Enroll the admin user, and import the new identity into the wallet.
            const enrollment = await ca.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' });
            const identity = X509WalletMixin.createIdentity('Org1MSP', enrollment.certificate, enrollment.key.toBytes());
            await wallet.import('admin', identity);
            return {
                status: 200,
                data: 'Successfully enrolled admin user "admin" and imported it into the wallet'
            };

        } catch (error) {
            return {
                status: 500,
                data: `Failed to enroll admin user "admin": ${error}`
            };
        }
    }

    /**
     * Register & enroll user with CA
     */
    async enrollUser(username) {
        try {
            console.log(username);
            console.log(`Wallet path: ${walletPath}`);

            // Check to see if we've already enrolled the user.
            const userExists = await wallet.exists(username);
            if (userExists) {
                return {
                    status: 422,
                    data: 'An identity for the user "' + username + '" already exists in the wallet'
                };
            }

            // Check to see if we've already enrolled the admin user.
            const adminExists = await wallet.exists('admin');
            if (!adminExists) {
                return {
                    status: 422,
                    data: 'An identity for the admin user "admin" does not exist in the wallet'
                };
            }

            // Create a new gateway for connecting to our peer node.
            const gateway = new Gateway();
            await gateway.connect(ccpPath, { wallet, identity: 'admin', discovery: { enabled: true, asLocalhost: true } });

            // Get the CA client object from the gateway for interacting with the CA.
            const ca = gateway.getClient().getCertificateAuthority();
            const adminIdentity = gateway.getCurrentIdentity();

            // Register the user, enroll the user, and import the new identity into the wallet.
            const secret = await ca.register({ affiliation: 'org1.department1', enrollmentID: username, role: 'client' }, adminIdentity);
            const enrollment = await ca.enroll({ enrollmentID: username, enrollmentSecret: secret });
            const userIdentity = X509WalletMixin.createIdentity('Org1MSP', enrollment.certificate, enrollment.key.toBytes());
            await wallet.import(username, userIdentity);
            return {
                status: 200,
                data: 'Successfully registered and enrolled admin user "' + username + '" and imported it into the wallet'
            };

        } catch (error) {
            return {
                status: 500,
                data: `Failed to register user "` + username + `": ${error}`
            };
        }
    }
}
module.exports = CAClient;