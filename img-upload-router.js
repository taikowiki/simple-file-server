const { defineDBHandler } = require('@yowza/db-handler');
const { YowzaServerResponse, YowzaServerRouter, YowzaServerError } = require('@yowza/server');
const { createDecipheriv } = require('crypto');
const dotenv = require('dotenv');
const mime = require('mime-types');
const fs = require('fs');

dotenv.config();

function decipher(encrypted, key) {
    const bufferKey = Buffer.from(key, 'hex');
    const decrypt = createDecipheriv('aes-256-gcm', bufferKey, bufferKey);
    const decrypted = decrypt.update(encrypted, 'hex', 'utf-8');

    return decrypted;
}

const getUserData = defineDBHandler((provider, providerId) => {
    return async(run) => {
        const result = await run("SELECT * FROM `user/data` WHERE `provider` = ? AND `providerId` = ?", [provider, providerId]);

        if(result.length === 0){
            return null;
        }
        else{
            return result[0];
        }
    }
})
const logFile = defineDBHandler((UUID, originalFileName, fileName) => {
    return async(run) => {
        return await run("INSERT INTO `file/log` (`UUID`, `originalFileName`, `fileName`) VALUES ()", [UUID, originalFileName, fileName])
    }
})

const imgUploadRouter = new YowzaServerRouter('/upload/img');

imgUploadRouter.addHandler(async (event) => {
    const userToken = event.request.cookie.get('auth-user');
    if(!userToken){
        return new YowzaServerError(403);
    }

    const userProviderData = JSON.parse(decipher(userToken, process.env.AUTH_KEY));
    if(!userProviderData){
        return new YowzaServerError(403);
    }

    const userData = await getUserData(userProviderData.provider, userProviderData.providerId);
    if(userData.grade < 9){
        return new YowzaServerError(401);
    }

    event.response.header.set('Access-Control-Allow-Origin', 'https://taiko.wiki');

    const formData = await event.request.form();
    const file = formData.get('file');

    let fileName;
    while(true){
        fileName = `${Date.now().toString(16)}${Math.random().toString(16)?.replaceAll('.', '')}.${file.mime ? mime.extension(file.mime) : ''}`;
        if(!fs.existsSync(__dirname + '/files/img/' + fileName)){
            break;
        }
    }

    fs.writeFileSync(__dirname + '/files/img/' + fileName, file.content);
    await logFile(userData.UUID, file.filename ?? 'undefined', fileName);

    return new YowzaServerResponse({
        type: 'json',
        content: JSON.stringify({fileName})
    })
})

module.exports = imgUploadRouter;