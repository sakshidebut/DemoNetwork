'use strict';

// Validate User data
function createUser(req, res, next) {
    // Check name
    req.checkBody('name')
        .exists().withMessage('The name field is required.')
        .notEmpty().withMessage('The name field is required.')
        .isString().withMessage('The name must be a string.');

    // Check email
    req.checkBody('email')
        .exists().withMessage('The email field is required.')
        .notEmpty().withMessage('The email field is required.')
        .isEmail().withMessage('Please enter a valid email address.')
        .isString().withMessage('The email must be a string.');

    // Check phone
    req.checkBody('phone')
        .exists().withMessage('The phone field is required.')
        .notEmpty().withMessage('The phone field is required.')
        .isNumeric().withMessage('The phone may only contain digits.')
        .isLength({ min: 10 }).withMessage('The phone must be 10 digits.');

    // Check address
    req.checkBody('address')
        .exists().withMessage('The address field is required.')
        .notEmpty().withMessage('The address field is required.')
        .isLength({ min: 5, max: 500 }).withMessage('The address must be between 5 and 500 characters.');


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
    // Check user id
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

module.exports = { createUser, getUser, allUsers, issueToken, checkSymbol, getToken, transferToken };