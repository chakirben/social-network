"use client"
import { useRouter } from 'next/navigation';
import { useEffect, useState, useRef, use } from "react";
import Post from "@/components/post";
import "../../../styles/global.css";
import GroupEventsPage from '@/components/events/groupEventsPage';
import Divider from '@/components/divider';
import CreateEvent from '@/components/createEventForm';
import Avatar from '@/components/avatar/avatar';
import { useUser } from '@/components/context/userContext';
import Header from '@/components/Header/header';
import InviteToGroups from "@/components/groups/infite_followers_to_join"


export default function GroupDetails({ groupId }) {
    const { user } = useUser();
    const [PostsGroup, setPostsGroup] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [events, setEvents] = useState([])
    const [activeTab, setActiveTab] = useState("posts");


    const [showInviteForm, setShowInviteForm] = useState(false);


    const [text, setText] = useState('');
    const [imageSrc, setImageSrc] = useState(null);
    const inputRef = useRef(null);


    const [error, setError] = useState("")


    useEffect(() => {
        console.log("Updated events:", events);
    }, [events]);
    const router = useRouter();

    const fetchPosts = async () => {
        setIsLoading(true);
        try {
            const rep = await fetch(`http://localhost:8080/api/PostsGroups?id=${groupId}`, {
                credentials: "include",
                headers: { "Content-Type": "application/json" },
            });
            if (rep.status === 404) {
                router.push('/404')
            }
            if (rep.status === 403) {
                setError('You are not allowed to access this group.')
            }
            const PostsGroupData = await rep.json();

            setPostsGroup(PostsGroupData || []);
        } catch (error) {
            console.error("Error fhomeetching posts:", error);
        } finally {
            setIsLoading(false);
        }
    };
    const fetchEvents = async () => {
        try {
            const res = await fetch(`http://localhost:8080/api/GetGroupEvents?id=${groupId}`, {
                credentials: "include",
                headers: { "Content-Type": "application/json" },
            });
            const data = await res.json();
            setEvents(data || []);


        } catch (error) {
            console.error("Error fetching events:", error);
        }
    };

    useEffect(() => {
        if (activeTab === "events") {
            fetchEvents();
        }
    }, [groupId, activeTab]);


    useEffect(() => {
        if (activeTab === "posts") {
            fetchPosts();
        }
    }, [groupId, activeTab]);


    const handleSubmit = async (e) => {
        e.preventDefault();
        const file = inputRef.current.files[0];
        const formData = new FormData();
        formData.append('content', text);
        formData.append('privacy', "inGroup");
        formData.append('groupId', groupId)
        if (file) {
            formData.append('image', file);
        }

        if (text.trim() == "") {
            return
        } else {

            try {
                const res = await homefetch('http://localhost:8080/api/CreatePost', {
                    method: 'POST',
                    body: formData,
                    credentials: 'include',

                });
                const result = await res.json();
                setPostsGroup((pv) => [result, ...pv])

                setText('');
                setImageSrc(null);
                inputRef.current.value = null;
            } catch (err) {
                console.error('Post failed:', err);
            }
        } home
    };


    const handleImageClick = () => {
        inputRef.current.click();
    };

    const handleFileChange = (e) => {
        const file = e.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = () => {
                setImageSrc(reader.result);
            };
            reader.readAsDataURL(file);
        }
    };

    if (isLoading) {
        return <div>Loading posts of the group...</div>;
    }

    const handlInviteClick = () => {
        setShowInviteForm(true)
    }


    return (
        <div className='postEventsInGroup'>

            {error != "" ? (
                <>{error}</>
            ) : (
                <>
                    <Header pageName={PostsGroup.GPTitle}
                        ele={
                            <button onClick={handlInviteClick} className="create-group-btn">
                                + INVITE
                            </button>
                        }
                    />
                    <div className='filterPostsAndEvents'>
                        <span
                            className={`postsSpan ${activeTab === "posts" ? "active" : ""}`}
                            onClick={() => setActiveTab("posts")}
                        >
                            Posts
                        </span>
                        <span
                            className={`eventsSpan ${activeTab === "events" ? "active" : ""}`}
                            onClick={() => setActiveTab("events")}
                        >
                            Events
                        </span>
                    </div>
                    <Divider />
                    {isLoading && <div>Loading...</div>}


                    {!isLoading && activeTab === "posts" && (
                        <>
                            <form className="creatPostForm" onSubmit={handleSubmit}>
                                <div className="searchBar df gp12 center">
                                    <Avatar url={user.avatar} name={user.firstName} />

                                    <input
                                        className="searchInput"
                                        placeholder="What’s happening ?"
                                        value={text}
                                        onChange={(e) => setText(e.target.value)}
                                    />
                                </div>

                                <div className="ImagePreviewBox">
                                    {imageSrc && <img src={imageSrc} alt="Preview" className="preview-img" />}
                                </div>

                                <Divider />

                                <div className='spB'>
                                    <div className='group'>
                                        <img style={{ width: "20px", height: "20px", cursor: "pointer" }}
                                            src="../../images/image.svg"
                                            className="upload-icon"
                                            onClick={handleImageClick}
                                        />
                                        <input
                                            type="file"
                                            accept="image/*"
                                            ref={inputRef}
                                            onChange={handleFileChange}
                                            style={{ display: 'none' }}
                                        />
                                    </div>
                                    <button type='submit'>post</button>
                                </div>
                            </form>
                            {Array.isArray(PostsGroup.allposts) && PostsGroup.allposts.length > 0 ?
                                PostsGroup.allposts.map((pst) => (
                                    <Post key={pst.id} pst={pst} />
                                )) : (
                                    <>{`No posts yet in this group creat a one.. (:`}</>
                                )}

                            {!isLoading && activeTab === "events" && (
                                <>
                                    <CreateEvent setEvents={setEvents} events={events} />
                                    {events.length > 0 ? (
                                        <GroupEventsPage id={groupId} events={events} setEvents={setEvents} />
                                    ) : (
                                        <div className="noEvents">
                                            <img className="noContent" src="/images/noContent.svg" alt="No content" />
                                            No events created, be the first
                                        </div>
                                    )}
                                </>
                            )}

                            {showInviteForm && (
                                <InviteToGroups groupId={groupId} onSkip={() => setShowInviteForm(false)} />
                            )}
                        </>
                    )}


                </>
            )
            }
        </div >
    );
}
