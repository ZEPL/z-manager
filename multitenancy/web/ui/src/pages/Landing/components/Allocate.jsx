import Rx    from "rx"
import React from "react"
import Material from "reusable/Material"
import {actions, mapState} from "pages/Landing"
import {
  Slider,
  TextField,
  Paper
} from "material-ui"

const labelStyle = {
  margin: "2em 0 1em",
  textAlign: "center"
}

const sliderStyle = {
  width: "70%",
  margin: "0 -69px",
  textAlign: "center",
  display: "inline-block",
}

const fieldStyle = {
  width: "60px",
  height: "48px",
  lineHeight: "48px",
  display: "inline-block"
}

const colStyle = {
  display: "block"
}

const rowStyle = {
  width: "240px",
}

const displayStyle = {
  width: "80px",
  display: "inline-block",
  verticalAlign: "top",
  textAlign: "center"
}

export class Allocate extends Material {
  render() {
    const {title, data, action, disabled} = this.props
    const {used, max, value, step, minValue} = data.toJS()

    const availableWidth = ((max - +used) * 100) / max || 100
    const usedWidth = 100 - availableWidth
    return (
      <div>
        <h2 style={labelStyle}>{title}</h2>
        <div className="row">
        <div style={{
          position: "relative",
          width: "70%",
          display: "inline-block"
        }}>
          <div style={{
            display: "inline-block",
            position: "relative",
            height: "122px",
            width: usedWidth + "%"
          }}>
            <div style={{
              position: "absolute",
              top: "50%",
              height: "2px",
              backgroundColor: "#FF4081",
              width: "100%"
            }}></div>
          </div>
          <Slider style={{
            width: availableWidth + "%",
            display: "inline-block",
          }}
            name="slider"
            disabled={disabled}
            max={max - +used}
            value={+value}
            step={step}
            onChange={(_, value) => action({
              cursor: data,
              value: +value
            })} />
          </div>
          <div style={displayStyle}>
            <label>Allocate</label>
            <TextField style={fieldStyle} className="col"
              disabled={this.props.disabled}>
              <input type="number"
                value={value}
                min={minValue}
                max={max}
                onChange={e => action({
                  cursor: data,
                  value: e.target.value
                })} />
            </TextField>
          </div>
          <div style={displayStyle}>
            <label style={colStyle}>Free</label>
            <label style={fieldStyle} >{+max - +used}</label>
          </div>
          <div style={displayStyle}>
            <label style={colStyle}>Total</label>
            <label style={fieldStyle} >{max}</label>
          </div>

        </div>
      </div>
    )
  }
}
