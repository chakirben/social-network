import { useState } from 'react';
import '../public/styles/global.css';
import '../public/styles/register.css';

export default function Register() {
  const [avatar, setAvatar] = useState(null);
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [dob, setDob] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [nickname, setNickname] = useState('');
  const [about, setAbout] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const handleAvatarChange = (e) => {
    if (e.target.files && e.target.files[0]) {
      setAvatar(URL.createObjectURL(e.target.files[0]));
    }
  };

  const validateForm = (firstName, lastName, dob, email, password) => {
    if (!firstName || !lastName || !dob || !email || !password) {
      return { valid: false, message: "All required fields must be filled" };
    }

    const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/;
    if (!emailRegex.test(email)) {
      return { valid: false, message: "Please enter a valid email address" };
    }

    if (password.length < 6) {
      return { valid: false, message: "Password must be at least 6 characters long" };
    }
    const validateDob = (dob) => {
      const age = (new Date() - new Date(dob)) / (1000 * 60 * 60 * 24 * 365);
      return age >= 16;
    };
    if (!validateDob(dob)) {
      return { valid: false, message: "user must have more than 16yo!" };
    }
    return { valid: true, message: "" };
  };

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
    if (avatar) {
      const avatarFile = document.querySelector('#avatarUpload').files[0];
      formData.append('avatar', avatarFile);
    }

    try {
      const response = await fetch(`http://localhost:8080/api/register`, {
        method: 'POST',
        body: formData, 
      });

      if (!response.ok) {
        throw new Error('Registration failed');
      }

      const data = await response.json();
    } catch (error) {
      console.error(error);
      setErrorMessage('An error occurred during registration');
    }
  };


  return (
    <div className="container">
      <div className="formContainer">
        <h1>Register</h1>
        <p>Enter your details</p>



        <div className="avatarContainer">
          <label htmlFor="avatarUpload" className="avatarLabel">
            {avatar ? (
              <img src={avatar} alt="Avatar" className="avatarImage" />
            ) : (
              <div className="avatarPlaceholder">👤</div>
            )}
          </label>
          <input
            id="avatarUpload"
            type="file"
            accept="image/*"
            onChange={handleAvatarChange}
            style={{ display: 'none' }}
          />
        </div>

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
          Don't have an account? <a href="/register">Register</a>
        </p>
      </div>
    </div>
  );
}