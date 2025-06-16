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
      console.log("Login success")
      Connect()

  
        async function FirstTimeUser () {
           try {
             const rep = await fetch("http://localhost:8080/api/getUserData", {
                credentials: "include" 
               })
               if (!rep.ok) {
                   throw new Error(`HTTP error! Status: ${res.status}`);
               }
                let data = await rep.json()
                localStorage.setItem("user", JSON.stringify(data));
                setUser(data)
           } catch (err) {
             console.log('errrrror', err)
           }
        }
        FirstTimeUser()

      router.push('/home');
    }
  }

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
