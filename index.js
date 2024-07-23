const YowzaServer = require('@yowza/server');
const fumenRouter = require('./fumen-router.js');
const imageRouter = require('./image-router.js')

const server = new YowzaServer.default();

server.addRouter(imageRouter);
server.addRouter(fumenRouter);

server.listen({
    http: {
        port: 80
    }
}, () => {console.log('listen on port 80.')})