import Rx          from "rx"
import Immutable   from "immutable"
import {fetch}     from "channels/fetch"
import endpoints   from "util/api"
import {wrap}      from "util/util"
import {request}   from "util/request"
import {userRoles} from "util/types"
import cookies     from "js-cookie"

import {loginChannel} from "pages/Login"
import {logoutChannel} from "pages/Navbar"

export const user = Rx.Observable.merge(
  request({endpoint: endpoints.whoiam}),

  loginChannel.flatMap(
    payload => request({payload, endpoint: endpoints.login})
  ).do(x => {cookies.set("username", x.login)}),

  logoutChannel.flatMap(
    payload => request({endpoint: endpoints.logout})
  ).map(x => {return {type: userRoles.guest}})
  .do(x => {
    cookies.remove("username")
    cookies.remove("port")
  })
).share()

// const f = fetch.mergeAll().startWith({containers: []})
// fetch.onNext(
//   request({endpoint: endpoints.getContainers})
//   .map(wrap("containers"))
// )

// export const db = Rx.Observable.combineLatest(
//   user,
//   fetch.mergeAll().startWith({}),
//   // request({
//   //   endpoint: endpoints.getContainers,
//   //   payload: {username: "all"}
//   // }),
//   // Rx.Observable.just({
//   //   cores: 72,
//   //   coresused: 24,
//   //   memory: 829440,
//   //   memoryused: 540672
//   // }),
//   // request({endpoint: endpoints.getCluster, error: "can't reach spark cluster"})
//   //   .map(x => x || wrap("cluster")(x)),
//   (user, fetch) => {
//     let __containers = [] // containers
//     let __cluster = {}
//
//     if (fetch.update) {
//       __containers = fetch.update.containers || __containers
//       __cluster = fetch.update.cluster || __cluster
//     }
// console.log(__cluster)
//     return {
//       user,
//       containers: __containers,
//       cluster: __cluster
//     }
//   }
// )


export const db = Rx.Observable.combineLatest(
  user,
  fetch.mergeAll().startWith({}),
  (user, fetch) => {
    return Immutable.fromJS({user}).merge(fetch)
  }
).scan(Immutable.Map(), (db, x) => {
  return db.merge(x)
})
.share()







