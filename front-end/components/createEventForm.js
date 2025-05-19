'use client'

import React, { useState } from 'react';
import Divider from './divider';
import { useUser } from './userContext';
import "../styles/creatEvent.css"
import { useParams } from 'next/navigation';

export default function CreateEvent({ setEvents, evnts }) {
  //  console.log("fffffff",evnts , setEvents);
  const { user } = useUser();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [err, setErr] = useState('');
  const [eventDate, setEventDate] = useState('');
  const { id } = useParams()


  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!title.trim() || !content.trim()) {
      setErr("Title and content are required");
      return;
    }

    const body = {
      title,
      description: content,
      groupId: parseInt(id, 10),
      eventDate: new Date(eventDate)
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
        setEventDate('');
        setErr('');
        const newEvent = await res.json()

        setEvents(prevEvents => [newEvent, ...prevEvents])
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
      <div className="df center gp12">
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
        <div className="calendar-section">
          <span className="eventsDate">Choose date</span>
          <input
            type="date"
            value={eventDate}
            onChange={(e) => setEventDate(e.target.value)}  // Update the eventDate state
            placeholder="Select an event date"
          />
        </div>

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
