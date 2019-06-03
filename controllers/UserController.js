'use strict';

const config = require('../config/config.js');
const FabricOperation = require('../controllers/FabricOperation.js');
const FabricController = new FabricOperation();

class UserController {

    async createUser(data) {
        // Invoke the chaincode function
        let response = await FabricController.invoke(config.user, config.channel, config.chaincode, 'createUser', data);
        return response;
    }

    async getUser(data) {
        // Invoke the chaincode function
        let response = await FabricController.invoke(config.user, config.channel, config.chaincode, 'getUser', data);
        return response;
    }

    async allUsers(data) {
        // Invoke the chaincode function
        let response = await FabricController.invoke(config.user, config.channel, config.chaincode, 'getUsers', data);
        return response;
    }

    async issueToken(data) {
        // Invoke the chaincode function
        let response = await FabricController.invoke(config.user, config.channel, config.chaincode, 'addAsset', data);
        return response;
    }

    async getToken(data) {
        // Invoke the chaincode function
        let response = await FabricController.invoke(config.user, config.channel, config.chaincode, 'getAssets', data);
        return response;
    }

    async transferToken(data) {
        // Invoke the chaincode function
        let response = await FabricController.invoke(config.user, config.channel, config.chaincode, 'transferAsset', data);
        return response;
    }
}

module.exports = UserController;