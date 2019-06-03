'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve(__dirname, '..', 'config', 'connection-org1.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

// Create a new file system based wallet for managing identities.
const walletPath = path.join(process.cwd(), 'wallet');
const wallet = new FileSystemWallet(walletPath);

class FabricOperation {

    /**
     * @api Query Transaction
     * @param string user username
     * @param string channel_name channel name to query
     * @param string chaincode_name  chaincode name
     * @param string function_name function name
     * @param json data data for query
     */
    async query(user, channel_name, chaincode_name, function_name, data = null) {
        try {
            this.userExists(user);

            // Create a new gateway for connecting to our peer node.
            const gateway = new Gateway();
            await gateway.connect(ccp, { wallet, identity: user, discovery: { enabled: false } });

            // Get the network (channel) our contract is deployed to.
            const network = await gateway.getNetwork(channel_name);

            // Get the contract from the network.
            const contract = network.getContract(chaincode_name);

            let result;

            // Submit the specified transaction.
            if (data) {
                result = await contract.evaluateTransaction(function_name, JSON.stringify(data));
            }
            else {
                result = await contract.evaluateTransaction(function_name);
            }

            // Disconnect from the gateway.
            await gateway.disconnect();
            return {
                status: 200,
                data: {
                    data: JSON.parse(result.toString())
                }
            };
        } catch (error) {
            return this.handleError(error);
        }
    }

    /**
     * @api Invoke Transaction
     * @param string user username
     * @param string channel_name channel name to query
     * @param string chaincode_name  chaincode name
     * @param string function_name function name
     * @param json data data for query
     */
    async invoke(user, channel_name, chaincode_name, function_name, data) {
        try {
            this.userExists(user);
            // Create a new gateway for connecting to our peer node.
            const gateway = new Gateway();
            await gateway.connect(ccp, { wallet, identity: user, discovery: { enabled: true } });

            // Get the network (channel) our contract is deployed to.
            const network = await gateway.getNetwork(channel_name);

            // Get the contract from the network.
            const contract = network.getContract(chaincode_name);

            // Submit the specified transaction.
            let result;
            result = await contract.submitTransaction(function_name, JSON.stringify(data));
            // Disconnect from the gateway.
            await gateway.disconnect();
            return {
                status: 200,
                data: {
                    data: JSON.parse(result.toString())
                }
            };
        } catch (error) {
            return this.handleError(error);
        }
    }

    /**
     * Check the user exists or not
     * @param string user username to be checked
     */
    async userExists(user) {
        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists(user);
        if (!userExists) {
            throw new Error(`Please enroll: ${user}`);
        }
        return true;
    }

    /**
     * Checks whether a string is JSON or not
     * @param {*} item
     */
    isJson(item) {
        item = typeof item !== 'string'
            ? JSON.stringify(item)
            : item;
        try {
            item = JSON.parse(item);
        } catch (e) {
            return false;
        }
        if (typeof item === 'object' && item !== null) {
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
        }
        return response;
    }
}

module.exports = FabricOperation;