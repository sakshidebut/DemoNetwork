var express = require('express');
var router = express.Router();
const config = require('../config/config.js');
const userController = require('../controllers/UserController.js');
const user_object = new userController();

//enroll admin
router.get('/enroll-admin', function (req, res, next) {
    user_object.enrollAdmin(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//enroll user1
router.get('/enroll-user', function (req, res, next) {
    user_object.enrollUser(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//create a new user
router.post('/create-user', function (req, res, next) {
    user_object.createUser(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//get user details
router.post('/get-user', function (req, res, next) {
    user_object.getUser(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//get all users
router.post('/all-users', function (req, res, next) {
    user_object.allUsers(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        res.status(result.status).json(result.data);
    });
});

//issue token
router.post('/issue-token', function (req, res, next) {
    user_object.issueToken(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//get tokens
router.post('/get-token', function (req, res, next) {
    user_object.getToken(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

//transfer token
router.post('/transfer-token', function (req, res, next) {
    user_object.transferToken(req.body).then(result => {
        res.status(result.status).json(result.data);
    }).catch(result => {
        //error handling
        res.status(result.status).json(result.data);
    });
});

module.exports = router;
