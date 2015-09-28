import React     from "react"
import Immutable from "immutable"
import Material  from "reusable/Material"
import {wrap}    from "util/util"

export class Stateful extends Material {
  constructor(props) {
    super(props.props)

    const defaultState = props.defaultState || {}

    this.channel = props.channel
      .scan(defaultState, (appState, reducer) => reducer(appState))
      .map(wrap('data'))
      .startWith(wrap('data', defaultState))
  }

  componentWillMount() {
    this.channel.subscribe(state => this.setState(state))
  }

  shouldComponentUpdate(props, state) {
    return !Immutable.fromJS(props).equals(Immutable.fromJS(this.props)) || 
           !Immutable.fromJS(state).equals(Immutable.fromJS(this.state))
  }

}
