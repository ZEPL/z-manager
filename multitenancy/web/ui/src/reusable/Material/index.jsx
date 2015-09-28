import React from "react"
import mui   from "material-ui"

const ThemeManager = new mui.Styles.ThemeManager()

export default class Material extends React.Component {
  constructor(props) {
    super(props)
  }

  getChildContext() {
    return {muiTheme: ThemeManager.getCurrentTheme()}
  }
}

Material.childContextTypes = {muiTheme: React.PropTypes.object}


