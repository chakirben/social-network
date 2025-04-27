
import '../public/login.css'
import "../public/style.css"
import {  useState } from 'react'
import { useRouter } from 'next/router'

export default function LoginPage() {
    const router = useRouter()
    const [errorMessage, setErrorMessage] = useState('')

    async function handleSubmit(event) {
        event.preventDefault()

        const formData = new FormData(event.currentTarget)
        const email = formData.get('email')?.toString().trim()
        const password = formData.get('password')?.toString().trim()

        if (!email || !password) {
            setErrorMessage('Please enter both email and password')
            return
        }

        const response = await fetch('https://localhost:8080/api/Login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password }),
        })

        if (response.ok) {
            router.push('/')
        } else {
            const errorText = await response.text()
            setErrorMessage(errorText || 'Failed to login')
            console.log('Login error:', errorText)
        }
    }

    return (
        <div className="tocenter">
            <div className='texts'>
                <img src='/images/logo.svg'></img>
                <div>Login to access premium posts</div>
            </div>

            <form onSubmit={handleSubmit}>
                <input type="email" name="email" placeholder="Email" required />
                <input type="password" name="password" placeholder="Password" required />
                {errorMessage && <div style={{ color: 'red' }}>{errorMessage}</div>}
                <button type="submit">Login</button>
            </form>

            <div>
                Don't have an account?{' '}
                <a href="/register">Register</a>
            </div>
        </div>
    )
}
