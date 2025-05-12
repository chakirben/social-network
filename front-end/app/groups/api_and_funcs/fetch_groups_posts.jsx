"use client"

import { useEffect, useState } from "react"
import Post from "@/components/post"
import CreatPostInGroup from "@/components/groups/creat_postgroup"
import GroupInfo from "@/components/groups/group_info"
export default function GroupDetails({ groupId , title }) {    
    const [PostsGroup, setPostsGroup] = useState([])
    const [isLoading, setIsLoading] = useState(true)

    useEffect(() => {
        const fetchposts = async () => {
            try {
                const rep = await fetch(`http://localhost:8080/api/PostsGroups?id=${groupId}`, {
                    credentials: "include",
                    headers: { "Content-Type": "application/json" },
                })
                const PostsGroupData = await rep.json()
                setPostsGroup(PostsGroupData || [])
                
            } catch (error){
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
            <GroupInfo title={title}/>
            <CreatPostInGroup gpid={groupId} />
            {PostsGroup.map((pst) => (
                <Post key={pst.id} pst={pst} />
            ))}
        </div>
    )
}   