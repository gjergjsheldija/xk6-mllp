import mllp from 'k6/x/mllp';

const client = new mllp.Client({
    host: '127.0.0.1',
    port: '5000'
});

export default function () {
    client.send('./sample.hl7');
}
