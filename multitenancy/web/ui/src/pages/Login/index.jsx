export * from "./style.scss"

import React    from "react"
import Immutable from "immutable"
import mui      from "material-ui"
import {Stateful} from "reusable/Stateful"
import {curry}       from 'util/util'
import {userRoles}   from 'util/types'
import {Input}       from './components/Input'

const ThemeManager = new mui.Styles.ThemeManager()
const {RaisedButton, Paper} = mui

const stateChannel = new Rx.Subject()
export const updateState = stateChannel.onNext.bind(stateChannel)

export const loginChannel = new Rx.Subject()
const attemptLogin = loginChannel.onNext.bind(loginChannel)

export const actions = {
  input: curry(({cursor, value}, oldState) => {
    const isValid = value => value.length > 3
    const getValidity = (cursor, state) => state.getIn([cursor, 'isValid'])

    return oldState.withMutations(state =>
      state
        .set(
          cursor,
          Immutable.fromJS({value, isValid: isValid(value)}))
        .set(
          'loginButtonIsActive',
          getValidity('username', state) && getValidity('password', state))
    )
  }),
  login: (event, state) => {
    event.preventDefault()
    const username = state.getIn(['username', 'value'])
    const password = state.getIn(['password', 'value'])

    return {login: username, password}
  },
}


class Login extends Stateful {
  constructor(props) {
    super({
      props,
      channel: stateChannel,
      defaultState: Immutable.fromJS({
        username: {
          value: '',
          isValid: false
        },
        password: {
          value: '',
          isValid: false
        },
        loginButtonIsActive: false
      })
    })

  }

  render() {
    const userType = this.props.db.getIn(['user', 'type'])
    const loginButtonIsActive = this.state.data.get('loginButtonIsActive')

    const {getCurrentQuery, transitionTo} = this.context.router

    if (userType && userType != userRoles.guest) {
      transitionTo(getCurrentQuery().nextPath || '/')
    }

    return (
      <div className="login-page">
        <h1 className="login-header">Z-Manager</h1>
        <Paper className="login-wrapper" zDepth={1}>
          <form className="login-form"
            onSubmit={e => attemptLogin(actions.login(e, this.state.data))}>
            <Input className="login-field"
              cursor="username"
              hintText="username"
              type="text" />
            <Input className="login-field"
              cursor="password"
              hintText="password"
              type="password" />
            <RaisedButton style={{margin: "2em auto"}}
              label="Log In"
              primary={true}
              disabled={!loginButtonIsActive}
              type="submit" />
          </form>
        </Paper>
      </div>
    )
  }
}

Login.contextTypes = {router: React.PropTypes.func}
Login.childContextTypes = {muiTheme: React.PropTypes.object}

export default Login
