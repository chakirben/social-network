'use client'
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import '../register/register.css';
import InitWs from '../websocket/websocket';

export default function Login() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const handleLogin = async (e) => {
    e.preventDefault();
    if (!email || !password) {
      setErrorMessage('Please fill in all fields.');
      return;
    }

    const response = await fetch(`http://localhost:8080/api/login`, {
      method: 'POST',
      credentials: 'include',  
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      const resp = await response.text();
      setErrorMessage(resp || 'Login failed.');
    } else {
      await InitWs();
      console.log("success");
      
    }
  }

  return (
    <div className="container">
      <div className="formContainer">
        <h1>Login</h1>
        <p>Welcome back! Please login to your account.</p>

        <form className="form flex col" onSubmit={handleLogin}>
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
          {errorMessage && <p className="error">{errorMessage}</p>}
          <button type="submit">Login</button>
        </form>

        <p className="registerLink">
          Don't have an account? <a href="/register">Register</a>
        </p>
      </div>
    </div>
  );
}
