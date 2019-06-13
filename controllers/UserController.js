'use strict';

const config = require('../config/config.js');
const FabricOperation = require('../controllers/FabricOperation.js');
const FabricController = new FabricOperation();
const CAClient = require('../controllers/CAClient.js');
const CAClientController = new CAClient();

class UserController {

    async getUser(data) {
        try {
            let result = await CAClientController.enrollUser(data.user);
            if (result.status = 200) {
                // Invoke the chaincode function
                let response = await FabricController.invoke(data.user, config.channel, config.chaincode, 'getUser', data);
                return response;
            }
            else {
                return {
                    status: result.status,
                    data: {
                        data: ''
                    }
                };
            }
        }
        catch (error) {
            return {
                status: 500,
                data: {
                    data: `Catch block error ": ${error}`
                }
            };
        }
    }

    async addAddress(data) {
        // Invoke the chaincode function
        let response = await FabricController.invoke(data.user, config.channel, config.chaincode, 'addAddress', data);
        return response;
    }

    async allUsers(data) {
        // Invoke the chaincode function
        let response = await FabricController.invoke(data.user, config.channel, config.chaincode, 'getUsers', data);
        return response;
    }

    async issueToken(data) {
        // Invoke the chaincode function
        let response = await FabricController.invoke(data.user, config.channel, config.chaincode, 'addAsset', data);
        return response;
    }

    async checkSymbol(data) {
        // Invoke the chaincode function
        let response = await FabricController.invoke(data.user, config.channel, config.chaincode, 'checkAsset', data);
        return response;
    }

    async getToken(data) {
        // Invoke the chaincode function
        let response = await FabricController.invoke(data.user, config.channel, config.chaincode, 'getAssets', data);
        return response;
    }

    async transferToken(data) {
        // Invoke the chaincode function
        let response = await FabricController.invoke(data.user, config.channel, config.chaincode, 'transferAsset', data);
        return response;
    }
}

module.exports = UserController;