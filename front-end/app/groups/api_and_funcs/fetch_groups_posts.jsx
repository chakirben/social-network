"use client"
import { useRouter } from 'next/navigation';
import { useEffect, useState } from "react"
import Post from "@/components/post"
import CreatPostInGroup from "@/components/groups/creat_postgroup"

export default function GroupDetails({ groupId, title }) {
    const [PostsGroup, setPostsGroup] = useState([])
    const [isLoading, setIsLoading] = useState(true)

    const router = useRouter();
    
    useEffect(() => {
        const fetchposts = async () => {
            try {
                const rep = await fetch(`http://localhost:8080/api/PostsGroups?id=${groupId}`, {
                    credentials: "include",
                    headers: { "Content-Type": "application/json" },
                })
                const PostsGroupData = await rep.json()
                setPostsGroup(PostsGroupData || [])

            } catch (error) {
                console.error("Error fetching groups:", error);
            } finally {
                setIsLoading(false);
            }


        }
        fetchposts()
    }, [groupId])

    if (isLoading) {
        return <div>Loading posts of the group...</div>;
    }

    return (
        <div>
            <div className="GPTitle">
                <div className="btnback">
                    <img src="./../images/arrow-left.svg"/>
                    <button className="backbtn" onClick={()=> router.push(`/groups`)}>back</button>
                </div>
                <div className="titleandimg" >
                    <img src="./../images/group.svg" />
                    <p>{title}</p>
                </div>
            </div>
            <CreatPostInGroup gpid={groupId} />
            {PostsGroup.map((pst) => (
                <Post key={pst.id} pst={pst} />
            ))}
        </div>
    )
}   