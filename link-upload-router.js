const { defineDBHandler } = require('@yowza/db-handler');
const { YowzaServerResponse, YowzaServerRouter, YowzaServerError } = require('@yowza/server');
const { default: axios } = require('axios');
const { createDecipheriv, randomUUID } = require('crypto');
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
    return async (run) => {
        const result = await run("SELECT * FROM `user/data` WHERE `provider` = ? AND `providerId` = ?", [provider, providerId]);

        if (result.length === 0) {
            return null;
        }
        else {
            return result[0];
        }
    }
})

const logFile = defineDBHandler((UUID, originalFileName, fileName) => {
    return async (run) => {
        return await run("INSERT INTO `file/log` (`UUID`, `originalFileName`, `fileName`) VALUES (?, ?, ?)", [UUID, originalFileName, fileName])
    }
})

const linkUploadRouter = new YowzaServerRouter('/upload/link');

linkUploadRouter.addHandler(async (event) => {
    const requestData = await event.request.json();

    const userToken = event.request.cookie.get('auth-user');

    if (!requestData.key || userToken) {
        return new YowzaServerError(403);
    }

    if (requestData.key && requestData.key !== process.env.API_KEY) {
        return new YowzaServerError(403);
    }

    if (userToken) {
        const userProviderData = JSON.parse(decipher(userToken, process.env.AUTH_KEY));
        if (!userProviderData) {
            return new YowzaServerError(403);
        }

        var userData = await getUserData(userProviderData.provider, userProviderData.providerId);
        if (userData.grade < 9) {
            return new YowzaServerError(401);
        }
    }

    if (typeof (requestData.url) !== "string") {
        return new YowzaServerError(400);
    }

    event.response.header.set('Access-Control-Allow-Origin', 'https://taiko.wiki');
    event.response.header.set('Access-Control-Allow-Credentials', 'true');

    try {
        var response = await axios({
            method: 'get',
            url: requestData.url,
            responseType: 'stream'
        })
    }
    catch {
        return new YowzaServerError(500);
    }

    let mimeType = mime.extension(response.headers['Content-Type']) || '';
    while (true) {
        var fileName = `${randomUUID()}${mimeType ? '.' + mimeType : ''}`;
        if (!fs.existsSync(__dirname + '/files/img/' + fileName)) {
            break;
        }
    }

    const writeStream = fs.createWriteStream(__dirname + '/files/img/' + fileName);
    try {
        await new Promise((res, rej) => {
            response.data.on('end', res);
            response.data.on('error', rej);
            response.data.pipe(writeStream);
        });
    }
    catch {
        return new YowzaServerError(500);
    }

    await logFile(userData?.UUID ?? 'server', requestData.url, fileName);

    return new YowzaServerResponse({
        type: 'json',
        content: JSON.stringify({ fileName })
    });
})

module.exports = linkUploadRouter;