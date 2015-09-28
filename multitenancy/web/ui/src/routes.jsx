import React from "react"
import Immutable from "immutable"
import {
  Route,
  DefaultRoute,
  NotFoundRoute
} from "react-router"

import enhanceHandler from "util/enhanceHandler"
import {userRoles}    from "util/types"

import endpoints from "util/api"
import {Landing} from "pages/Landing"
import Login     from "pages/Login"
import {Admin}   from "pages/Admin"
import NotFound  from "pages/404"
import Container from "pages/Container"

const rules = Immutable.Map()

// const GetLanding = Landing
const GetLanding = enhanceHandler(
  Landing,
  rules
    .set(userRoles.guest, {redirect: "login"})
)

const GetAdmin = enhanceHandler(
  Admin,
  rules
    .set(userRoles.guest, {redirect: "login"})
// FIXME: currently showing to users, but needs to be shown only to admins
    // .set(userRoles.user, {
      // redirect: "/404",
    //   endpoints: [["containers", endpoints.getContainers]]
    // })
)


export default (
  <Route path="/" handler={Container}>
    <DefaultRoute name="landing" handler={GetLanding} />

    <Route name="login"  path="login"  handler={Login} />
    <Route name="admin" path="admin" handler={Admin} />

    <NotFoundRoute handler={NotFound} />
  </Route>
)
