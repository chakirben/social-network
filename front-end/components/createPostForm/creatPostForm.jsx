'use client'

import React, { useState, useRef, useEffect } from 'react';
import Divider from '../divider';
import UserData from "@/components/UserData";
import { useUser } from '../context/userContext';
import styles from './create.module.css';
import Avatar from '../avatar/avatar';

export default function CreatePost({ newpost }) {
  const inputRef = useRef(null);
  const [imageSrc, setImageSrc] = useState(null);
  const [selectedOption, setSelectedOption] = useState('public');
  const [text, setText] = useState('');
  const [selectedUsers, setSelectedUsers] = useState([]);
  const [users, setUsers] = useState([]);
  const [collapsed, setCollapsed] = useState(false);
  const [err, setErr] = useState("");
  const [loadingUsers, setLoadingUsers] = useState(false);
  const { user } = useUser();

  const fetchUsers = async () => {
    setLoadingUsers(true);
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/getFollowersList`, { credentials: "include" });
      const data = await response.json();
      setUsers(data);
    } catch (error) {
      console.error("Error fetching followers:", error);
      setErr("Failed to load followers list");
    } finally {
      setLoadingUsers(false);
    }
  };

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
    // Clear any previous error when privacy changes
    setErr("");
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    // Validate selected users only if privacy is "private"
    if (selectedOption === "private" && selectedUsers.length === 0) {
      setErr("Please select at least 1 user");
      return;
    }

    const file = inputRef.current.files[0];
    const formData = new FormData();
    formData.append('content', text);
    formData.append('privacy', selectedOption);

    // Append selected users only if privacy is private
    if (selectedOption === "private") {
      selectedUsers.forEach(id => {
        formData.append('selectedUsers', id);
      });
    }

    if (file) {
      formData.append('image', file);
    }

    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/CreatePost`, {
        method: 'POST',
        body: formData,
        credentials: 'include',
      });

      if (res.ok) {
        const result = await res.json();
        newpost(result);

        setText('');
        setSelectedUsers([]);
        setErr('');
        setSelectedOption('public');
        setImageSrc(null);
        inputRef.current.value = null;
      } else {
        const errorText = await res.text();
        setErr(`Post failed: ${errorText}`);
      }
    } catch (err) {
      console.error('Post failed:', err);
      setErr('Post failed: Network error');
    }
  };

  const showUsersList = async (e) => {
    e.preventDefault();
    if (!collapsed) {
      // Only fetch when opening the modal
      await fetchUsers();
    }
    setCollapsed(!collapsed);
  };

  return (
    <form className="creatPostForm" onSubmit={handleSubmit}>
      <div className="df center">
        <Avatar url={user?.avatar} name={user?.firstName} />

        <input
          className={styles.searchInput}
          placeholder="What's happening?"
          value={text}
          onChange={(e) => setText(e.target.value)}
        />
      </div>

      <div className="ImagePreviewBox">
        {imageSrc && <img src={imageSrc} alt="Preview" className="preview-img" />}
      </div>

      <Divider />

      <div className="spB pd8">
        <div className="group">
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
              <button className="thiary" onClick={showUsersList}>
                {selectedUsers.length ? `${selectedUsers.length} selected users ✔️` : "+ Select users"}
              </button>
              {collapsed && (
                <div className="FriendList">
                  <div className="df sb center">
                    <h4>Select users</h4>
                    <img
                      src="/images/close.svg"
                      className="icn"
                      alt="Close"
                      onClick={(e) => {
                        e.preventDefault();
                        setCollapsed(false);
                      }}
                    />
                  </div>
                  <div className="userList df cl gp12 start">
                    {loadingUsers ? (
                      <div className="df gp12">Loading followers...</div>
                    ) : users?.length === 0 ? (
                      <div className="df gp12">No followers</div>
                    ) : (
                      users?.map((user) => (
                        <div key={user.id} className="df gp12">
                          <input
                            type="checkbox"
                            checked={selectedUsers.includes(user.id)}
                            onChange={() => toggleUser(user.id)}
                            className="checkBox"
                          />
                          <UserData usr={user} />
                        </div>
                      ))
                    )}
                  </div>
                  <button
                    onClick={(e) => {
                      e.preventDefault();
                      setCollapsed(false);
                    }}
                  >
                    Select
                  </button>
                </div>
              )}
            </>
          )}
        </div>

        <button
          type="submit"
          disabled={!text.trim()}
          className={!text.trim() ? 'button-disabled' : 'button-active'}
        >
          Post
        </button>
      </div>

      <div className="err">{err}</div>
    </form>
  );
}