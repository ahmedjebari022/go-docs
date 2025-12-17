import { LucideOctagonPause } from 'lucide-react'
import React, { useCallback, useMemo, useState, useEffect, useRef} from 'react'
import { createEditor} from 'slate'
import {Slate, Editable, withReact, DefaultElement} from 'slate-react'

const initialValue = [
    {
        type: 'paragraphe',
        children: [{text: 'A line of text in  a pragrah.'}],
    }
]

function TextEditor({documentId,handleUpdate}) {
    const isRemote = useRef(false)
    const socketRef = useRef(null)
    useEffect(() => { 
        const socket = new WebSocket(`//localhost:8081/ws/${documentId}`)
        socket.onmessage = (event)=>{
            console.log("received",event.data)
            const ops = JSON.parse(event.data)
            isRemote.current = true 
            ops.forEach(op => {
               editor.apply(op) 
            })
        }

        socket.onopen = () => {
            console.log("Connected")
        }
        socket.onerror = (error) => {
            console.log(error)
        }
        socket.onclose = (event) => {
            console.log(`closed by ${event}`)
        }
        socketRef.current = socket
        
        return () => {
            socket.close()
        }
    },[])
    
    const [editor] = useState(() => withReact(createEditor()))

    function logOerations(value){
        if(isRemote.current){
            isRemote.current = false
            return
        }
        handleUpdate(value)
        const ops = editor.operations.filter(
            op => op.type !== 'set_selection'
        )
        if (ops.length > 0){
            socketRef.current.send(JSON.stringify(ops))
        }
    }
  return (
    <Slate 
    editor={editor} 
    initialValue={initialValue}
    onChange={logOerations} 
    >
        <Editable />
    </Slate>
  )
}

export default TextEditor