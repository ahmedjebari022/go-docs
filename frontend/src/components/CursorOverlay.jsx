import { useRemoteCursorOverlayPositions } from '@slate-yjs/react'
import React, { useRef, useEffect, useState } from 'react'

export const RemoteCursorOverlay = ({ children, className, style, awareness }) => {
  const containerRef = useRef(null)
  const [cursors] = useRemoteCursorOverlayPositions({
    containerRef,
  })

  const [, setTick] = useState(0)
  useEffect(() => {
    if (!awareness) return
    const handler = () => setTick(t => t + 1)
    awareness.on('change', handler)
    return () => awareness.off('change', handler)
  }, [awareness])

  return (
    <div className={className} style={{ position: 'relative', ...style }} ref={containerRef}>
      {children}
      {cursors.map((cursor) => {
        const clientState = awareness?.getStates().get(cursor.clientId)
        
        const name = clientState?.name || 'Guest'
        const color = clientState?.color || '#3b82f6'

        return (
          <React.Fragment key={cursor.clientId}>
            {/* Render the caret */}
            {cursor.caretPosition && (
              <div
                style={{
                  ...cursor.caretPosition,
                  position: 'absolute',
                  backgroundColor: color,
                  width: '2px',
                  zIndex: 10,
                }}
              >
                <div
                  style={{
                    position: 'absolute',
                    top: -20,
                    left: -2,
                    backgroundColor: color,
                    color: 'white',
                    fontSize: '12px',
                    padding: '2px 4px',
                    borderRadius: '4px',
                    whiteSpace: 'nowrap',
                    pointerEvents: 'none',
                  }}
                >
                  {name}
                </div>
              </div>
            )}
            
            {/* Render the selection rectangles */}
            {cursor.selectionRects.map((rect, i) => (
              <div
                key={i}
                style={{
                  ...rect,
                  position: 'absolute',
                  backgroundColor: color,
                  opacity: 0.2,
                  pointerEvents: 'none',
                }}
              />
            ))}
          </React.Fragment>
        )
      })}
    </div>
  )
}