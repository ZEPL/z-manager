import __curry from "lodash/function/curry"

export const curry = __curry

export const wrap = curry((cursor, data) => {
  let obj = {}
  obj[cursor] = data
  return obj
})
