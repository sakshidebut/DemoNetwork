'use strict';

const config = require('../config/config.js');
const FabricCAServices = require('fabric-ca-client');
const { FileSystemWallet, Gateway, X509WalletMixin } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve(__dirname, '..', 'config', 'connection-org1.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);
// Create a new file system based wallet for managing identities.
const walletPath = path.join(process.cwd(), 'wallet');
const wallet = new FileSystemWallet(walletPath);

class UserController {

    async enrollAdmin() {
        try {
            // Create a new CA client for interacting with the CA.
            const caInfo = ccp.certificateAuthorities['ca.example.com'];
            const caTLSCACertsPath = path.resolve(__dirname, caInfo.tlsCACerts.path);
            const caTLSCACerts = fs.readFileSync(caTLSCACertsPath);
            const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

            console.log(`Wallet path: ${walletPath}`);

            // Check to see if we've already enrolled the admin user.
            const adminExists = await wallet.exists('admin');
            if (adminExists) {
                return 'An identity for the admin user "admin" already exists in the wallet';
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

    async enrollUser() {
        try {
            console.log(`Wallet path: ${walletPath}`);

            // Check to see if we've already enrolled the user.
            const userExists = await wallet.exists(config.user);
            if (userExists) {
                return 'An identity for the user "' + config.user + '" already exists in the wallet';
            }

            // Check to see if we've already enrolled the admin user.
            const adminExists = await wallet.exists('admin');
            if (!adminExists) {
                return 'An identity for the admin user "admin" does not exist in the wallet';
            }

            // Create a new gateway for connecting to our peer node.
            const gateway = new Gateway();
            await gateway.connect(ccpPath, { wallet, identity: 'admin', discovery: { enabled: true, asLocalhost: true } });

            // Get the CA client object from the gateway for interacting with the CA.
            const ca = gateway.getClient().getCertificateAuthority();
            const adminIdentity = gateway.getCurrentIdentity();

            // Register the user, enroll the user, and import the new identity into the wallet.
            const secret = await ca.register({ affiliation: 'org1.department1', enrollmentID: config.user, role: 'client' }, adminIdentity);
            const enrollment = await ca.enroll({ enrollmentID: config.user, enrollmentSecret: secret });
            const userIdentity = X509WalletMixin.createIdentity('Org1MSP', enrollment.certificate, enrollment.key.toBytes());
            await wallet.import(config.user, userIdentity);
            return {
                status: 200,
                data: 'Successfully registered and enrolled admin user "' + config.user + '" and imported it into the wallet'
            };

        } catch (error) {
            return {
                status: 500,
                data: `Failed to register user "` + config.user + `": ${error}`
            };
        }
    }

    async createUser(data) {
        try {
            console.log(`Wallet path: ${walletPath}`);

            // Check to see if we've already enrolled the user.
            const userExists = await wallet.exists(config.user);
            if (!userExists) {
                return 'An identity for the user "' + config.user + '" does not exist in the wallet';
            }

            // Create a new gateway for connecting to our peer node.
            const gateway = new Gateway();
            await gateway.connect(ccpPath, { wallet, identity: config.user, discovery: { enabled: true, asLocalhost: true } });

            // Get the network (channel) our contract is deployed to.
            const network = await gateway.getNetwork(config.channel);

            // Get the contract from the network.
            const contract = network.getContract(config.chaincode);
            let result;
            // Submit the specified transaction.
            // createUser transaction - requires 4 arguments
            result = await contract.submitTransaction('createUser', data.name, data.email, data.phone, data.address);
            // Disconnect from the gateway.
            await gateway.disconnect();
            return {
                status: 200,
                data: {
                    data: JSON.parse(result.toString())
                }
            };
        } catch (error) {
            return this.handleError(error)
        }
    }

    async getUser(data) {
        try {
            console.log(`Wallet path: ${walletPath}`);

            // Check to see if we've already enrolled the user.
            const userExists = await wallet.exists(config.user);
            if (!userExists) {
                return 'An identity for the user "' + config.user + '" does not exist in the wallet';
            }

            // Create a new gateway for connecting to our peer node.
            const gateway = new Gateway();
            await gateway.connect(ccpPath, { wallet, identity: config.user, discovery: { enabled: true, asLocalhost: true } });

            // Get the network (channel) our contract is deployed to.
            const network = await gateway.getNetwork(config.channel);

            // Get the contract from the network.
            const contract = network.getContract(config.chaincode);

            let result;
            // Submit the specified transaction.
            // getUser transaction - requires 1 arguments
            result = await contract.submitTransaction('getUser', data.id);
            // Disconnect from the gateway.
            await gateway.disconnect();
            return {
                status: 200,
                data: {
                    data: JSON.parse(result.toString())
                }
            };
        } catch (error) {
            return this.handleError(error)
        }
    }

    async allUsers(data) {
        try {
            console.log(`Wallet path: ${walletPath}`);

            // Check to see if we've already enrolled the user.
            const userExists = await wallet.exists(config.user);
            if (!userExists) {
                return 'An identity for the user "' + config.user + '" does not exist in the wallet';
            }

            // Create a new gateway for connecting to our peer node.
            const gateway = new Gateway();
            await gateway.connect(ccpPath, { wallet, identity: config.user, discovery: { enabled: true, asLocalhost: true } });

            // Get the network (channel) our contract is deployed to.
            const network = await gateway.getNetwork(config.channel);

            // Get the contract from the network.
            const contract = network.getContract(config.chaincode);

            let result;
            // Submit the specified transaction.
            // getUser transaction - requires 1 arguments
            result = await contract.submitTransaction('getUsers', data.id);
            // Disconnect from the gateway.
            await gateway.disconnect();
            return {
                status: 200,
                data: {
                    data: JSON.parse(result.toString())
                }
            };
        } catch (error) {
            return this.handleError(error)
        }
    }

    async issueToken(data) {
        try {
            console.log(`Wallet path: ${walletPath}`);

            // Check to see if we've already enrolled the user.
            const userExists = await wallet.exists(config.user);
            if (!userExists) {
                return 'An identity for the user "' + config.user + '" does not exist in the wallet';
            }

            // Create a new gateway for connecting to our peer node.
            const gateway = new Gateway();
            await gateway.connect(ccpPath, { wallet, identity: config.user, discovery: { enabled: true, asLocalhost: true } });

            // Get the network (channel) our contract is deployed to.
            const network = await gateway.getNetwork(config.channel);

            // Get the contract from the network.
            const contract = network.getContract(config.chaincode);

            let result;
            // Submit the specified transaction.
            // getUser transaction - requires 1 arguments
            result = await contract.submitTransaction('addAsset', data.id, data.code, data.quantity);
            // Disconnect from the gateway.
            await gateway.disconnect();
            return {
                status: 200,
                data: {
                    data: {
                        message: 'Asset has been saved successfully.'
                    }
                }
            };
        } catch (error) {
            return this.handleError(error)
        }
    }

    async getToken(data) {
        try {
            console.log(`Wallet path: ${walletPath}`);

            // Check to see if we've already enrolled the user.
            const userExists = await wallet.exists(config.user);
            if (!userExists) {
                return 'An identity for the user "' + config.user + '" does not exist in the wallet';
            }

            // Create a new gateway for connecting to our peer node.
            const gateway = new Gateway();
            await gateway.connect(ccpPath, { wallet, identity: config.user, discovery: { enabled: true, asLocalhost: true } });

            // Get the network (channel) our contract is deployed to.
            const network = await gateway.getNetwork(config.channel);

            // Get the contract from the network.
            const contract = network.getContract(config.chaincode);

            let result;
            // Submit the specified transaction.
            // getUser transaction - requires 1 arguments
            result = await contract.submitTransaction('getAssets', data.id);
            // Disconnect from the gateway.
            await gateway.disconnect();
            return {
                status: 200,
                data: {
                    data: JSON.parse(result.toString())
                }
            };

        } catch (error) {
            return this.handleError(error)
        }
    }

    async transferToken(data) {
        try {
            console.log(`Wallet path: ${walletPath}`);

            // Check to see if we've already enrolled the user.
            const userExists = await wallet.exists(config.user);
            if (!userExists) {
                return 'An identity for the user "' + config.user + '" does not exist in the wallet';
            }

            // Create a new gateway for connecting to our peer node.
            const gateway = new Gateway();
            await gateway.connect(ccpPath, { wallet, identity: config.user, discovery: { enabled: true, asLocalhost: true } });

            // Get the network (channel) our contract is deployed to.
            const network = await gateway.getNetwork(config.channel);

            // Get the contract from the network.
            const contract = network.getContract(config.chaincode);

            let result;
            // Submit the specified transaction.
            // getUser transaction - requires 1 arguments
            result = await contract.submitTransaction('transferAsset', data.from_id, data.to_id, data.code, data.quantity);
            // Disconnect from the gateway.
            await gateway.disconnect();
            return {
                status: 200,
                data: {
                    data: {
                        message: "Assets has been trasnfered successfully!"
                    }
                }
            };
        } catch (error) {
            return this.handleError(error)
        }
    }

    /**
     * Checks whether a string is JSON or not
     * @param {*} item 
     */
    isJson(item) {
        item = typeof item !== "string"
            ? JSON.stringify(item)
            : item;
        try {
            item = JSON.parse(item);
        } catch (e) {
            return false;
        }
        if (typeof item === "object" && item !== null) {
            return item;
        }
        return false;
    }

    /**
     * Handles the errors occured during invoke or query
     * @param {*} error 
     */
    handleError(error) {
        let response = {
            status: 500,
            data: {
                message: error.message
            }
        };
        // check for the chaincode response
        if (error.hasOwnProperty('endorsements')) {
            // chaincode is executed and has thrown some error
            let endorsements = error.endorsements;
            if (endorsements.length) {
                // get the details of the error
                let errors = this.isJson(endorsements[0].message);
                if (errors) {
                    // make the response
                    response.data.message = errors.msg;
                    response.status = errors.code;
                }
                else {
                    response.data.message = endorsements[0].message;
                }

                // check if the error has extra data
                if (errors.hasOwnProperty('details')) {
                    response.data.errors = errors.details;
                }
            }
            return response;
        }
    }
}

module.exports = UserController;