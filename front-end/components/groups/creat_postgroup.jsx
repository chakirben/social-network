'use client'

export default function CreatPostInGroup( { gpid } ) {
  
 


  /*user-icon.png*/
  return (
    <form className="creatPostForm" onSubmit={handleSubmit}>
      <div className="searchBar">
        <img src="/user-icon.png" />
        <input
          className="searchInput"
          placeholder="What’s happening ?"
          value={text}
          onChange={(e) => setText(e.target.value)}
        />
      </div>

      <div className="ImagePreviewBox">
        {imageSrc && <img src={imageSrc} alt="Preview" className="preview-img" />}
      </div>

      <Divider />

      <div className='spB'>
        <div className='group'>
          <img
            src="./images/image.svg"
            className="upload-icon"
            onClick={handleImageClick}
          />
          <input
            type="file"
            accept="image/*"
            ref={inputRef}
            onChange={handleFileChange}
            style={{ display: 'none' }}
          />
        </div>
        <button type='submit'>post</button>
      </div>
    </form>
  );
}
