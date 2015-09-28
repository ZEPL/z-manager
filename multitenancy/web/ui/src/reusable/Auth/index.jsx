import React from "react"

export class Auth extends React.Component {
  render() {
    const isAllowed = this.props.roles.indexOf(this.props.type) > -1

    return isAllowed?
      <div {...this.props}>{this.props.children}</div>
      : null
  }
}
