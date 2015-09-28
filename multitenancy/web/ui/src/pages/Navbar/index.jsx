import React       from "react"
import Material    from "reusable/Material"
import {Auth}        from "reusable/Auth"
import {Button}    from "reusable/Button"
import {userRoles} from "util/types"

export const logoutChannel = new Rx.Subject()
export class Navbar extends Material {
  render() {
    const user = this.props.db.get("user")
    const username = user.get("login")
    const buttonStyle = {
      margin: "0 5px",
      verticalAlign: "top"
    }

    const isAdmin = this.props.routerState.path.indexOf("admin") > -1

    return (
      <div style={{paddingTop: "100px"}}>
        <div style={{position: "absolute", top: "20px", right: "20px"}}>
          <Auth roles={[userRoles.user, userRoles.admin]} {...user.toJS()}>
            <label style={{marginRight: "20px"}}>
              {username}
            </label>
            <Auth roles={[userRoles.admin]} {...user.toJS()}
              style={{display: "inline-block"}}>
              <Button
                linkButton={true}
                href={isAdmin? "#/" : "#/admin"}
                label={isAdmin? "Home page" : "Admin page"}
                style={buttonStyle} />
            </Auth>
            <Button
              style={buttonStyle}
              primary={true}
              label="Logout"
              onClick={_ => logoutChannel.onNext(true)} />
          </Auth>
        </div>
      </div>
    )
  }
}
