var newman = require('newman'); // require newman in your project

// call newman.run to pass `options` object and wait for callback
newman.run({
    collection: require('./song-collection-app.postman_collection.json'),
    reporters: 'cli'
}, function (err) {
    if (err) { throw err; }
    console.log('collection run complete!');
}).on('beforeItem', function (err, args) { // on start of run, log to console
    // Teardown alternatives for CI running might be pointing to different DB tables
    // within a test environment. That would avoid collision when running tests concurrently
    // in the same DB. Still, it's tricky, the environment setup for a big team would need
    // a bit more thought.
    console.log('TODO: Teardown database.');
});