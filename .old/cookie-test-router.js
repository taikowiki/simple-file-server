const { YowzaServerResponse, YowzaServerRouter } = require('@yowza/server');

const cookieTestRouter = new YowzaServerRouter('/cookie-test');

cookieTestRouter.addHandler(async (event) => {
    console.log(Array.from(event.request.cookie.entries()))
    return new YowzaServerResponse({
        type: 'raw',
        content: event.request.cookie.get('auth-user') ?? 'no cookie'
    })
})

module.exports = cookieTestRouter;