var express = require('express');
var router = express.Router();
const config = require('../config/config.js');
const FabricCAServices = require('fabric-ca-client');
const { FileSystemWallet, Gateway, X509WalletMixin } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve(__dirname, '..', 'config', 'connection-org1.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

class UserController {

    async enrollAdmin() {
        try {
            // Create a new CA client for interacting with the CA.
            const caInfo = ccp.certificateAuthorities['ca.org1.india.com'];
            const caTLSCACertsPath = path.resolve(__dirname, caInfo.tlsCACerts.path);
            const caTLSCACerts = fs.readFileSync(caTLSCACertsPath);
            const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);
    
            // Create a new file system based wallet for managing identities.
            const walletPath = path.join(process.cwd(), 'wallet');
            const wallet = new FileSystemWallet(walletPath);
            console.log(`Wallet path: ${walletPath}`);
    
            // Check to see if we've already enrolled the admin user.
            const adminExists = await wallet.exists('admin');
            if (adminExists) {
                console.log('An identity for the admin user "admin" already exists in the wallet');
                return;
            }
    
            // Enroll the admin user, and import the new identity into the wallet.
            const enrollment = await ca.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' });
            const identity = X509WalletMixin.createIdentity('Org1MSP', enrollment.certificate, enrollment.key.toBytes());
            await wallet.import('admin', identity);
            return 'Successfully enrolled admin user "admin" and imported it into the wallet';
    
        } catch (error) {
            console.error(`Failed to enroll admin user "admin": ${error}`);
            process.exit(1);
        }
    }

    async enrollUser() {
        try {
            // Create a new file system based wallet for managing identities.
          const walletPath = path.join(process.cwd(), 'wallet');
          const wallet = new FileSystemWallet(walletPath);
          console.log(`Wallet path: ${walletPath}`);
  
          // Check to see if we've already enrolled the user.
          const userExists = await wallet.exists(config.user);
          if (userExists) {
              console.log('An identity for the user "'+config.user+'" already exists in the wallet');
              return;
          }
  
          // Check to see if we've already enrolled the admin user.
          const adminExists = await wallet.exists('admin');
          if (!adminExists) {
              console.log('An identity for the admin user "admin" does not exist in the wallet');
              console.log('Run the enrollAdmin.js application before retrying');
              return;
          }
  
          // Create a new gateway for connecting to our peer node.
          const gateway = new Gateway();
          await gateway.connect(ccpPath, { wallet, identity: 'admin', discovery: { enabled: true } });
  
          // Get the CA client object from the gateway for interacting with the CA.
          const ca = gateway.getClient().getCertificateAuthority();
          const adminIdentity = gateway.getCurrentIdentity();
  
          // Register the user, enroll the user, and import the new identity into the wallet.
          const secret = await ca.register({ affiliation: 'org1.department1', enrollmentID: config.user, role: 'client' }, adminIdentity);
          console.log(config.user);
          const enrollment = await ca.enroll({ enrollmentID: config.user, enrollmentSecret: secret });
          const userIdentity = X509WalletMixin.createIdentity('Org1MSP', enrollment.certificate, enrollment.key.toBytes());
          await wallet.import(config.user, userIdentity);
          return 'Successfully registered and enrolled admin user "'+config.user+'" and imported it into the wallet';
    
        } catch (error) {
          console.error(`Failed to register user "`+config.user+`": ${error}`);
          process.exit(1);
        }
    }
}

module.exports = UserController;