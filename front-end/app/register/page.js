'use client'
import { useState } from 'react';
import '../../styles/global.css';
import "./register.css"
//import '../public/styles/register.css';
import { useRouter } from 'next/navigation';
import AvatarUpload from '@/public/components/avatar';
import { validateForm } from '@/public/utils/formValidation';

export default function Register() {
  const router = useRouter();
  const [avatar, setAvatar] = useState(null);
  const [avatarFile, setAvatarFile] = useState(null);
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [dob, setDob] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [nickname, setNickname] = useState('');
  const [about, setAbout] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const handleRegister = async (e) => {
    e.preventDefault();

    const validation = validateForm(firstName, lastName, dob, email, password);

    if (!validation.valid) {
      setErrorMessage(validation.message);
      return;
    }

    const formData = new FormData();
    formData.append('firstName', firstName);
    formData.append('lastName', lastName);
    formData.append('dob', dob);
    formData.append('email', email);
    formData.append('password', password);
    formData.append('nickname', nickname);
    formData.append('about', about);
    if (avatarFile) {
      formData.append('avatar', avatarFile);
    }

    const response = await fetch(`http://localhost:8080/api/register`, {
      method: 'POST',
      credentials: 'include',
      body: formData,
    });

    if (!response.ok) {
      let resp = await response.text();
      setErrorMessage(resp);
    } else {
      router.push('/home');
    }
  }

  return (
    <div className="container">
      <div className="formContainer">
        <h1>Register</h1>
        <p>Enter your details</p>

        <AvatarUpload avatar={avatar} onAvatarChange={(url, file) => {
          setAvatar(url);
          setAvatarFile(file);
        }} />

        <form className="form flex col" onSubmit={handleRegister}>
          <div className="flex">
            <input
              type="text"
              className="grw"
              placeholder="First name"
              value={firstName}
              onChange={(e) => setFirstName(e.target.value)}
              required
            />
            <input
              type="text"
              className="grw"
              placeholder="Last name"
              value={lastName}
              onChange={(e) => setLastName(e.target.value)}
              required
            />
          </div>
          <input
            type="date"
            placeholder="Date of birth"
            value={dob}
            onChange={(e) => setDob(e.target.value)}
            required
          />
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <input
            type="text"
            placeholder="Nickname"
            value={nickname}
            onChange={(e) => setNickname(e.target.value)}
          />
          <textarea
            placeholder="About"
            value={about}
            onChange={(e) => setAbout(e.target.value)}
          />
          {errorMessage && <p className="error">{errorMessage}</p>}
          <button type="submit">Register</button>
        </form>

        <p className="registerLink">
          Already have an account? <a href="/login">Login</a>
        </p>
      </div>
    </div>
  );
}
