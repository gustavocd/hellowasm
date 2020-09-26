import React from 'react'
import AceEditor from 'react-ace'

import 'ace-builds/src-noconflict/mode-golang'
import 'ace-builds/src-noconflict/theme-monokai'
import 'ace-builds/src-noconflict/ext-language_tools'

function LiveEditor(props) {
  const { onChange, value } = props

  return (
    <AceEditor
      mode="golang"
      theme="monokai"
      width="100%"
      height="70vh"
      fontSize={16}
      value={value}
      onChange={onChange}
      name="live_editor"
      editorProps={{ $blockScrolling: true }}
      setOptions={{
        enableBasicAutocompletion: true,
        enableLiveAutocompletion: true,
        enableSnippets: true
      }}
    />
  )
}

export default LiveEditor
