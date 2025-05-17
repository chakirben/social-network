'use client'

import React, { useState } from 'react';
import Divider from './divider';
import { useUser } from './userContext';
import "../styles/creatEvent.css"

export default function CreateEvent() {
  const { user } = useUser();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [err, setErr] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!title.trim() || !content.trim()) {
      setErr("Title and content are required");
      return;
    }

    const body = {
      title,
      content
    };

    try {
      const res = await fetch('http://localhost:8080/api/CreateEvent', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(body),
      });

      if (res.ok) {
        setTitle('');
        setContent('');
        setErr('');
        console.log('Event created successfully');
      } else {
        const errorText = await res.text();
        setErr(errorText || 'Failed to create event');
      }
    } catch (error) {
      console.error('Create event error:', error);
      setErr('Something went wrong.');
    }
  };

  return (
    <form className="creatPostForm">
      <div className="df center">
        <img className="avatar" src={user ? `http://localhost:8080/${user.avatar}` : ''} />
        <input
          className="searchInput"
          placeholder="Event title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
      </div>

      <textarea
        className="textArea"
        placeholder="Event details..."
        value={content}
        onChange={(e) => setContent(e.target.value)}
      />

      <Divider />

      <div className="spB pd8">
        <div></div>
        <button
          type="submit"
          onClick={handleSubmit}
          disabled={!title.trim() || !content.trim()}
          className={!title.trim() || !content.trim() ? 'button-disabled' : 'button-active'}
        >
          Create Event
        </button>
      </div>

      <div className='err'>{err}</div>
    </form>
  );
}
