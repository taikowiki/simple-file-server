const { YowzaServerError, YowzaServerResponse, YowzaServerRouter } = require('@yowza/server');
const { existsSync, createReadStream } = require('fs');

const imageRouter = new YowzaServerRouter('/img/:fileName');

imageRouter.addHandler(async (event) => {
    const fileName = event.params.fileName;
    const filePath = __dirname + '/files/img/' + fileName;

    if (!existsSync(filePath)){
        return new YowzaServerError(404);
    }

    return new YowzaServerResponse({
        type: 'file',
        content: createReadStream(filePath)
    })
})

module.exports = imageRouter;