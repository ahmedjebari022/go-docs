import { forwardRef, useEffect, useImperativeHandle, useMemo, useRef, useState } from 'react'
import { createEditor, Editor, Transforms } from 'slate'
import { Editable, Slate, withReact } from 'slate-react'
import { withCursors, withYjs, YjsEditor } from '@slate-yjs/core'
import { WebsocketProvider } from 'y-websocket'
import { RemoteCursorOverlay } from './CursorOverlay'
import * as Y from 'yjs'

// Default fallback
const defaultInitialValue = [{ type: 'paragraph', children: [{ text: '' }] }]

export const TextEditor = forwardRef(({userEmail, documentId, initialContent}, ref) => {
  const [connected, setConnected] = useState(false)
  const [sharedType, setSharedType] = useState()
  const [provider, setProvider] = useState()
  const slateEditorRef = useRef({})

  const randomColor = useMemo(() => '#' + Math.floor(Math.random()*16777215).toString(16).padStart(6, '0'), [])
  
  useEffect(() => {
    const yDoc = new Y.Doc()
    const sharedDoc = yDoc.get('slate', Y.XmlText)
    const yProvider = new WebsocketProvider('ws://localhost:1234', documentId, yDoc)

    yProvider.on('sync', (isSynced) => {
      if (isSynced) {
        // CRITICAL FIX: Hydrate Yjs directly here, before React even sees it
        // We check if the Yjs document is empty.
        if (sharedDoc.toDelta().length === 0) {
            // If we have DB content, we insert it into Yjs
            if (initialContent && initialContent.length > 0) {
                console.log("Hydrating Yjs from DB content")
                // We can't use Slate transforms here easily without the editor instance.
                // But we can rely on the SlateEditor component to do it safely now that we know we are synced.
            }
        }
        setConnected(true)
      }
    })

    yProvider.awareness.setLocalStateField('color', randomColor)
    yProvider.awareness.setLocalStateField('name', userEmail)

    setSharedType(sharedDoc)
    setProvider(yProvider)

    return () => {
      yDoc?.destroy()
      yProvider?.destroy()
    }
  }, [documentId, userEmail, randomColor]) // Added dependencies

  useImperativeHandle(ref, ()=>({
    save: () => {
      if (slateEditorRef.current){
        return slateEditorRef.current.getSaveData()
      } 
      return null
    }
  }))

  if (!connected || !sharedType || !provider) {
    return <div>Loading Editor...</div>
  }

  return (
    <SlateEditor 
      ref={slateEditorRef} 
      sharedType={sharedType} 
      initialContent={initialContent} 
      provider={provider} 
    />
  )
})

const SlateEditor = forwardRef(({ sharedType, provider, initialContent }, ref) => {
  const editor = useMemo(() => {
    const e = withReact(withCursors(withYjs(createEditor(), sharedType), provider.awareness))
    const { normalizeNode } = e
    e.normalizeNode = (entry, options) => {
      const [node] = entry
      if (!Editor.isEditor(node) || node.children.length > 0) {
        return normalizeNode(entry, options)
      }
      Transforms.insertNodes(editor, defaultInitialValue, { at: [0] })
    }
    return e
  }, [sharedType, provider])

  // Connect Yjs to Slate
  useEffect(() => {
    YjsEditor.connect(editor)
    return () => YjsEditor.disconnect(editor)
  }, [editor])

  // HYDRATION LOGIC (The Fix)
  // We run this ONCE when the component mounts.
  useEffect(() => {
    // 1. Check if Yjs is empty (length 0)
    if (sharedType.length === 0 && initialContent && initialContent.length > 0) {
        console.log("Injecting DB content into empty Yjs document")
        
        // 2. Use Slate transforms to insert the data. 
        // Because Yjs is connected, this will update Yjs automatically.
        Editor.withoutNormalizing(editor, () => {
            // Clear any default empty paragraph Slate might have created
            editor.children.forEach((node) => {
                Transforms.removeNodes(editor, { at: [0] })
            })
            // Insert the DB content
            Transforms.insertNodes(editor, initialContent, { at: [0] })
        })
    }
  }, []) // Run once on mount. Yjs is already connected because parent waited for 'connected' state.

  useImperativeHandle(ref, () => ({
    getSaveData: () => editor.children
  }))

  return (
    <Slate editor={editor} initialValue={defaultInitialValue}>
      <RemoteCursorOverlay awareness={provider.awareness}>
        <Editable className="min-h-[500px]" />
      </RemoteCursorOverlay>
    </Slate>
  )
})


