'use client'
import SideBar from "@/components/sidebar";
import { useEffect, useState } from "react";


export default function Profile() {
    const [profile, setProfile] = useState(null)
    useEffect(() => {
        async function fetchProfile() {
            try {
                const res = await fetch('http://localhost:3000/api/profile', { credentials: "include" })
                if (!res.ok) {
                    throw new Error(`HTTP error! Status: ${res.status}`);
                }
                let data = await res.json()
                console.log(data);

                setProfile(data)
            } catch (err) {
                console.log('errrrror', err)
            }
        }
        fetchProfile()
    }, [])

    return (
        <div className="profileContainer">
            <SideBar />

            {profile ? (
                <pre>{JSON.stringify(profile, null, 2)}</pre>
            ) : (
                <p>Loading...</p>
            )}
        </div>
    )
}

export default async function Profile() {

    const response = await fetch(`http://localhost:8080/api/profile`, {
        method: 'POST',
        credentials: 'include',  
        headers: {
          'Content-Type': 'application/json',
        },
        //here we should put the nickname of the owner of profile
        body: JSON.stringify({nickname: "iafriad"}),
    });
  
    if (!response.ok) {
        const resp = await response.text();
        console.log("Error get Data :", resp);
    } else {
        const data = await response.json();
        console.log("Success:");
        console.log(data);
    }
}