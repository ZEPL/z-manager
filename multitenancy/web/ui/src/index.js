export * from 'common.scss'

import Rx       from 'rx'
import React    from 'react'
import {db}     from 'channels/db'
import {router} from 'channels/router'
import {errors} from 'channels/error'

Rx.Observable.combineLatest(
  db,
  errors,
  router,
  Array
).subscribe(([db, errors, {Handler, routerState}]) => {
  React.render(
    <Handler db={db} errors={errors} routerState={routerState} />,
    document.querySelector('.js-app')
  )
}, err => console.log(err))

