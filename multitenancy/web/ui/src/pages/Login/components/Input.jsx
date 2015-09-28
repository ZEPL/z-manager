import React from 'react'
import {
  updateState,
  actions
} from '../index'
import mui      from "material-ui"
import Material from "reusable/Material"

const {TextField} = mui

export class Input extends Material {
  render() {
    return (
      <TextField className={this.props.className}
        hintText={this.props.hintText}
        onChange={event => updateState(actions.input({
          cursor: this.props.cursor,
          value: event.target.value
        }))}>
        <input type={this.props.type} />
      </TextField>
    )
  }
}
