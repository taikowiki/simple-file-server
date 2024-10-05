const YowzaServer = require('@yowza/server');
const fumenRouter = require('./fumen-router.js');
const imageRouter = require('./image-router.js');
const cookieRouter = require('./cookie-test-router.js')
const imgUploadRouter = require('./img-upload-router.js');
const linkUploadRouter = require('./link-upload-router.js');

const server = new YowzaServer.default();

server.addMiddleware(async(event) => {
    return event;
})

server.addRouter(imageRouter);
server.addRouter(fumenRouter);
server.addRouter(cookieRouter)
server.addRouter(imgUploadRouter);
server.addRouter(linkUploadRouter);

server.listen({
    http: {
        port: 3000
    }
}, () => {console.log('listen on port 80.')})