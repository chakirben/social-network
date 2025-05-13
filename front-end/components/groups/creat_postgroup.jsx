'use client'

import React, { useState, useRef } from 'react';
import Divider from './../divider';

export default function CreatPostInGroup( { gpid } ) {
  const inputRef = useRef(null);
  const [imageSrc, setImageSrc] = useState(null);
 
  const [text, setText] = useState('');

  const handleImageClick = () => {
    inputRef.current.click();
  };

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = () => {
        setImageSrc(reader.result);
      };
      reader.readAsDataURL(file);
    }
  };


  const handleSubmit = async (e) => {
    e.preventDefault();
    const file = inputRef.current.files[0];
    const formData = new FormData();
    formData.append('content', text);
    formData.append('privacy', "inGroup");
    formData.append('groupId', gpid)
    if (file) {
      formData.append('image', file);
    }

    try {
      const res = await fetch('http://localhost:8080/api/CreatePost', {
        method: 'POST',
        body: formData,
        credentials: 'include',
      });
      const result = await res.json();
      console.log('Post submitted:', result);
      setText('');
  
      setImageSrc(null);
      inputRef.current.value = null;
    } catch (err) {
      console.error('Post failed:', err);
    }
  };
  
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
