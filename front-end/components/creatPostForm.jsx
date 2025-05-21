'use client'

import React, { useState, useRef, useEffect } from 'react';
import Divider from './divider';
import UserData from "@/components/UserData";
import { useUser } from './userContext';
export default function CreatePost({newpost}) {
  const inputRef = useRef(null);
  const [imageSrc, setImageSrc] = useState(null);
  const [selectedOption, setSelectedOption] = useState('public');
  const [text, setText] = useState('');
  const [selectedUsers, setSelectedUsers] = useState([]);
  const [users, setUsers] = useState([]);
  const [collapsed, setCollapsed] = useState(false);
  const [err , setErr]  =  useState("")
  const { user, setUser } = useUser();

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await fetch("http://localhost:8080/api/getFollowersList", { credentials: "include" });
        const data = await response.json();
        setUsers(data);
      } catch (error) {
        console.error("Error fetching posts:", error);
      }
    };
    fetchUsers();
  }, []);

  const handleImageClick = () => {
    inputRef.current.click();
  };

  const toggleUser = (id) => {

    setSelectedUsers((prev) =>
      prev.includes(id) ? prev.filter((uid) => uid !== id) : [...prev, id]
    );
    console.log("selected users: ", selectedUsers);

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
    if (selectedOption === "only"  && !selectedUsers.length){
      setErr("please select at least 1 user")
      return
    }
    const file = inputRef.current.files[0];
    const formData = new FormData();
    formData.append('content', text);
    formData.append('privacy', "inGroup");
    selectedUsers.forEach(id => {
      formData.append('selectedUsers', id);
    });
    if (file) {
      formData.append('image', file);
    }

    try {
      const res = await fetch('http://localhost:8080/api/CreatePost', {
        method: 'POST',
        body: formData,
        credentials: 'include',
      });
      if (res.ok) {
        console.log('Post submitted:', res);
        setText('');
        setSelectedUsers(null)
        setErr(null)
        setSelectedOption('public');
        setImageSrc(null);
        inputRef.current.value = null;
      }

       const result = await res.json();
       newpost(result);

    } catch (err) {
      console.error('Post failed:', err);
    }
  };

  const showUsersList = (e) => {
    e.preventDefault();
    setCollapsed(!collapsed);
  };

  return (
    <form className="creatPostForm">
      <div className="df center">
        <img className="avatar" src={user ? `http://localhost:8080/${user.avatar}`: ""} />
        <input
          className="searchInput"
          placeholder="What's happening ?"
          value={text}
          onChange={(e) => setText(e.target.value)}
        />
      </div>

      <div className="ImagePreviewBox">
        {imageSrc && <img src={imageSrc} alt="Preview" className="preview-img" />}
      </div>

      <Divider />

      <div className='spB pd8'>
        <div className='group'>
          <img
            src="./images/image.svg"
            className="upload-icon"
            onClick={handleImageClick}
            alt="Upload"
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
            <option value="almostPrivate">Followers</option>
            <option value="private">Only</option>
          </select>
          {selectedOption === "private" && (
            <>
              <button className='thiary' onClick={(e) => showUsersList(e)}> {selectedUsers && selectedUsers.length ? `${selectedUsers.length} selected users ✔️` : "+ Select users"}</button>
              {collapsed && (
                <div className='FriendList'>
                  <div className='df sb center'>
                    <h4>Select users</h4>
                    <img src='/images/close.svg' className='icn' alt="Close" onClick={(e) => { e.preventDefault(); setCollapsed(false); }} />
                  </div>
                  <div className='userList df cl gp12 start'>
                    {!users? <div className='df gp12' >no follow</div> : 
                    users.map((user) => (
                      <div key={user.id} className='df gp12'>
                        <input
                          type="checkbox"
                          checked={selectedUsers.includes(user.id)}
                          onChange={() => toggleUser(user.id)}
                          className="checkBox"
                        />
                        <UserData usr={user} />
                      </div>
                    )) }
                  </div>
                  <button onClick={(e) => { e.preventDefault(); setCollapsed(false); }}>Select</button>
                </div>
              )}
            </>
          )}
        </div>
          <button
            type='submit'
            onClick={handleSubmit}
            disabled={!text.trim()}
            className={!text.trim() ? 'button-disabled' : 'button-active'}
          >
            Post
          </button>

      </div>
      <div className='err'>{err}</div>
    </form>
  );
}