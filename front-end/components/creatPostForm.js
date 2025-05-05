'use client'

import React, { useState, useRef } from 'react';
import Divider from './divider';

export default function CreatePost() {
  const inputRef = useRef(null);
  const [imageSrc, setImageSrc] = useState(null);
  const [selectedOption, setSelectedOption] = useState('');

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

  return (
    <div className="creatPostForm">
      <div className="searchBar">
        <img src="/user-icon.png" />
        <input
          className="searchInput"
          placeholder="What’s happening ?"
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

    </div>
  );
}
