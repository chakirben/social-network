'use client'

import React, { useState, useRef } from 'react';
import Divider from './divider';

export default function CreatePost() {
  const inputRef = useRef(null);
  const [imageSrc, setImageSrc] = useState(null);
  const [selectedOption, setSelectedOption] = useState('public');
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

  const handleChange = (e) => {
    setSelectedOption(e.target.value);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const file = inputRef.current.files[0];
    const formData = new FormData();
    formData.append('text', text);
    formData.append('audience', selectedOption);
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
      setSelectedOption('public');
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
          <select className="Myselect" value={selectedOption} onChange={handleChange}>
            <option value="public">Public</option>
            <option value="Followers">Followers</option>
            <option value="Only">Only</option>
          </select>
        </div>
        <button type='submit'>post</button>
      </div>
    </form>
  );
}
