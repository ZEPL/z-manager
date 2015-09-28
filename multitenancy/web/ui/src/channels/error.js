import Rx from 'rx'

export const error = new Rx.BehaviorSubject(null)
export const errors = Rx.Observable.merge(
  error,
  error.map(x => null).delay(3000)
).scan([], (acc, err) => err? acc.concat(err) : acc.slice(1))
