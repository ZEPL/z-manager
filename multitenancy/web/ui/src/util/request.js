import Rx                    from "rx"
import superagent            from "superagent"
import {error as errorCh}    from "channels/error"

let host = location.protocol + '//' + location.host + '/api/v1/'

if (__DEV__) {
  host = "http://localhost:3000/api/v1/"
}

export const request = ({endpoint, params, payload, error}) => {
  const [method, end] = endpoint

  const req = superagent[method](host + end(params))
    .send(payload)
    .withCredentials()

  return fromSuperagent(req)
    .catch(err => {
      errorCh.onNext(error || err)
      return Rx.Observable.just({})
    })
    .share()
}

function fromSuperagent(request) {
  return Rx.Observable.create(observer => {
    request.end((err, res) => {
      if (err) {
        let __err
// FIXME: remove this custom errors after backend passes user-friendly messages
        if (request.url.indexOf("cluster") > -1) {
          res.text = JSON.stringify({
            text: "Cannot connect to Spark: " + res.statusText
          })
        } else if (request.url.indexOf("list") > -1) {
          res.text = JSON.stringify({
            text: "Cannot connect to Docker: " + res.statusText
          })
        }

        try {
          let text = JSON.parse(res.text)
          __err = text.text || text.error
        } catch(e) {
          __err = err
        }

        observer.onError(__err)
      } else {
// FIXME: this is a dirty hack that shouldn't be there; blame backend!
        let result
        if (res.body) {
          result = res.body
        } else {
          try {
            result = JSON.parse(res.text)
          } catch(e) {
            result = {}
          }
        }
        observer.onNext(result);
      }

      observer.onCompleted();
    })
  })
}
