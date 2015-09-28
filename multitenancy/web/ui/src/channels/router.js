import _         from 'lodash'
import Rx        from 'rx'
import Router    from 'react-router'
import routes    from 'routes'
import {wrap}    from 'util/util'
import {request} from 'util/request'
import {fetch}   from 'channels/fetch'
import {user}    from 'channels/db'
import {
  transition as transitionCh
} from 'channels/transition'

const __router = Router.create({routes})

export const router = Rx.Observable.fromEventPattern(
  h => __router.run(h),
  _.identity,
  args => {
    const Handler     = args[0]
    const routerState = args[1]
    
    return {Handler, routerState}
  }
)

Rx.Observable.combineLatest(
  user,
  transitionCh,
  ({login, type}, {transition, params, go, config}) => {
    const {redirect, endpoints} = config.get(type) || config.get('_') || {}
    params.user || (params.user = login)

    return {
      go,
      params, 
      redirect,
      endpoints,
      transition
    }
  }
).subscribe(({transition, redirect, params, go, endpoints}) => {
/* Heuristics:
 * needs to send a request? => do send
 *      success? => show target view
 *      error? => redirect to 404
 * or needs redirect? => do redirect
 * or _ => show target view
 */

  if (endpoints && endpoints.length > 0) {
    endpoints.map(([name, endpoint]) =>
      fetch.onNext(request({endpoint, params}).map(wrap(name)))
    )
    go();

//   if (endpoints && endpoints.length > 0) {
//     Rx.Observable.zipArray.apply(this,
//       endpoints.map(([_, endpoint]) =>
//         request({endpoint, params})
//       )
//     ).subscribe(
//       xs => {
//         const state = xs.reduce((state, body, i) => {
//           state[endpoints[i][0]] = body; 
//           return state;
//         }, {})
//         go();
//         console.log(state)
//         // resetState(state)
//       },
// // handle errors
//       ({status}) => {
//         if (status == 400 || status == 404) {
//           router.transitionTo('/404');
//         } else {
//           router.transitionTo(redirect, params);
//         }
//       }
//    )
  } else if (redirect) {
    const query = redirect == 'login'? {nextPath: transition.path} : {}
    __router.transitionTo(redirect, params, query)
  } else {
    go()
  }
})

