import { useState, useRef } from "react"

export default function CommentInput({ id, setComments }) {

    const [cmnt, setcmnt] = useState('')
    const [image, setImage] = useState(null)
    const fileInputRef = useRef(null)

    const handleImageClick = () => {
        fileInputRef.current.click()
    }

    const handleImageChange = (e) => {
        const file = e.target.files[0]
        if (file) {
            setImage(file)
        }
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        if (!cmnt) return
        const formData = new FormData()
        formData.append('content', cmnt)
        formData.append('postId', id)
        if (image) {
            formData.append('image', image)
        }

        try {
            const res = await fetch('http://localhost:8080/api/SetComment', {
                credentials: 'include',
                method: 'POST',
                body: formData
            })
            const result = await res.json()
            console.log('Response:', result)
            setcmnt('')
            setComments(prevComments => [result, ...prevComments || []])
            setImage(null)
        } catch (err) {
            console.error('Error submitting comment:', err)
        }
    }

    return (
        <form className="commentInput spB df cl" onSubmit={handleSubmit}>
            {image && (
                <div className="imagePrvC">
                    <img
                        className="imagePreview"
                        src={URL.createObjectURL(image)}
                        alt="preview"
                    />
                    <img
                        className="closeImage"
                        src="/images/close.svg"
                        alt="close"
                        onClick={() => setImage(null)}
                        style={{ cursor: 'pointer' }}
                    />
                </div>
            )}
            <div className="strech df spB center">
                <div className="df gp6 center">
                    <img
                        className="icn1"
                        src="/images/image.svg"
                        alt="upload"
                        onClick={handleImageClick}
                        style={{ cursor: 'pointer' }}
                    />
                    <input
                        type="file"
                        ref={fileInputRef}
                        style={{ display: 'none' }}
                        accept="image/*"
                        onChange={handleImageChange}
                    />
                    <input
                        type="text"
                        placeholder="comment on this post"
                        value={cmnt}
                        onChange={(e) => setcmnt(e.target.value)}
                        className="searchInput"
                    />
                </div>
                <button
                    type='submit'
                    onClick={handleSubmit}
                    disabled={!cmnt.trim()}
                    className={!cmnt.trim() ? 'button-disabled' : 'button-active'}
                >comment</button>
            </div>
        </form>
    )
}
