import React from 'react';
import Rx from 'rx'

export class Alert extends React.Component {
  render() {
    const {message} = this.props.children

    const style={
      zIndex: '10',
      lineHeight: '24px',
      fontSize: '16px',
      padding: '20px',
      position: 'absolute',
      right: '20px',
      top: this.props.index *  70 + 'px',
      backgroundColor: '#FF4081',
      color: '#fff',
      textTransform: 'uppercase',
      cursor: 'pointer',
      boxShadow: '0px 1px 6px rgba(0, 0, 0, 0.12), 0px 1px 4px rgba(0, 0, 0, 0.24)',
      display: message? 'block' : 'none'
    }

    return (
      <p style={style} {...this.props}>
        {message.toString()}
      </p>
    )
  }
}



