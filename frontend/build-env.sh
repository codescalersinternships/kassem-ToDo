#!/bin/sh


file="src/config.js"
echo $file
configs='
const BACKEND_URL= "http://192.168.59.100:30008/";
export default BACKEND_URL
'
echo $configs > $file
npm run build
npm start 
