import { useEffect, useMemo, useState } from 'react'
import { createEditor, Editor, Transforms } from 'slate'
import { Editable, Slate, withReact } from 'slate-react'
import { withCursors, withYjs, YjsEditor } from '@slate-yjs/core'
import { WebsocketProvider } from 'y-websocket'
import { RemoteCursorOverlay } from './CursorOverlay'

import * as Y from 'yjs'

const initialValue = [
  {
    children: [{ text: '' }],
  },
]

export const TextEditor = ({userEmail, documentId}) => {
  const [connected, setConnected] = useState(false)
  const [sharedType, setSharedType] = useState()
  const [provider, setProvider] = useState()

  const randomColor = useMemo(() => '#' + Math.floor(Math.random()*16777215).toString(16).padStart(6, '0'), [])

  useEffect(() => {
    const yDoc = new Y.Doc()
    const sharedDoc = yDoc.get('slate', Y.XmlText)

    const yProvider = new WebsocketProvider('ws://localhost:1234', documentId, yDoc)


    yProvider.awareness.on('change', changes => {
      console.log(Array.from(yProvider.awareness.getStates().values()))
    })

    yProvider.on('status', event => {
      console.log(event.status) 
      setConnected(event.status === 'connected')
    })

    setSharedType(sharedDoc)
    setProvider(yProvider)

    return () => {
      yDoc?.destroy()
      yProvider?.destroy()
    }
  }, [documentId])

  useEffect(() => {
    if (provider && userEmail) {
      provider.awareness.setLocalStateField('color', randomColor)
      provider.awareness.setLocalStateField('name', userEmail)
    }
  }, [provider, userEmail, randomColor])

  if (!connected || !sharedType || !provider) {
    return <div>Loading…</div>
  }

  return <SlateEditor sharedType={sharedType} provider={provider} />
}

const SlateEditor = ({ sharedType, provider }) => {
  const editor = useMemo(() => {
    const e = withReact(withCursors(withYjs(createEditor(), sharedType), provider.awareness))

    const { normalizeNode } = e
    e.normalizeNode = (entry, options) => {
      const [node] = entry

      if (!Editor.isEditor(node) || node.children.length > 0) {
        return normalizeNode(entry, options)
      }

      Transforms.insertNodes(editor, initialValue, { at: [0] })
    }

    return e
  }, [sharedType, provider])

  useEffect(() => {
    YjsEditor.connect(editor)
    return () => YjsEditor.disconnect(editor)
  }, [editor])

  return (
    <Slate editor={editor} initialValue={initialValue}>
      <RemoteCursorOverlay awareness={provider.awareness}>
        <Editable />
      </RemoteCursorOverlay>
    </Slate>
  )
}


