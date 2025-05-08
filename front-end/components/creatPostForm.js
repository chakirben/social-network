'use client'

import React, { useState, useRef } from 'react';
import Divider from './divider';
import UserData from "@/components/UserData";
export default function CreatePost() {
  const inputRef = useRef(null);
  const [imageSrc, setImageSrc] = useState(null);
  const [selectedOption, setSelectedOption] = useState('public');
  const [text, setText] = useState('');
  const [selectedUsers, setSelectedUsers] = useState([]);
  const [collaped, setcollapsed] = useState('');
  const handleImageClick = () => {
    inputRef.current.click();
  };
  const toggleUser = (id) => {
    setSelectedUsers((prev) =>
      prev.includes(id) ? prev.filter((uid) => uid !== id) : [...prev, id]
    );
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
    console.log(e.target.value);

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
      }); 16
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
  const showUsersList = () => {
    console.log("dkhm");
    if (collaped) {

      setcollapsed(false)
    } else {
      setcollapsed(true)

    }
  }
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
            <option value="followers">Followers</option>
            <option value="only">Only</option>
          </select>
          {selectedOption === "only" ? (
            <>
              <button className='thiary' onClick={showUsersList}>+ Select users</button>
              {collaped ? (
                <div className='FriendList' >
                  <div className='df sb center'>
                    <h4>Select users</h4>
                    <img src='/images/close.svg' className='icn' onClick={showUsersList}></img>
                  </div>
                  <div className='userList df cl gp12 start'>
                    <div className='df gp12'>
                      <input
                        type="checkbox"
                        checked={selectedUsers.includes(1)}
                        onChange={() => toggleUser(1)}
                        className="checkBox"
                      />
                      <UserData usr={{ firstName: "chakir", lastName: "ben", image: "", followers: "15" }} />
                    </div>
                    <div className='df gp12'>
                      <input
                        type="checkbox"
                        checked={selectedUsers.includes(1)}
                        onChange={() => toggleUser(1)}
                        className="checkBox"
                      />
                      <UserData usr={{ firstName: "chakir", lastName: "ben", image: "", followers: "15" }} />
                    </div>
                    <div className='df gp12'>
                      <input
                        type="checkbox"
                        checked={selectedUsers.includes(1)}
                        onChange={() => toggleUser(1)}
                        className="checkBox"
                      />
                      <UserData usr={{ firstName: "chakir", lastName: "ben", image: "", followers: "15" }} />
                    </div>
                    <div className='df gp12'>
                      <input
                        type="checkbox"
                        checked={selectedUsers.includes(1)}
                        onChange={() => toggleUser(1)}
                        className="checkBox"
                      />
                      <UserData usr={{ firstName: "chakir", lastName: "ben", image: "", followers: "15" }} />
                    </div>  
                  </div>
                  <button>select</button>
                </div>
              ) : ""}
            </>
          ) : ""}
          <button type='submit'>post</button>
        </div>
      </div>
    </form>
  );
}