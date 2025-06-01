"use client"
import { useRouter } from 'next/navigation';
import { useEffect, useState , useRef} from "react";
import Post from "@/components/post";
import "../../../styles/global.css";
import GroupEventsPage from '@/components/events/groupEventsPage';
import Divider from '@/components/divider';
import CreateEvent from '@/components/createEventForm';


export default function GroupDetails({ groupId, title }) {
    const [PostsGroup, setPostsGroup] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [events, setEvents] = useState([])
    const [activeTab, setActiveTab] = useState("posts");

    const [text, setText] = useState('');
    const [imageSrc, setImageSrc] = useState(null);
    const inputRef = useRef(null);

    

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
            const PostsGroupData = await rep.json();
            setPostsGroup(PostsGroupData || []);
        } catch (error) {
            console.error("Error fetching posts:", error);
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
                const res = await fetch('http://localhost:8080/api/CreatePost', {
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
        }
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



    return (
        <div className='postEventsInGroup'>
            <div className="GPTitle">
                <div className="divbtnbackgroups">
                    <img src="./../images/arrow-left.svg"  onClick={() => router.push(`/groups`)}/>
                    <button className="backbtngroups" onClick={() => router.push(`/groups`)}>back</button>
                </div>
                <div className="titleandimg">
                    <img src="./../images/group.svg" />
                    <p>{title}</p>
                </div>
            </div>

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
                        <div className="searchBar">
                            <img src="./../../../public/images/user-icon.png" />
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
                                <img
                                    src="./images/image.svg"
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
                    {PostsGroup.map((pst) => (
                        <Post key={pst.id} pst={pst} />
                    ))}
                </>
            )}

            {!isLoading && activeTab === "events" && (
                <>
                    <CreateEvent setEvents={setEvents} evnts={events} />
                    <GroupEventsPage id={groupId} events={events} />
                </>
            )}
        </div>
    );
}
