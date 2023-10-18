import {Hl7} from 'k6/x/mllp';

const config = {
    host: '127.0.0.1',
    port: '2575'
}

const client = new Hl7(config);

export const options = {
    vus: 100,
    duration: '30s',
};

 export default function () {
    client.send('./examples/sample.hl7');
 }
