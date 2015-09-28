import {transition as channel} from "channels/transition"

export default (Handler, config) => {
  Handler.willTransitionTo = (transition, params, query, go) => {
    channel.onNext({
      config: config || {}, 
      transition, 
      params, 
      query, 
      go 
    })
  }

  return Handler
}

