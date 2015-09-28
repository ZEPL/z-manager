import React from 'react'
import {
  RaisedButton
} from "material-ui"

const buttonStyle = {
  margin: "3em 0.5em 1em",
  verticalAlign: "top"
}

export class Button extends React.Component {
  render() {
    return (
      <RaisedButton style={buttonStyle} {...this.props} />
    )
  }
}
