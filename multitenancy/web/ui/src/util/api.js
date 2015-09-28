const {GET, POST, PUT, DELETE} = {
  GET: "get",
  POST: "post",
  PUT: "put",
  DELETE: "delete"
}

export default {
  whoiam: [GET, () => "users/whoiam"],
  login: [PUT, () => "users/login"],
  logout: [POST, () => "users/logout"],
  getCluster: [GET, () => "cluster"],
  getContainers: [POST, () => "containers/list"],
  createInstance: [POST, () => "containers/create"],
  deleteInstance: [POST, () => "containers/delete"]
}

