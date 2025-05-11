
'use client'
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import FollowButton from './FollowButton';

export default function ProfileClient({ session, searchParams }) {
    const [profileId, setProfileId] = useState(null);

    const { id } = useParams() || searchParams;
    useEffect(() => {
        setProfileId(id);
    }, [id]);



    const handleFetch = async () => {
        const response = await fetch('http://localhost:8080/api/profile', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ id: id }),
        });
        
        if (!response.ok) {
            const resp = await response.text();
            console.error('Error get Data:', resp);
        } else {
            const data = await response.json();
            console.log('Success:', data);
        }
    }

    useEffect(() => {
        console.log("Profile ID from Params: ", id);
        console.log("Session ID: ", session);
        handleFetch();
    }, []);
    

    return (
        <div>
            <FollowButton session={session} id={id} />
        </div>
    );
}