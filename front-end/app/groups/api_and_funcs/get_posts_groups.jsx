"use client"

import { useEffect, useState } from "react"
import Post from "@/components/post"
import CreatPostInGroup from "@/components/groups/creat_postgroup"
export default function GroupDetails({groupId , back}) {
    const [PostsGroup , setPostsGroup] = useState([])
    const [isLoading, setIsLoading] = useState(true);
        useEffect(() =>{
        const fetchposts = async () => {
            try {
                const rep = await fetch("http://localhost:8080/api/PostsGroups", {
                    credentials: "include" ,
                    method: "POST",
                    headers : {
                        "Content-Type" : "application/json"
                    },
                    body : JSON.stringify({
                        groupId : groupId
                    })
                })
                const PostsGroupData  = await rep.json()
                console.log(PostsGroupData)

                
                setPostsGroup(PostsGroupData || [])
            } catch {
                console.error("Error fetching groups:", error);
            } finally {
                setIsLoading(false);
            }


        } 
        fetchposts()
    },[])

    if (isLoading) {
        return <div>Loading posts of the group...</div>;
    }

    if (PostsGroup.length === 0) {
        return <div>There is no posts of this group Creat a one...</div>;
    }

    return (
        <div>
            <CreatPostInGroup gpid={groupId}/>
            {PostsGroup.map((pst)=> (
                <Post key={pst.id} pst={pst} />
            ))}
        </div>
    )
}   