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
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

//enroll user1
router.get('/enroll-user', function (req, res, next) {
    CAClientController.enrollUser(req.body).then(result => {
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

//create a new user
router.post('/create-user', validator.createUser, function (req, res, next) {
    user_object.createUser(req.body).then(result => {
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

//get user details
router.post('/get-user', validator.getUser, function (req, res, next) {
    user_object.getUser(req.body).then(result => {
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

//get all users
router.post('/all-users', validator.allUsers, function (req, res, next) {
    user_object.allUsers(req.body).then(result => {
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

//issue token
router.post('/issue-token', validator.issueToken, function (req, res, next) {
    user_object.issueToken(req.body).then(result => {
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

//get tokens
router.post('/get-token', validator.getToken, function (req, res, next) {
    user_object.getToken(req.body).then(result => {
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

//transfer token
router.post('/transfer-token', validator.transferToken, function (req, res, next) {
    user_object.transferToken(req.body).then(result => {
        res.status(200).send(result);
    }).catch(err => {
        //error handling
        res.status(err.httpStatus || 500).send({ message: err.message, status: 0 });
    });
});

module.exports = router;
