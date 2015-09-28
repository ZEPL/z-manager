export * from "./style.scss"

import React     from "react"
import Immutable from "immutable"
import Material  from "reusable/Material"
import {request} from "util/request"
import {wrap}    from "util/util"
import endpoints from "util/api"
import {fetch}   from "channels/fetch"
import {
  Table
} from "material-ui"

export class Admin extends Material {
  componentDidMount() {
    request({
      endpoint: endpoints.getContainers,
      payload: {username: "all"}
    })
  }

  componentWillReceiveProps(props) {
  console.log(props.db.toJS())
  }

  render() {
    const containers = this.props.db.get("containers")
    .sortBy(x => +x.get("cores"))
    .toJS()

    return (
      <div>
        <h1 style={{textAlign: "center"}}>Running instances</h1>
        <div className="table-responsive-vertical shadow-z-1">
          <table
            className="table table-bordered table-hover table-mc-deep-orange">
            <thead>
              <tr>
                {
                  [
                    "Username",
                    "Container Id",
                    "Cores",
                    "Memory",
                    "Port",
                    "Delete"
                  ].map(x => <td>{x}</td>)
                }
              </tr>
            </thead>
            <tbody>
              {containers.map(
                ({containerId, cores, memory, port, username}, i) =>
                  <tr
                    onClick={e => {
                      let targetId = void 0
                      try {
                        targetId = this.refs["delete-" + i].getDOMNode().id
                      } catch(e) {}

                      if (e.target.id == targetId) {
                        fetch.onNext(request({
                          endpoint: endpoints.deleteInstance,
                          payload: {
                            containerId,
                            username
                          }
                        }).flatMap(x => {
                          return request({
                            endpoint: endpoints.getContainers,
                            payload: {username: "all"}
                          })}).map(wrap("update"))
                        )
                      }
                    }}>
                    <td>{username}</td>
                    <td>{containerId}</td>
                    <td>{cores}</td>
                    <td>{memory}</td>
                    <td>{port}</td>
                    <td
                      id="admin-table-delete"
                      ref={"delete-" + i}
                      style={{
                        color: "red",
                        textAlign: "center",
                        cursor: "pointer"
                      }}>X</td>
                  </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    )
  }
}

