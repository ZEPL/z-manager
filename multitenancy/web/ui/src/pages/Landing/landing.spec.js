import Immutable from "immutable"
import {actions} from "./index.jsx"

const defaultState = Immutable.fromJS({
  cores: {
    value: 5,
    allocated: 3,
    max: 8,
    step: 1
  },
  memory: {
    value: 600,
    allocated: 512,
    max: 1024,
    step: 10
  }
})

describe("landing page", () => {

  it("resets state", () => {
    const newState = defaultState.setIn(["cores", "value"], 0)

    expect(
      actions.resetState(newState, defaultState).equals(newState)
    ).toBe(true)
  })

  it("updates value", () => {
    const payload = {
      cursor: defaultState.get("cores"),
      value: 0
    }
    const newState = defaultState.setIn(["cores", "value"], 0)

    expect(
      actions.updateValue(payload, defaultState).equals(newState)
    ).toBe(true)
  })

})



