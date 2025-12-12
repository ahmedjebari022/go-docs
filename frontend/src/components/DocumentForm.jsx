import React, { useState } from 'react'
import documentService from '../services/documentService'

function DocumentForm() {
    const [name, setName] = useState
    function handleChange(e){
        setName(e.target.value)
    }
    async function handleSubmit(e){
        e.preventDefault()
        try {
            const res = await documentService.create(name)
            if (res.status === 200){
                console.log("created succesfully"+res)
            }
        } catch (error) {
            console.log(error)            
        }
    }
    return (
        <form onSubmit={handleSubmit}>
                <input type="text" value={name} onChange={handleChange} />
                <input type="submit" />
        </form>
    )
}

export default DocumentForm