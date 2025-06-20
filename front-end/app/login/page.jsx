'use client'
import { useState, useContext, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { WebSocketContext } from '@/components/context/wsContext';

import { useUser } from '@/components/context/userContext'
import '../register/register.css';
import styles from './login.module.css'
export default function Login() {
  const { setUser } = useUser();

  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const { Connect } = useContext(WebSocketContext);

  const handleLogin = async (e) => {
    e.preventDefault();
    if (!email || !password) {
      setErrorMessage('Please fill in all fields.');
      return;
    }
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/login`, {
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
        return;
      }

      if (Connect) {
        Connect();
      }

      const fetchUser = async () => {
        try {
          const rep = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/getUserData`, {
            credentials: "include",
          });

          if (!rep.ok) {
            throw new Error(`HTTP error! Status: ${rep.status}`);
          }

          const data = await rep.json();
          setUser(data);
        } catch (err) {
          console.log("Error fetching user:", err);
        }

      };
      fetchUser();


      console.log("Login successful");

      router.push('/home');

    } catch (err) {
      console.error("Login fetch failed", err);
      setErrorMessage("Network error. Please try again.");
    }
  };
  return (
    <div className="container">
      <div className="formContainer">
        <div className={styles.header}>
          <h1>Login</h1>
          <p>Welcome back! Please login.</p>
        </div>
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
          Don't have an account? <a onClick={() => { router.push("/register") }}>Register</a>
        </p>
      </div>
    </div>
  );
}
