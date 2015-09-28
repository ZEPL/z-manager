export * from "./style.scss"

import React          from "react"
import {RouteHandler} from "react-router"
import {Alert}        from "reusable/Alert"
import {Navbar}       from "pages/Navbar"

class Container extends React.Component {
  render() {
    return (
      <div className="container">
        {
          this.props.errors.map(
            (x, i) =>
              <Alert index={i}>
                {{message: x}}
              </Alert>
          )
        }

        <Navbar {...this.props} />
        <RouteHandler {...this.props} />
      </div>
    )
  }
}

export default Container
