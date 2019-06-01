var express = require('express');
var router = express.Router();
const config = require('../config/config.js');
// const FabricCAServices = require('fabric-ca-client');
const { FileSystemWallet, Gateway, X509WalletMixin } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const userController = require('../controllers/UserController.js');
const user_object = new userController();
const ccpPath = path.resolve(__dirname, '..', 'config', 'connection-org1.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

//enroll admin
router.get('/enroll-admin', function (req, res, next) {
    user_object.enrollAdmin(req.body).then(result => {
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

//enroll user1
router.get('/enroll-user', function (req, res, next) {
    user_object.enrollUser(req.body).then(result => {
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

//create a new user
router.post('/create-user', function (req, res, next) {
    user_object.createUser(req.body).then(result => {
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

//get user details
router.post('/get-user', function (req, res, next) {
    user_object.getUser(req.body).then(result => {
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

//get all users
router.post('/all-users', function (req, res, next) {
    user_object.allUsers(req.body).then(result => {
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

//issue token
router.post('/issue-token', function (req, res, next) {
    user_object.issueToken(req.body).then(result => {
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

//get tokens
router.post('/get-token', function (req, res, next) {
    user_object.getToken(req.body).then(result => {
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

//transfer token
router.post('/transfer-token', function (req, res, next) {
    user_object.transferToken(req.body).then(result => {
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

module.exports = router;
