const { YowzaServerError, YowzaServerResponse, YowzaServerRouter } = require('@yowza/server');
const { existsSync, createReadStream } = require('fs');

const fumenRouter = new YowzaServerRouter('/fumen/:songNo/:difficulty');

fumenRouter.addHandler(async (event) => {
    const fileName = event.params.songNo + '/' + event.params.difficulty + '.png';
    const filePath = __dirname + '/files/fumen/' + fileName;

    if (!existsSync(filePath)){
        return new YowzaServerError(404);
    }

    return new YowzaServerResponse({
        type: 'file',
        content: createReadStream(filePath)
    })
})

module.exports = fumenRouter;