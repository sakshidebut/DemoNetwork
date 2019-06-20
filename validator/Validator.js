'use strict';

// Validate createUser API
function createUser(req, res, next) {
    // Check address
    req.checkBody('address')
        .exists().withMessage('The address field is required.')
        .notEmpty().withMessage('The address field is required.');

    // validation errors
    let error = req.validationErrors();
    if (error) {
        let message = error[0].msg;
        res.status(422).json({ message: message, key: error[0].param });
    } else {
        next();
    }
}

// Validate getUser API
function getUser(req, res, next) {
    // Check address
    req.checkBody('secret')
        .exists().withMessage('The secret field is required.')
        .notEmpty().withMessage('The secret field is required.');

    // validation errors
    let error = req.validationErrors();
    if (error) {
        let message = error[0].msg;
        res.status(422).json({ message: message, key: error[0].param });
    } else {
        next();
    }
}

// Validate allUsers API
function allUsers(req, res, next) {
    // Check user id
    req.checkBody('id')
        .exists().withMessage('The id field is required.')
        .notEmpty().withMessage('The id field is required.');

    // validation errors
    let error = req.validationErrors();
    if (error) {
        let message = error[0].msg;
        res.status(422).json({ message: message, key: error[0].param });
    } else {
        next();
    }
}

// Validate issueToken API
function issueToken(req, res, next) {
    // Check user id
    req.checkBody('user_id')
        .exists().withMessage('The user id field is required.')
        .notEmpty().withMessage('The user id field is required.');

    // Check code
    req.checkBody('code')
        .exists().withMessage('The code field is required.')
        .notEmpty().withMessage('The code field is required.');

    // Check quantity
    req.checkBody('quantity')
        .exists().withMessage('The quantity field is required.')
        .notEmpty().withMessage('The quantity field is required.')
        .isNumeric().withMessage('The quantity may only contain digits.');

    // Check label
    req.checkBody('label')
        .exists().withMessage('The label field is required.')
        .notEmpty().withMessage('The label field is required.');

    // validation errors
    let error = req.validationErrors();
    if (error) {
        let message = error[0].msg;
        res.status(422).json({ message: message, key: error[0].param });
    } else {
        next();
    }
}

// Validate checkSymbol API
function checkSymbol(req, res, next) {

    // Check code
    req.checkBody('code')
        .exists().withMessage('The code field is required.')
        .notEmpty().withMessage('The code field is required.');

    // validation errors
    let error = req.validationErrors();
    if (error) {
        let message = error[0].msg;
        res.status(422).json({ message: message, key: error[0].param });
    } else {
        next();
    }
}

// Validate getToken API
function getToken(req, res, next) {
    // Check user id
    req.checkBody('id')
        .exists().withMessage('The id field is required.')
        .notEmpty().withMessage('The id field is required.');

    // validation errors
    let error = req.validationErrors();
    if (error) {
        let message = error[0].msg;
        res.status(422).json({ message: message, key: error[0].param });
    } else {
        next();
    }
}

// Validate transferToken API
function transferToken(req, res, next) {
    // Check from id
    req.checkBody('from_id')
        .exists().withMessage('The from id field is required.')
        .notEmpty().withMessage('The from id field is required.');

    // Check to id
    req.checkBody('to_id')
        .exists().withMessage('The to id field is required.')
        .notEmpty().withMessage('The to id field is required.');

    // Check code
    req.checkBody('code')
        .exists().withMessage('The code field is required.')
        .notEmpty().withMessage('The code field is required.');

    // Check quantity
    req.checkBody('quantity')
        .exists().withMessage('The quantity field is required.')
        .notEmpty().withMessage('The quantity field is required.')
        .isNumeric().withMessage('The quantity may only contain digits.');

    // validation errors
    let error = req.validationErrors();
    if (error) {
        let message = error[0].msg;
        res.status(422).json({ message: message, key: error[0].param });
    } else {
        next();
    }
}

// Validate addAddress API
function addAddress(req, res, next) {
    // Check user id
    req.checkBody('user_id')
        .exists().withMessage('The user id field is required.')
        .notEmpty().withMessage('The user id field is required.');

    req.checkBody('value')
        .exists().withMessage('The value field is required.')
        .notEmpty().withMessage('The value field is required.');

    // Check label
    req.checkBody('label')
        .exists().withMessage('The label field is required.')
        .notEmpty().withMessage('The label field is required.');

    // validation errors
    let error = req.validationErrors();
    if (error) {
        let message = error[0].msg;
        res.status(422).json({ message: message, key: error[0].param });
    } else {
        next();
    }
}

// Validate sendCoins API
function sendCoins(req, res, next) {
    // Check from id
    req.checkBody('from_id')
        .exists().withMessage('The from id field is required.')
        .notEmpty().withMessage('The from id field is required.');

    // Check to id
    req.checkBody('to_id')
        .exists().withMessage('The to id field is required.')
        .notEmpty().withMessage('The to id field is required.');

    // Check quantity
    req.checkBody('quantity')
        .exists().withMessage('The quantity field is required.')
        .notEmpty().withMessage('The quantity field is required.')
        .isNumeric().withMessage('The quantity may only contain digits.');

    // Check label
    req.checkBody('label')
        .exists().withMessage('The label field is required.')
        .notEmpty().withMessage('The label field is required.');

    // validation errors
    let error = req.validationErrors();
    if (error) {
        let message = error[0].msg;
        res.status(422).json({ message: message, key: error[0].param });
    } else {
        next();
    }
}

// Validate addAddress API
function checkAddressLabel(req, res, next) {
    // Check user id
    req.checkBody('user_id')
        .exists().withMessage('The user id field is required.')
        .notEmpty().withMessage('The user id field is required.');

    // validation errors
    let error = req.validationErrors();
    if (error) {
        let message = error[0].msg;
        res.status(422).json({ message: message, key: error[0].param });
    } else {
        next();
    }
}

module.exports = { createUser, getUser, allUsers, issueToken, checkSymbol, getToken, transferToken, sendCoins, addAddress, checkAddressLabel };