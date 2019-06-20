'use strict';

const config = require('../config/config.js');
const FabricOperation = require('../controllers/FabricOperation.js');
const FabricController = new FabricOperation();
const CAClient = require('../controllers/CAClient.js');
const CAClientController = new CAClient();

class UserController {

    async createUser(data) {
        try {
            const secret = this.makeid(12);
            let result = await CAClientController.enrollUser(data.user, secret);
            if (result.status = 200) {
                data.secret = result.secret;
                data.identity = data.user;

                console.log(data);

                // Invoke the chaincode function
                let response = await FabricController.invoke(data.user, config.channel, config.chaincode, 'createUser', data);
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

    async getUser(data) {
        const str = data.secret;
        const split_str = str.split('-#');
        console.log(split_str);
        console.log(split_str[1]);
        // Invoke the chaincode function
        let response = await FabricController.invoke(split_str[1], config.channel, config.chaincode, 'getUser', data);
        return response;
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

    async sendCoins(data) {
        // Invoke the chaincode function
        let response = await FabricController.invoke(data.user, config.channel, config.chaincode, 'sendBalance', data);
        return response;
    }

    async checkAddressLabel(data) {
        // Invoke the chaincode function
        let response = await FabricController.invoke(data.user, config.channel, config.chaincode, 'getLabel', data);
        return response;
    }

    makeid(length) {
        var result = '';
        var characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
        var charactersLength = characters.length;
        for (var i = 0; i < length; i++) {
            result += characters.charAt(Math.floor(Math.random() * charactersLength));
        }
        return result;
    }
}

module.exports = UserController;