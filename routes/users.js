'use strict';

var express = require('express');
var router = express.Router();
const userController = require('../controllers/UserController.js');
const user_object = new userController();
const CAClient = require('../controllers/CAClient.js');
const CAClientController = new CAClient();
const validator = require('../validator/Validator');

//enroll admin
router.get('/enroll-admin', function (req, res, next) {
    CAClientController.enrollAdmin(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//create user details
router.post('/create-user', validator.createUser, function (req, res, next) {
    user_object.createUser(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//get user details
router.post('/get-user', validator.getUser, function (req, res, next) {
    user_object.getUser(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//add user's addresses
router.post('/add-address', validator.addAddress, function (req, res, next) {
    user_object.addAddress(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//get all users
router.post('/all-users', validator.allUsers, function (req, res, next) {
    user_object.allUsers(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//issue token
router.post('/issue-token', validator.issueToken, function (req, res, next) {
    user_object.issueToken(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//check symbol
router.post('/check-symbol', validator.checkSymbol, function (req, res, next) {
    user_object.checkSymbol(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//get tokens
router.post('/get-token', validator.getToken, function (req, res, next) {
    user_object.getToken(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//transfer token
router.post('/transfer-token', validator.transferToken, function (req, res, next) {
    user_object.transferToken(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//send coins
router.post('/send-coins', validator.sendCoins, function (req, res, next) {
    user_object.sendCoins(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//check addresslabel
router.post('/check-addresslabel', validator.checkAddressLabel, function (req, res, next) {
    user_object.checkAddressLabel(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

module.exports = router;
