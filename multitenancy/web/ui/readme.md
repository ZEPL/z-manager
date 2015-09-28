Z-Manager front-end
===================

## Table of contents
1. Installation
2. Development
3. Production

## 1. Installation
+ Install `node` and `npm`
+ Run `(sudo) npm install`

## 2. Development
Tech stack:
+ [npm](https://www.npmjs.com/) for package management
+ [webpack](http://webpack.github.io/) for source processing and bundling
+ [ES6](https://github.com/lukehoban/es6features) via [babel](https://babeljs.io/) as programming language
+ [react](https://facebook.github.io/react), [react-router](http://rackt.github.io/react-router/) and [rxjs](https://reactive-extensions.github.io/RxJS/) as core js stack
+ [jasmine](http://jasmine.github.io/) for testing

`package.json` has `scripts` that can be run via `npm run script-name`, for
example `npm run dev` will run a script under `scripts: {dev: ...}`.

To develop:
+ `npm run dev`
+ navigate to http://localhost:9090/
+ changes in the code will trigger automatic recompilation

To test:
+ `npm run test` will find all files that end on `.spec.js` and run them

## 3. Production
+ `npm run build` to build frontend webapp that is ready for serving from `./public`
