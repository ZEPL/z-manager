import React from "react"
import {Button} from "reusable/Button"
import {request} from "util/request"
import endpoints from "util/api"
import cookies from "js-cookie"

const actions = {
  toggleProgressIndicator: value => {
    return {isProgress: value}
  }
}

export class ButtonGroup extends React.Component {
  render() {
    const {
      data,
      db,
      update,
      fetch,
      disabled,
      getDefaultState
    } = this.props

    const username = db.getIn(["user", "login"])
    const containerId = data.get("containerId")
    const workers = data.get("workers") || {size: 1}
    const cores = data.getIn(["cores", "value"]).toString()
    const memory = data.getIn(["memory", "value"])
    const isCreated = data.get("isCreated")

    const formatMemory = (mem, size) => parseInt(mem * 1024 / size) + "m"

    const reqDelete = request({
      endpoint: endpoints.deleteInstance,
      payload: {containerId, username}
    })

    const reqContainers = request({
      endpoint: endpoints.getContainers,
      payload: {username}
    })

    const reqCreate = request({
      endpoint: endpoints.createInstance,
      payload: {
        username,
        cores,
        memory: formatMemory(memory, workers.size)
      }
    })

    return (
      isCreated?
        <Created {...{reqDelete, reqContainers, update, fetch}} />
        :
        <NotCreated {...{
          db,
          reqCreate,
          reqContainers,
          getDefaultState,
          fetch,
          disabled,
          update
        }} />
    )
  }
}

class Created extends React.Component {
  render() {
    const {
      fetch,
      reqDelete,
      reqContainers,
      update
    } = this.props

    const handleDeleteInstance = _ => {
      fetch(
        reqDelete.flatMap(_ => reqContainers)
          .do(update(actions.toggleProgressIndicator(true)))
          .do(_ => cookies.remove("port"))
      )
    }

    return (
      <div style={{textAlign: "right"}}>
        <Button label="Delete instance" onClick={handleDeleteInstance} />
        <Button label="Open Zeppelin"
          primary={true}
          linkButton={true}
          href="zeppelin" />
      </div>
    )
  }
}

class NotCreated extends React.Component {
  render() {
    const {
      db,
      fetch,
      reqCreate,
      reqContainers,
      update,
      disabled,
      getDefaultState
    } = this.props

    const handleStateReset = _ => {
      update(getDefaultState(this.props))
    }

    const handleCreateInstance = _ => {
      fetch(
        reqCreate
          .flatMap(({port}) => {
            return reqContainers.do(_ => cookies.set("port", port))
          })
          .do(update(actions.toggleProgressIndicator(true)))
      )
    }
    return (
      <div style={{textAlign: "right"}}>
        <Button label="Reset"
          disabled={disabled}
          onClick={handleStateReset} />
        <Button label="Create"
          primary={true}
          disabled={disabled}
          onClick={handleCreateInstance} />
      </div>
    )
  }
}
