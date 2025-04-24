import fs from 'fs';
import path from 'path';

export default {
  server: {
    // https: {
    //   key: fs.readFileSync(path.resolve(__dirname, 'cert/key.pem')),
    //   cert: fs.readFileSync(path.resolve(__dirname, 'cert/cert.pem')),
    // },
    
    allowedHosts: ['.ngrok-free.app'],
    port: 5173,
    host: 'localhost',
  },
};