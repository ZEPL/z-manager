export * from "./style.scss"

import _          from "lodash"
import Rx         from "rx"
import React      from "react"
import Immutable  from "immutable"
import {Stateful} from "reusable/Stateful"
import {Allocate} from "./components/Allocate"
import {Progress} from "./components/Progress"
import {request}  from "util/request"
import {fetch}    from "channels/fetch"
import endpoints  from "util/api"
import {ButtonGroup} from "./components/ButtonGroup.jsx"
import cookies from "js-cookie"
import {
  Paper
} from "material-ui"

const getDefaultState = props => {
  const {db} = props
  const containers = (db.get("containers") || Immutable.List())
    .filter(x => x.get("username") == db.getIn(["user", "login"]))

// set port cookie after login in case it doesn't exist
  const port = containers.getIn([0, "port"])

  if (port && !cookies.get(port)) {
    cookies.set("port", port)
  }



  const cluster = db.get("cluster") || Immutable.Map()

  const {
    workers,
    cores,
    coresused,
    memory,
    memoryused
  } = cluster.toJS()

  const isDisabled = isNaN(cores) || isNaN(memory)

  return Immutable.fromJS({
    isDisabled,
    isCreated: !!containers.size,
    isProgress: false,
    containerId: containers.getIn([0, "containerId"]),
    workers,
    cores: {
      value: 1,
      minValue: 1,
      used: parseInt(coresused),
      max: parseInt(cores),
      step: 1
    },
    memory: {
      value: 0.1,
      minValue: 0.1,
      used: Math.round(memoryused / 1024 * 100) / 100,
      max: Math.round(memory / 1024 * 100) / 100,
      step: 0.1
    }
  })
}

const state = new Rx.Subject()

const makeReducer = f => newState => state.onNext(f(newState))

export const mergeState = makeReducer(
  newState => currentState => currentState.merge(newState))

export const mapState = makeReducer(
  ({oldData, newData}) => currentState => currentState.map(
    x => x === oldData? newData : x))

export const actions = {
  updateValue: ({cursor, value}) => {
    const {max, minValue} = cursor.toJS()

    const newValue = Math.min(Math.max(minValue, value), max)
    const newState = cursor.set("value", newValue)

    return {oldData: cursor, newData: newState}
  }
}


export class Landing extends Stateful {
  constructor(props) {
    super({
      props,
      channel: state,
      defaultState: getDefaultState(props)
    })
  }

  shouldComponentUpdate() {return true}

  componentDidMount() {
    fetch.onNext(
      // Rx.Observable.just({
      //   workers: [{}],
      //   cores: 72,
      //   coresused: 15,
      //   memory: 400,
      //   memoryused: 200
      // })
      request({endpoint: endpoints.getCluster})
      .flatMap(cluster =>
        // Rx.Observable.just({containers: {}})
          request({
            endpoint: endpoints.getContainers,
            payload: {username: this.props.db.getIn(["user", "login"])}
          })
          .map(({containers}) => {
            return {
              cluster,
              containers
            }
          })
      )
    )
  }

  componentWillReceiveProps(props) {
    mergeState(getDefaultState(props))
  }

  render() {
    const containerStyle = {
       margin: "0 auto",
       maxWidth: "960px",
       padding: "40px"
    }

    const data = this.state.data
    const isCreated = data.get("isCreated")
    const cores = data.get("cores")
    const memory = data.get("memory")
    const workers = data.get("workers")
    const isDisabled = data.get("isDisabled")
    const isProgress = data.get("isProgress")

    return (
      <div>
        <Paper style={containerStyle}>
          <Allocate title="Cores"
            action={_.compose(mapState, actions.updateValue)}
            disabled={isDisabled || isCreated}
            data={cores} />
          <Allocate title="Memory (GB)"
            action={_.compose(mapState, actions.updateValue)}
            disabled={isDisabled || isCreated}
            data={memory} />
          {
            isProgress? <Progress style={{textAlign: "right"}} />
            :
            <ButtonGroup {...this.props}
              data={data}
              disabled={isDisabled}
              getDefaultState={getDefaultState}
              update={mergeState}
              fetch={fetch.onNext.bind(fetch)}
            />
          }

        </Paper>
      </div>
    )
  }
}

