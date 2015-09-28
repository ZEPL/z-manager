import React from 'react'
import {
  CircularProgress
} from "material-ui"

const buttonStyle = {
  margin: "0 1em",
  verticalAlign: "top"
}

export class Progress extends React.Component {
  render() {
    return (
      <div {...this.props}>
        <p style={{display: "inline-block", marginTop: "0.5em"}}>
          Working. Please wait...
        </p>
        <CircularProgress style={buttonStyle} mode="indeterminate" size={0.5} />
      </div>
    )
  }
}
